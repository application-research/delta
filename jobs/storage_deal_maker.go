package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
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
	"gorm.io/gorm"
	"strings"
	"time"
)

type StorageDealMakerProcessor struct {
	Context   context.Context
	LightNode *core.DeltaNode
	Content   *core.Content
	PieceComm *core.PieceCommitment
}

func (i StorageDealMakerProcessor) Run() error {
	err := i.makeStorageDeal(i.Content, i.PieceComm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func NewStorageDealMakerProcessor(ln *core.DeltaNode, content core.Content, commitment core.PieceCommitment) IProcessor {
	return &StorageDealMakerProcessor{
		LightNode: ln,
		Content:   &content,
		PieceComm: &commitment,
		Context:   context.Background(),
	}
}

func (i *StorageDealMakerProcessor) makeStorageDeal(content *core.Content, pieceComm *core.PieceCommitment) error {

	var minerAddress = i.GetAssignedMinerForContent(*content).Address
	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
		Status: utils.CONTENT_DEAL_MAKING_PROPOSAL, //"making-deal-proposal",
	})

	bal, err := i.LightNode.FilClient.Balance(i.Context)
	if err != nil {
		return err
	}

	fmt.Println(i.LightNode.FilClient.ClientAddr)
	fmt.Println("balance", bal.Balance.String())
	fmt.Println("escrow", bal.MarketEscrow.String())

	pCid, err := cid.Decode(pieceComm.Cid)
	if err != nil {
		fmt.Println("piece cid decode", err)
	}

	priceBigInt, err := types.BigFromString("0")
	var DealDuration = 1555200 - (2880 * 21)
	duration := abi.ChainEpoch(DealDuration)

	prop, err := i.LightNode.FilClient.MakeDeal(i.Context, minerAddress, pCid, priceBigInt, abi.PaddedPieceSize(pieceComm.PaddedPieceSize), duration, true)
	fmt.Println(prop)
	if err != nil {
		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm))
		return err
	}

	propnd, err := cborutil.AsIpld(prop.DealProposal)
	if err != nil {
		fmt.Println("proposal nd", err)
		return err
	}

	dealUUID := uuid.New()
	proto, err := i.LightNode.FilClient.DealProtocolForMiner(i.Context, minerAddress)
	if err != nil {
		fmt.Println("deal protocol for miner", err)
	}

	deal := &core.ContentDeal{
		Content:             content.ID,
		PropCid:             propnd.Cid().String(),
		DealUUID:            dealUUID.String(),
		Miner:               minerAddress.String(),
		Verified:            true,
		DealProtocolVersion: proto,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		//MinerVersion:        ask.MinerVersion,
	}
	if err := i.LightNode.DB.Create(deal).Error; err != nil {
		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm))
		return xerrors.Errorf("failed to create database entry for deal: %w", err)
	}

	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
		Status: utils.CONTENT_DEAL_SENDING_PROPOSAL, //"sending-deal-proposal",
	})

	propPhase, err := i.sendProposalV120(i.Context, *prop, propnd.Cid(), dealUUID, uint(deal.ID))

	if propPhase == true && err != nil {

		// TODO: better error handling
		if strings.Contains(err.Error(), "deal proposal is identical") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			return err
		}
		if strings.Contains(err.Error(), " piece size less than minimum required size") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			return err
		}

		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm))

		return err
	}

	if propPhase == false && err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status: utils.CONTENT_DEAL_PROPOSAL_SENT, //"deal-proposal-sent",
		})
	}

	// Online - transfer it
	if propPhase == false && content.ConnectionMode == "online" {
		propCid, err := cid.Decode(deal.PropCid)
		contentCid, err := cid.Decode(content.Cid)
		channelId, err := i.LightNode.FilClient.StartDataTransfer(i.Context, i.GetAssignedMinerForContent(*content).Address, propCid, contentCid)
		if err != nil {
			i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm))
			return err
		}
		content.PieceCommitmentId = pieceComm.ID
		pieceComm.Status = utils.COMMP_STATUS_COMITTED        //"committed"
		content.Status = utils.DEAL_STATUS_TRANSFER_STARTED   //"transfer-started"
		deal.LastMessage = utils.DEAL_STATUS_TRANSFER_STARTED //"transfer-started"
		deal.DTChan = channelId.String()
		i.LightNode.DB.Transaction(func(tx *gorm.DB) error {
			tx.Model(&core.PieceCommitment{}).Where("id = ?", pieceComm.ID).Save(pieceComm)
			tx.Model(&core.Content{}).Where("id = ?", content.ID).Save(content)
			tx.Model(&core.ContentDeal{}).Where("id = ?", deal.ID).Save(deal)
			return nil
		})

		// subscribe to data transfer events
		i.LightNode.Dispatcher.AddJob(NewDataTransferStatusListenerProcessor(i.LightNode))

	} else {

		// reprocess
		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm))
	}

	return nil

}

type MinerAddress struct {
	Address address.Address
}

func (i *StorageDealMakerProcessor) GetAssignedMinerForContent(content core.Content) MinerAddress {
	var storageMinerAssignment core.ContentMinerAssignment
	i.LightNode.DB.Model(&core.ContentMinerAssignment{}).Where("content = ?", content.ID).Find(&storageMinerAssignment)
	fmt.Println("storageMinerAssignment", storageMinerAssignment.ID)
	if storageMinerAssignment.ID != 0 {
		address.CurrentNetwork = address.Mainnet
		a, err := address.NewFromString(storageMinerAssignment.Miner)
		if err != nil {
			fmt.Println("error on miner address", err, a)
		}
		return MinerAddress{Address: a}
	}
	return i.GetStorageProviders()[0]
}

func (i *StorageDealMakerProcessor) GetStorageProviders() []MinerAddress {
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

func (i *StorageDealMakerProcessor) sendProposalV120(ctx context.Context, netprop network.Proposal, propCid cid.Cid, dealUUID uuid.UUID, dbid uint) (bool, error) {
	// Create an auth token to be used in the request
	authToken, err := httptransport.GenerateAuthToken()
	if err != nil {
		return false, xerrors.Errorf("generating auth token for deal: %w", err)
	}

	rootCid := netprop.Piece.Root
	size := netprop.Piece.RawBlockSize
	var announceAddr multiaddr.Multiaddr

	if len(i.LightNode.Node.Config.AnnounceAddrs) == 0 {
		return false, xerrors.Errorf("cannot serve deal data: no announce address configured for estuary node")
	}

	addrstr := i.LightNode.Node.Config.AnnounceAddrs[1] + "/p2p/" + i.LightNode.Node.Host.ID().String()
	announceAddr, err = multiaddr.NewMultiaddr(addrstr)
	if err != nil {
		return false, xerrors.Errorf("cannot parse announce address '%s': %w", addrstr, err)
	}

	// Add an auth token for the data to the auth DB
	err = i.LightNode.FilClient.Libp2pTransferMgr.PrepareForDataRequest(ctx, dbid, authToken, propCid, rootCid, size)
	if err != nil {
		return false, xerrors.Errorf("preparing for data request: %w", err)
	}

	// Send the deal proposal to the storage provider
	propPhase, err := i.LightNode.FilClient.SendProposalV120(ctx, dbid, netprop, dealUUID, announceAddr, authToken)

	if err != nil {
		i.LightNode.FilClient.Libp2pTransferMgr.CleanupPreparedRequest(i.Context, dbid, authToken)
		if strings.Contains(err.Error(), "deal proposal is identical") { // don't put it back on the queue
			return false, err
		}
	}
	return propPhase, err
}
