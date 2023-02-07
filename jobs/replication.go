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
)

// ReplicationProcessor check replication if exists
type ReplicationProcessor struct {
	Processor
}

func NewReplicationProcessor(ln *core.LightNode) ReplicationProcessor {
	return ReplicationProcessor{
		Processor{
			LightNode: ln,
		},
	}
}

func (r *ReplicationProcessor) Run() {

	// get all piece comm record that are replication-assigned
	var contents []core.Content
	r.LightNode.DB.Model(&core.Content{}).Where("status = ?", "piece-assigned").Find(&contents)

	// for each content, get piece commitment and process
	for _, content := range contents {
		var pieceCommitments []core.PieceCommitment
		r.LightNode.DB.Model(&core.PieceCommitment{}).Where("piece_commitment_id = ?", content.PieceCommitmentId).Find(&pieceCommitments)
		for _, pieceCommitment := range pieceCommitments {
			err := r.makeStorageDeal(content, pieceCommitment)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func (r *ReplicationProcessor) makeStorageDeal(content core.Content, pieceComm core.PieceCommitment) error {

	//var pieceCommitment core.PieceCommitment
	//r.LightNode.DB.Model(&core.PieceCommitment{}).Where("piece = ?", bucketReplicationRequests.PieceCommitment).Find(&pieceCommitment)

	bal, err := r.LightNode.Filclient.Balance(context.Background())
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

	prop, err := r.LightNode.Filclient.MakeDeal(context.Background(), r.GetStorageProviders()[0].Address, pCid, priceBigInt, abi.PaddedPieceSize(pieceComm.PaddedPieceSize), duration, true)

	if err != nil {
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

	proto, err := r.LightNode.Filclient.DealProtocolForMiner(context.Background(), r.GetStorageProviders()[0].Address)
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
		return xerrors.Errorf("failed to create database entry for deal: %w", err)
	}

	fmt.Println(r.GetStorageProviders()[0].Address)
	propPhase, err := r.sendProposalV120(context.Background(), *prop, propnd.Cid(), dealUUID, deal.ID)

	if propPhase == true {
		pieceComm.Status = "complete"
		r.LightNode.DB.Save(&pieceComm)
	}

	return nil

}

type MinerAddress struct {
	Address address.Address
}

func (r *ReplicationProcessor) GetStorageProviders() []MinerAddress {
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

func (r *ReplicationProcessor) sendProposalV120(ctx context.Context, netprop network.Proposal, propCid cid.Cid, dealUUID uuid.UUID, dbid uint) (bool, error) {
	// In deal protocol v120 the transfer will be initiated by the
	// storage provider (a pull transfer) so we need to prepare for
	// the data request

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

	fmt.Println("AnnounceAddrs", r.LightNode.Node.Config.AnnounceAddrs[0])
	addrstr := r.LightNode.Node.Config.AnnounceAddrs[0] + "/p2p/" + r.LightNode.Node.Host.ID().String()
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
	return propPhase, err
}
