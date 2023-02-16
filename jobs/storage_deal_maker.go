package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	fc "github.com/application-research/filclient"
	smtypes "github.com/filecoin-project/boost/storagemarket/types"
	"github.com/filecoin-project/boost/transport/httptransport"
	boosttypes "github.com/filecoin-project/boost/transport/types"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
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

	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
		Status: utils.CONTENT_DEAL_MAKING_PROPOSAL, //"making-deal-proposal",
	})

	// any error here, fail the content
	var minerAddress = i.GetAssignedMinerForContent(*content).Address
	var filClient, err = i.GetAssignedWalletForContent(*content)
	var dealProposal = i.GetDealProposalForContent(*content)

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
	}

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
	}

	priceBigInt, err := types.BigFromString("0")

	var dealDuration = utils.DEFAULT_DURATION
	if dealProposal.ID != 0 {
		dealDuration = int(dealProposal.Duration)
	}
	duration := abi.ChainEpoch(dealDuration)
	payloadCid, err := cid.Decode(pieceComm.Cid)
	pieceCid, err := cid.Decode(pieceComm.Piece)

	prop, err := filClient.MakeDealWithOptions(i.Context, minerAddress, payloadCid, priceBigInt, duration,
		fc.DealWithVerified(true),
		fc.DealWithFastRetrieval(false), // this should be a parameter.
		fc.DealWithPieceInfo(fc.DealPieceInfo{
			Cid:         pieceCid,
			Size:        abi.PaddedPieceSize(pieceComm.PaddedPieceSize),
			PayloadSize: uint64(pieceComm.Size),
		}),
	)

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
	}

	dealProp := prop.DealProposal
	if dealProposal.StartEpoch != 0 {
		dealProp.Proposal.StartEpoch = abi.ChainEpoch(dealProposal.StartEpoch)
		dealProp.Proposal.EndEpoch = dealProp.Proposal.StartEpoch + duration
	}

	propnd, err := cborutil.AsIpld(dealProp)
	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
	}

	dealUUID := uuid.New()
	proto, err := filClient.DealProtocolForMiner(i.Context, minerAddress)
	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
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
		i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)
		return xerrors.Errorf("failed to create database entry for deal: %w", err)
	}

	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
		Status: utils.CONTENT_DEAL_SENDING_PROPOSAL, //"sending-deal-proposal",
	})

	// 	send the proposal over
	propPhase, err := i.sendProposalV120(i.Context, *prop, propnd.Cid(), dealUUID, uint(deal.ID))

	if propPhase == true && err != nil {

		if strings.Contains(err.Error(), "deal proposal is identical") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})

			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}

		if strings.Contains(err.Error(), "deal duration out of bounds") {
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})

			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}

		if strings.Contains(err.Error(), "storage price per epoch less than asking price") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}
		if strings.Contains(err.Error(), " piece size less than minimum required size") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}

		if strings.Contains(err.Error(), " invalid deal end epoch") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}

		if strings.Contains(err.Error(), "could not load link") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}

		if strings.Contains(err.Error(), "proposal PieceCID had wrong prefix") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}

		if strings.Contains(err.Error(), "proposal piece size is invalid") { // don't put it back on the queue
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(core.ContentDeal{
				LastMessage: err.Error(),
				Failed:      true, // mark it as failed
			})
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
			})
			return err
		}
		i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)

		return err
	}

	if propPhase == false && err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(core.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_SENT, //"deal-proposal-sent",
			LastMessage: err.Error(),
		})
	}

	// Online - transfer it. Offline, proposal is enough
	if propPhase == false && content.ConnectionMode == "online" {

		propCid, err := cid.Decode(deal.PropCid)
		contentCid, err := cid.Decode(content.Cid)
		channelId, err := filClient.StartDataTransfer(i.Context, i.GetAssignedMinerForContent(*content).Address, propCid, contentCid)

		// if this is online then the user/sp expects the data to be transferred. if it fails, re-try.
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
		i.LightNode.Dispatcher.AddJobAndDispatch(NewDataTransferStatusListenerProcessor(i.LightNode), 1)
	}

	return nil

}

func (i *StorageDealMakerProcessor) CatchFailures() {

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

type WalletMeta struct {
	KeyType    string `json:"key_type"`
	PrivateKey string `json:"private_key"`
}

func (i *StorageDealMakerProcessor) GetDealProposalForContent(content core.Content) core.ContentDealProposalParameters {
	var contentDealProposalParameters core.ContentDealProposalParameters
	i.LightNode.DB.Model(&core.ContentDealProposalParameters{}).Where("content = ?", content.ID).Find(&contentDealProposalParameters)

	return contentDealProposalParameters
}

func (i *StorageDealMakerProcessor) GetAssignedWalletForContent(content core.Content) (*fc.FilClient, error) {
	api, _, err := core.LotusConnection(utils.LOTUS_API)
	if err != nil {
		return nil, err
	}

	var storageWalletAssignment core.ContentWalletAssignment
	i.LightNode.DB.Model(&core.ContentWalletAssignment{}).Where("content = ?", content.ID).Find(&storageWalletAssignment)

	if storageWalletAssignment.ID != 0 {
		newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())

		var walletMeta WalletMeta

		json.Unmarshal([]byte(storageWalletAssignment.Wallet), &walletMeta)
		unhexPkey, err := hex.DecodeString(walletMeta.PrivateKey)
		decodedPkey, err := base64.StdEncoding.DecodeString(string(unhexPkey))

		if err != nil {
			fmt.Println("error on unhex", err)
			return nil, err
		}

		newWalletAddr, err := newWallet.WalletImport(context.Background(), &types.KeyInfo{
			Type:       types.KeyType(walletMeta.KeyType),
			PrivateKey: decodedPkey,
		})

		if err != nil {
			fmt.Println("error on wallet import", err)
			return nil, err
		}
		// new filclient just for this request
		fc, err := fc.NewClient(i.LightNode.Node.Host, api, newWallet, newWalletAddr, i.LightNode.Node.Blockstore, i.LightNode.Node.Datastore, i.LightNode.Node.Config.DatastoreDir.Directory)
		if err != nil {
			fmt.Println("error on filclient", err)
			return nil, err
		}
		return fc, err
	}

	return i.LightNode.FilClient, err
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
	transferParams, err := json.Marshal(boosttypes.HttpRequest{
		URL: "libp2p://" + announceAddr.String(),
		Headers: map[string]string{
			"Authorization": httptransport.BasicAuthHeader("", authToken),
		},
	})

	// Send the deal proposal to the storage provider
	//propPhase, err := i.LightNode.FilClient.SendProposalV120(ctx, dbid, netprop, dealUUID, announceAddr, authToken)
	//if err != nil {
	//	i.LightNode.FilClient.Libp2pTransferMgr.CleanupPreparedRequest(i.Context, dbid, authToken)
	//	if strings.Contains(err.Error(), "deal proposal is identical") { // don't put it back on the queue
	//		return false, err
	//	}
	//}
	var propPhase bool
	//var err error
	if i.Content.ConnectionMode == "offline" {
		propPhase, err = i.LightNode.FilClient.SendProposalV120WithOptions(
			ctx, netprop,
			fc.ProposalV120WithDealUUID(dealUUID),
			fc.ProposalV120WithLibp2pTransfer(announceAddr, authToken, dbid),
			fc.ProposalV120WithOffline(true),
			fc.ProposalV120WithTransfer(smtypes.Transfer{
				Type:     "libp2p",
				ClientID: fmt.Sprintf("%d", dbid),
				Params:   transferParams,
				Size:     uint64(i.Content.Size),
			}),
		)
	} else {
		propPhase, err = i.LightNode.FilClient.SendProposalV120WithOptions(
			ctx, netprop,
			fc.ProposalV120WithDealUUID(dealUUID),
			fc.ProposalV120WithLibp2pTransfer(announceAddr, authToken, dbid),
			fc.ProposalV120WithSkipIPNIAnnounce(false),
			fc.ProposalV120WithTransfer(smtypes.Transfer{
				Type:     "libp2p",
				ClientID: fmt.Sprintf("%d", dbid),
				Params:   transferParams,
				Size:     netprop.Piece.RawBlockSize,
			}),
		)
	}

	return propPhase, err
}
