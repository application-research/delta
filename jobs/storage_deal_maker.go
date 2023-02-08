package jobs

import (
	"context"
	"fc-deal-making-service/core"
	"fmt"
	"github.com/filecoin-project/boost/transport/httptransport"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/xerrors"
	"strings"
)

type StorageDealMakerProcessor struct {
	Processor
}

type ItemReplicationProcessor struct {
	ReplicationProcessor
}

func (i ItemReplicationProcessor) Run() error {
	err := i.makeStorageDeal(i.Content, i.PieceComm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewItemReplicationProcessor(ln *core.LightNode, content core.Content, commitment core.PieceCommitment) IProcessor {
	return &ItemReplicationProcessor{
		ReplicationProcessor{
			LightNode: ln,
			Content:   &content,
			PieceComm: &commitment,
			Context:   context.Background(),
		},
	}
}

func NewStorageDealMakerProcessor(ln *core.LightNode) IProcessor {
	return &StorageDealMakerProcessor{
		Processor{
			LightNode: ln,
		},
	}
}

func (r *StorageDealMakerProcessor) Run() error {

	// get all piece comm record that are replication-assigned
	dispatcher := core.CreateNewDispatcher()
	var contents []core.Content
	r.LightNode.DB.Model(&core.Content{}).Where("status = ?", "piece-assigned").Find(&contents)

	// for each content, get piece commitment and process
	for _, content := range contents {
		var pieceCommitments []core.PieceCommitment
		r.LightNode.DB.Model(&core.PieceCommitment{}).Where("id = ?", content.PieceCommitmentId).Find(&pieceCommitments)

		for _, pieceCommitment := range pieceCommitments {
			item := NewItemReplicationProcessor(r.LightNode, content, pieceCommitment)
			dispatcher.AddJob(item)
		}
	}
	dispatcher.Start(10)
	return nil
}

func (r *ItemReplicationProcessor) makeStorageDeal(content *core.Content, pieceComm *core.PieceCommitment) error {

	//var pieceCommitment core.PieceCommitment
	//r.LightNode.DB.Model(&core.PieceCommitment{}).Where("piece = ?", bucketReplicationRequests.PieceCommitment).Find(&pieceCommitment)

	bal, err := r.LightNode.Filclient.Balance(r.Context)
	if err != nil {
		return err
	}

	fmt.Println(r.LightNode.Filclient.ClientAddr)

	fmt.Println("balance", bal.Balance.String())
	fmt.Println("escrow", bal.MarketEscrow.String())

	// 6 deals
	fmt.Println("pieceCommitment.Cid", pieceComm.Cid)
	pCid, err := cid.Decode(pieceComm.Cid)
	if err != nil {

	}
	fmt.Println("pCid", pCid.String())
	priceBigInt, err := types.BigFromString("000000000000002")
	var DealDuration = 1555200 - (2880 * 21)
	duration := abi.ChainEpoch(DealDuration)

	prop, err := r.LightNode.Filclient.MakeDeal(r.Context, r.GetStorageProviders()[0].Address, pCid, priceBigInt, abi.PaddedPieceSize(pieceComm.PaddedPieceSize), duration, true)
	fmt.Println(prop)
	if err != nil {
		fmt.Println("Prop", err)
		r.LightNode.Dispatcher.AddJob(NewItemReplicationProcessor(r.LightNode, *content, *pieceComm))
		return err
	}

	propnd, err := cborutil.AsIpld(prop.DealProposal)
	if err != nil {
		fmt.Println("Propnd", err)
		return err
	}

	dealUUID := uuid.New()

	if err != nil {
		fmt.Println("Prop", err)
		//return err
	}

	proto, err := r.LightNode.Filclient.DealProtocolForMiner(r.Context, r.GetStorageProviders()[0].Address)
	if err != nil {
		fmt.Println("Proto ", err)
		//return err
	}
	fmt.Println(proto)

	deal := &core.ContentDeal{
		Content:             content.ID,
		PropCid:             propnd.Cid().String(),
		DealUUID:            dealUUID.String(),
		Miner:               r.GetStorageProviders()[0].Address.String(),
		Verified:            true,
		DealProtocolVersion: proto,
		//MinerVersion:        ask.MinerVersion,
	}
	if err := r.LightNode.DB.Create(deal).Error; err != nil {
		r.LightNode.Dispatcher.AddJob(NewItemReplicationProcessor(r.LightNode, *content, *pieceComm))
		return xerrors.Errorf("failed to create database entry for deal: %w", err)
	}

	fmt.Println(r.GetStorageProviders()[0].Address)
	propPhase, err := r.sendProposalV120(r.Context, *prop, propnd.Cid(), dealUUID, deal.ID)
	if err != nil {
		fmt.Println("PropPhase", err)
		// get this request back to the queue
		r.LightNode.Dispatcher.AddJob(NewItemReplicationProcessor(r.LightNode, *content, *pieceComm))
		return err
	}
	fmt.Println("propPhase", propPhase)
	if propPhase == true {
		pieceComm.Status = "complete"
		content.Status = "replication-complete"
		r.LightNode.DB.Save(&pieceComm)
		r.LightNode.DB.Save(&content)
	} else {
		r.LightNode.Dispatcher.AddJob(NewItemReplicationProcessor(r.LightNode, *content, *pieceComm))
	}

	return nil

}

type MinerAddress struct {
	Address address.Address
}

func (r *ItemReplicationProcessor) GetStorageProviders() []MinerAddress {
	var storageProviders []MinerAddress
	for _, s := range mainnetMinerStrs {
		address.CurrentNetwork = address.Mainnet
		a, err := address.NewFromString(s)
		if err != nil {
			fmt.Println("error on miner address", err, a)
		}
		storageProviders = append(storageProviders, MinerAddress{Address: a})
	}
	return storageProviders
}

var mainnetMinerStrs = []string{
	"f01963614",
}

func (r *ItemReplicationProcessor) sendProposalV120(ctx context.Context, netprop network.Proposal, propCid cid.Cid, dealUUID uuid.UUID, dbid uint) (bool, error) {
	// In deal protocol v120 the transfer will be initiated by the

	// Create an auth token to be used in the request
	authToken, err := httptransport.GenerateAuthToken()
	if err != nil {
		return false, xerrors.Errorf("generating auth token for deal: %w", err)
	}

	rootCid := netprop.Piece.Root
	size := netprop.Piece.RawBlockSize
	var announceAddr multiaddr.Multiaddr

	if len(r.LightNode.Node.Config.AnnounceAddrs) == 0 {
		return false, xerrors.Errorf("cannot serve deal data: no announce address configured for estuary node")
	}

	fmt.Println("AnnounceAddrs", r.LightNode.Node.Config.AnnounceAddrs[1])
	addrstr := r.LightNode.Node.Config.AnnounceAddrs[1] + "/p2p/" + r.LightNode.Node.Host.ID().String()
	announceAddr, err = multiaddr.NewMultiaddr(addrstr)
	if err != nil {
		return false, xerrors.Errorf("cannot parse announce address '%s': %w", addrstr, err)
	}

	// Add an auth token for the data to the auth DB
	err = r.LightNode.Filclient.Libp2pTransferMgr.PrepareForDataRequest(ctx, dbid, authToken, propCid, rootCid, size)
	if err != nil {
		return false, xerrors.Errorf("preparing for data request: %w", err)
	}

	// Send the deal proposal to the storage provider
	propPhase, err := r.LightNode.Filclient.SendProposalV120(ctx, dbid, netprop, dealUUID, announceAddr, authToken)

	if err != nil {
		//  deal proposal is identical
		// if err includes message  deal proposal is identical
		if strings.ContainsAny(err.Error(), "deal proposal is identical") { // don't put it back on the queue
			fmt.Println("identical!!!", err)
			return true, err
		}

		r.LightNode.Filclient.Libp2pTransferMgr.CleanupPreparedRequest(r.Context, dbid, authToken)
	}
	return propPhase, err
}
