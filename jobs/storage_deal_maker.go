package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	fc "github.com/application-research/filclient"
	smtypes "github.com/filecoin-project/boost/storagemarket/types"
	"github.com/filecoin-project/boost/transport/httptransport"
	boosttypes "github.com/filecoin-project/boost/transport/types"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-fil-markets/storagemarket/network"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/market"
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

// WalletMeta `WalletMeta` is a struct with two fields, `KeyType` and `PrivateKey`, both of which are strings.
// @property {string} KeyType - The type of key used to sign the transaction.
// @property {string} PrivateKey - The private key of the wallet.
type WalletMeta struct {
	KeyType    string `json:"key_type"`
	PrivateKey string `json:"private_key"`
}

// MinerAddress `MinerAddress` is a type that represents a miner address.
// @property Address - The address of the miner.
type MinerAddress struct {
	Address address.Address
}

// Genesis Miner
// this will never be used I want to keep it here to honor `Jason Cihelka` for his work.
var mainnetMinerStrs = []string{
	"f01963614",
}

// StorageDealMakerProcessor `StorageDealMakerProcessor` is a struct that contains a `context.Context`, a `*core.DeltaNode`, a `*model.Content`, and
// a `*model.PieceCommitment`.
// @property Context - The context of the deal maker processor.
// @property LightNode - The light node that is making the deal.
// @property Content - The content that the storage deal is for.
// @property PieceComm - The piece commitment that the miner is trying to prove.
type StorageDealMakerProcessor struct {
	// It creates a new `StorageDealMakerProcessor` object, which is a type of `IProcessor` object
	Context   context.Context
	LightNode *core.DeltaNode
	Content   *model.Content
	PieceComm *model.PieceCommitment
}

// NewStorageDealMakerProcessor It creates a new `StorageDealMakerProcessor` object, which is a type of `IProcessor` object
func NewStorageDealMakerProcessor(ln *core.DeltaNode, content model.Content, commitment model.PieceCommitment) IProcessor {
	return &StorageDealMakerProcessor{
		LightNode: ln,
		Content:   &content,
		PieceComm: &commitment,
		Context:   context.Background(),
	}
}

// Run The above code is a function that is part of the StorageDealMakerProcessor struct. It is a function that is called when
// the StorageDealMakerProcessor is run. It calls the makeStorageDeal function, which is defined in the same file.
func (i StorageDealMakerProcessor) Run() error {
	err := i.makeStorageDeal(i.Content, i.PieceComm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Making a deal with the miner.
func (i *StorageDealMakerProcessor) makeStorageDeal(content *model.Content, pieceComm *model.PieceCommitment) error {

	// update the status
	var contentToUpdate model.Content
	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Find(&contentToUpdate)
	contentToUpdate.Status = utils.CONTENT_DEAL_MAKING_PROPOSAL //"making-deal-proposal"
	contentToUpdate.UpdatedAt = time.Now()
	i.LightNode.DB.Save(&contentToUpdate)

	// any error here, fail the content
	var miner, errOnMinerAddr = i.GetAssignedMinerForContent(*content)
	if errOnMinerAddr != nil {
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		contentToUpdate.LastMessage = errOnMinerAddr.Error()
		contentToUpdate.UpdatedAt = time.Now()
		return errOnMinerAddr
	}
	minerAddress := miner.Address

	// get filclient instance for content.
	var filClient, errOnFilc = i.GetAssignedFilclientForContent(*content)

	if errOnFilc != nil {
		contentToUpdate.UpdatedAt = time.Now()
		contentToUpdate.LastMessage = errOnFilc.Error()
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		i.LightNode.DB.Save(&contentToUpdate)
		return errOnFilc
	}

	// prep the proposal
	var dealProposal, errOnDealPrep = i.GetDealProposalForContent(*content)

	if errOnDealPrep != nil {
		contentToUpdate.UpdatedAt = time.Now()
		contentToUpdate.LastMessage = errOnDealPrep.Error()
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		i.LightNode.DB.Save(&contentToUpdate)
		return errOnDealPrep
	}

	var priceBigInt types.BigInt
	if !dealProposal.VerifiedDeal {
		unverifiedDealPrice, errPrice := types.BigFromString(dealProposal.UnverifiedDealMaxPrice)
		if errPrice != nil {
			contentToUpdate.UpdatedAt = time.Now()
			contentToUpdate.LastMessage = errPrice.Error()
			contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
			i.LightNode.DB.Save(&contentToUpdate)
			return errPrice
		}
		bigIntBalance, errBalance := i.LightNode.LotusApiNode.WalletBalance(context.Background(), filClient.ClientAddr)
		if errBalance != nil {
			contentToUpdate.UpdatedAt = time.Now()
			contentToUpdate.LastMessage = errBalance.Error()
			contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
			i.LightNode.DB.Save(&contentToUpdate)
			return errBalance
		}
		// check if the balance is enough
		if unverifiedDealPrice.GreaterThan(bigIntBalance) {
			contentToUpdate.UpdatedAt = time.Now()
			contentToUpdate.LastMessage = "insufficient funds"
			contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
			contentToUpdate.AutoRetry = false                           // stop retrying if the balance is not enough. it won't work.
			i.LightNode.DB.Save(&contentToUpdate)
			return xerrors.New("insufficient funds")
		}

		priceBigInt = unverifiedDealPrice
	} else {
		verifiedPrice, errVerPrice := types.BigFromString("0")
		if errVerPrice != nil {
			contentToUpdate.UpdatedAt = time.Now()
			contentToUpdate.LastMessage = errVerPrice.Error()
			contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
			i.LightNode.DB.Save(&contentToUpdate)
			return errVerPrice
		}
		priceBigInt = verifiedPrice
	}

	var dealDuration = utils.DEFAULT_DURATION
	if dealProposal.ID != 0 {
		dealDuration = int(dealProposal.Duration)
	}
	duration := abi.ChainEpoch(dealDuration)
	payloadCid, err := cid.Decode(pieceComm.Cid)
	if err != nil {
		contentToUpdate.UpdatedAt = time.Now()
		contentToUpdate.LastMessage = err.Error()
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		i.LightNode.DB.Save(&contentToUpdate)
	}

	pieceCid, err := cid.Decode(pieceComm.Piece)
	if err != nil {
		contentToUpdate.UpdatedAt = time.Now()
		contentToUpdate.LastMessage = err.Error()
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		i.LightNode.DB.Save(&contentToUpdate)

	}

	// label deal
	label, err := market.NewLabelFromString(dealProposal.Label)
	if err != nil {
		contentToUpdate.UpdatedAt = time.Now()
		contentToUpdate.LastMessage = err.Error()
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		i.LightNode.DB.Save(&contentToUpdate)
	}

	prop, err := filClient.MakeDealWithOptions(i.Context, minerAddress, payloadCid, priceBigInt, duration,
		fc.DealWithVerified(dealProposal.VerifiedDeal),
		fc.DealWithFastRetrieval(!dealProposal.RemoveUnsealedCopy),
		fc.DealWithLabel(label),
		fc.DealWithStartEpoch(abi.ChainEpoch(dealProposal.StartEpoch)),
		fc.DealWithEndEpoch(abi.ChainEpoch(dealProposal.EndEpoch)),
		fc.DealWithPricePerEpoch(priceBigInt),
		fc.DealWithPieceInfo(fc.DealPieceInfo{
			Cid:         pieceCid,
			Size:        abi.PaddedPieceSize(pieceComm.PaddedPieceSize),
			PayloadSize: uint64(pieceComm.Size),
		}),
	)
	if err != nil {
		contentToUpdate.UpdatedAt = time.Now()
		contentToUpdate.LastMessage = err.Error()
		contentToUpdate.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED //"failed"
		i.LightNode.DB.Save(&contentToUpdate)
		return err
	}

	if err != nil {
		fmt.Println(err)
		switch {
		case strings.Contains(err.Error(), "miner connection failed: failed to dial"),
			strings.Contains(err.Error(), "opening stream to miner: failed to open stream to peer: protocol not supported"),
			strings.Contains(err.Error(), "error getting deal protocol for miner connecting"):
			if content.AutoRetry {
				// check the auto retry limit if it's reached then stop retrying
				var dealCount int64
				i.LightNode.DB.Model(&model.ContentDeal{}).Where("content = ?", content.ID).Count(&dealCount)
				if int(dealCount) >= i.LightNode.Config.Common.MaxAutoRetry {
					content.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED
					content.LastMessage = "Retry limit reached"
					content.AutoRetry = false
					content.UpdatedAt = time.Now()
					i.LightNode.DB.Save(&content)
					return nil
				}

				minerAssignService := core.NewMinerAssignmentService(*i.LightNode)
				provider, errOnPv := minerAssignService.GetSPWithGivenBytes(i.Content.Size)

				if errOnPv != nil {
					content.Status = utils.CONTENT_DEAL_PROPOSAL_FAILED
					content.LastMessage = err.Error()
					content.UpdatedAt = time.Now()
					i.LightNode.DB.Save(&content)
					return errOnPv
				}

				// create new content miner record.
				contentMiner := &model.ContentMiner{
					Content:   content.ID,
					Miner:     provider.Address,
					UpdatedAt: time.Now(),
					CreatedAt: time.Now(),
				}
				if err := i.LightNode.DB.Create(contentMiner).Error; err != nil {
					return xerrors.Errorf("failed to create database entry for content miner: %w", err)
				}

				// and dispatch the job again
				i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)
			}
		default:
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
				UpdatedAt:   time.Now(),
			})
		}
		return err
	}

	dealProp := prop.DealProposal
	prop.FastRetrieval = !dealProposal.RemoveUnsealedCopy

	if dealProposal.StartEpoch != 0 {
		dealProp.Proposal.StartEpoch = abi.ChainEpoch(dealProposal.StartEpoch)
	}

	if dealProposal.EndEpoch != 0 {
		dealProp.Proposal.EndEpoch = abi.ChainEpoch(dealProposal.EndEpoch)
	}
	propnd, err := cborutil.AsIpld(dealProp)

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
			UpdatedAt:   time.Now(),
		})
		return err
	}

	dealUUID := uuid.New()
	proto, err := filClient.DealProtocolForMiner(i.Context, minerAddress)
	if err != nil {
		switch {
		// there are cases where error occurs before the deal is even sent to the miner.
		case strings.Contains(err.Error(), "failed to send request: stream reset"),
			strings.Contains(err.Error(), "proposal piece size is invalid"),
			strings.Contains(err.Error(), "piece size less than minimum required size"),
			strings.Contains(err.Error(), "storage price per epoch less than asking price"),
			strings.Contains(err.Error(), "miner connection failed: failed to dial"),
			strings.Contains(err.Error(), "failed to dial"),
			strings.Contains(err.Error(), "deal proposal is identical to deal"),
			strings.Contains(err.Error(), "provider has insufficient funds to accept deal"),
			strings.Contains(err.Error(), "opening stream to miner: failed to open stream to peer: protocol not supported"),
			strings.Contains(err.Error(), "miner is not considering online storage deals"),
			strings.Contains(err.Error(), "miner is not accepting unverified storage deals"),
			strings.Contains(err.Error(), "Deal rejected | Under maintenance, retry later"),
			strings.Contains(err.Error(), "Deal rejected | Price below acceptance for such deal"),
			strings.Contains(err.Error(), "Deal rejected | Such deal is not accepted"),
			strings.Contains(err.Error(), "Deal rejected | Error | like a spam"),
			strings.Contains(err.Error(), "does not support any deal making protocol"),
			strings.Contains(err.Error(), "failed validation: server error: getting chain head"),
			strings.Contains(err.Error(), "Error 2 (Worker balance too low)"),
			strings.Contains(err.Error(), "send proposal rpc:"):
			if content.AutoRetry {
				var dealCount int64
				i.LightNode.DB.Model(&model.ContentDeal{}).Where("content = ?", content.ID).Count(&dealCount)
				if int(dealCount) >= i.LightNode.Config.Common.MaxAutoRetry {
					i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
						Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
						LastMessage: "Retry limit reached",
						AutoRetry:   false,
						UpdatedAt:   time.Now(),
					})
					return nil
				}
				minerAssignService := core.NewMinerAssignmentService(*i.LightNode)
				provider, errOnPv := minerAssignService.GetSPWithGivenBytes(content.Size)
				if errOnPv != nil {
					// just fail it then
					i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
						Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
						LastMessage: err.Error(),
						UpdatedAt:   time.Now(),
					})
					return errOnPv
				}

				// create new content miner record.
				contentMiner := &model.ContentMiner{
					Content:   content.ID,
					Miner:     provider.Address,
					UpdatedAt: time.Now(),
					CreatedAt: time.Now(),
				}
				if err := i.LightNode.DB.Create(contentMiner).Error; err != nil {
					return xerrors.Errorf("failed to create database entry for content miner: %w", err)
				}

				// and dispatch the job again
				i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)
			}
		default:
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
				Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
				LastMessage: err.Error(),
				UpdatedAt:   time.Now(),
			})
		}
		return err
	}
	deal := &model.ContentDeal{
		Content:             content.ID,
		PropCid:             propnd.Cid().String(),
		DealUUID:            dealUUID.String(),
		Miner:               minerAddress.String(),
		Verified:            true,
		DealProtocolVersion: string(proto),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	if err := i.LightNode.DB.Create(deal).Error; err != nil {
		i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)
		return xerrors.Errorf("failed to create database entry for deal: %w", err)
	}

	propString := propnd.String()

	// 	log and send the proposal over
	i.LightNode.DB.Create(&model.ContentDealProposal{
		Content:   content.ID,
		Meta:      propString,
		Signed:    propString,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
		Status: utils.CONTENT_DEAL_SENDING_PROPOSAL, //"sending-deal-proposal",
	})

	// send the proposal.
	_, errProp := i.sendProposalV120(i.Context, *prop, propnd.Cid(), dealUUID, uint(deal.ID), dealProposal)

	// check all errors
	if errProp != nil {
		contentToUpdate = model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED,
			LastMessage: errProp.Error(),
			UpdatedAt:   time.Now(),
		}
		contentDealToUpdate := model.ContentDeal{
			LastMessage: errProp.Error(),
			Failed:      true,
			UpdatedAt:   time.Now(),
		}
		switch {
		// we only retry if the error is one of these
		case strings.Contains(errProp.Error(), "failed to send request: stream reset"),
			strings.Contains(errProp.Error(), "proposal piece size is invalid"),
			strings.Contains(errProp.Error(), "piece size less than minimum required size"),
			strings.Contains(errProp.Error(), "storage price per epoch less than asking price"),
			strings.Contains(errProp.Error(), "miner connection failed: failed to dial"),
			strings.Contains(errProp.Error(), "failed to dial"),
			strings.Contains(errProp.Error(), "connection limited. rate: Wait(n=3) would exceed context deadline"),
			strings.Contains(errProp.Error(), "deal proposal is identical to deal"),
			strings.Contains(errProp.Error(), "provider has insufficient funds to accept deal"),
			strings.Contains(errProp.Error(), "opening stream to miner: failed to open stream to peer: protocol not supported"),
			strings.Contains(errProp.Error(), "miner is not considering online storage deals"),
			strings.Contains(errProp.Error(), "miner is not accepting unverified storage deals"),
			strings.Contains(errProp.Error(), "Deal rejected | Under maintenance, retry later"),
			strings.Contains(errProp.Error(), "Deal rejected | Price below acceptance for such deal"),
			strings.Contains(errProp.Error(), "Deal rejected | Such deal is not accepted"),
			strings.Contains(errProp.Error(), "Deal rejected | Error | like a spam"),
			strings.Contains(errProp.Error(), "failed validation: server error: getting chain head"),
			strings.Contains(errProp.Error(), "Error 2 (Worker balance too low)"),
			strings.Contains(errProp.Error(), "send proposal rpc:"):

			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(&contentDealToUpdate)
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(&contentToUpdate)

			// re-assign a miner
			if content.AutoRetry {
				// check the auto retry limit if it's reached then stop retrying
				var dealCount int64
				i.LightNode.DB.Model(&model.ContentDeal{}).Where("content = ?", content.ID).Count(&dealCount)
				if int(dealCount) >= i.LightNode.Config.Common.MaxAutoRetry {
					i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
						Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
						LastMessage: "Retry limit reached",
						AutoRetry:   false,
						UpdatedAt:   time.Now(),
					})
					return nil
				}
				minerAssignService := core.NewMinerAssignmentService(*i.LightNode)
				provider, errOnPv := minerAssignService.GetSPWithGivenBytes(content.Size)
				if errOnPv != nil {
					// just fail it then
					i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(&contentDealToUpdate)
					i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(&contentToUpdate)
					return errOnPv
				}

				// create new content miner record.
				contentMiner := &model.ContentMiner{
					Content:   content.ID,
					Miner:     provider.Address,
					UpdatedAt: time.Now(),
					CreatedAt: time.Now(),
				}
				if err := i.LightNode.DB.Create(contentMiner).Error; err != nil {
					return xerrors.Errorf("failed to create database entry for content miner: %w", err)
				}

				// and dispatch the job again
				i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)
			}
			return errProp
		default:
			fmt.Println("default", errProp.Error())
			i.LightNode.DB.Model(&deal).Where("id = ?", deal.ID).Updates(&contentDealToUpdate)
			i.LightNode.DB.Model(&content).Where("id = ?", content.ID).
				Updates(model.Content{
					Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED,
					LastMessage: errProp.Error(),
					UpdatedAt:   time.Now(),
				})
		}
		return errProp
	}

	// if this is e2e, then we need to start the data transfer.
	if errProp == nil && content.ConnectionMode == utils.CONNECTION_MODE_E2E {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status: utils.CONTENT_DEAL_PROPOSAL_SENT, //"sending-deal-proposal",
		})
		propCid, err := cid.Decode(deal.PropCid)
		contentCid, err := cid.Decode(content.Cid)

		assignedMinerForContent, err := i.GetAssignedMinerForContent(*content)
		if err != nil {
			return err
		}
		channelId, err := filClient.StartDataTransfer(i.Context, assignedMinerForContent.Address, propCid, contentCid)

		// if this is online then the user/sp expects the data to be transferred. if it fails, re-try.
		if err != nil {
			i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, *content, *pieceComm), 1)
			return err
		}

		content.PieceCommitmentId = pieceComm.ID
		pieceComm.Status = utils.COMMP_STATUS_COMITTED           //"committed"
		content.Status = utils.DEAL_STATUS_TRANSFER_STARTED      //"transfer-started"
		content.LastMessage = utils.DEAL_STATUS_TRANSFER_STARTED //"transfer-started"
		deal.LastMessage = utils.DEAL_STATUS_TRANSFER_STARTED    //"transfer-started"

		pieceComm.UpdatedAt = time.Now()
		content.UpdatedAt = time.Now()
		deal.UpdatedAt = time.Now()
		deal.TransferStarted = time.Now()

		deal.DTChan = channelId.String()
		i.LightNode.DB.Transaction(func(tx *gorm.DB) error {
			tx.Model(&model.PieceCommitment{}).Where("id = ?", pieceComm.ID).Save(pieceComm)
			tx.Model(&model.Content{}).Where("id = ?", content.ID).Save(content)
			tx.Model(&model.ContentDeal{}).Where("id = ?", deal.ID).Save(deal)
			return nil
		})

	}

	//	if this is import, then we need to mark the deal as deal_proposal_sent.
	if errProp == nil && content.ConnectionMode == utils.CONNECTION_MODE_IMPORT {
		pieceComm.Status = utils.COMMP_STATUS_COMITTED //"committed"
		content.Status = utils.CONTENT_DEAL_PROPOSAL_SENT
		deal.LastMessage = utils.CONTENT_DEAL_PROPOSAL_SENT

		pieceComm.UpdatedAt = time.Now()
		content.UpdatedAt = time.Now()
		deal.UpdatedAt = time.Now()

		i.LightNode.DB.Transaction(func(tx *gorm.DB) error {
			tx.Model(&pieceComm).Where("id = ?", pieceComm.ID).Save(pieceComm)
			tx.Model(&content).Where("id = ?", content.ID).Save(content)
			tx.Model(&deal).Where("id = ?", deal.ID).Save(deal)
			return nil
		})
	}

	return nil

}

// GetAssignedMinerForContent Getting the miner address for the content.
func (i *StorageDealMakerProcessor) GetAssignedMinerForContent(content model.Content) (MinerAddress, error) {
	var storageMinerAssignment model.ContentMiner
	err := i.LightNode.DB.Model(&model.ContentMiner{}).Where("content = ?", content.ID).Order("created_at desc").First(&storageMinerAssignment).Error
	if err != nil {
		return MinerAddress{}, err
	}
	if storageMinerAssignment.ID != 0 {
		address.CurrentNetwork = address.Mainnet
		a, err := address.NewFromString(storageMinerAssignment.Miner)
		if err != nil {
			fmt.Println("error on miner address", err, a)
		}
		return MinerAddress{Address: a}, nil
	}
	return i.GetStorageProviders()[0], nil
}

// Getting the content deal proposal parameters for a given content.
func (i *StorageDealMakerProcessor) GetDealProposalForContent(content model.Content) (model.ContentDealProposalParameters, error) {
	var contentDealProposalParameters model.ContentDealProposalParameters
	err := i.LightNode.DB.Model(&model.ContentDealProposalParameters{}).Where("content = ?", content.ID).Order("created_at desc").First(&contentDealProposalParameters).Error
	if err != nil {
		return model.ContentDealProposalParameters{}, err
	}
	return contentDealProposalParameters, err
}

// Creating a new filclient for the content.
func (i *StorageDealMakerProcessor) GetAssignedFilclientForContent(content model.Content) (*fc.FilClient, error) {
	api := i.LightNode.LotusApiNode
	var storageWalletAssignment model.ContentWallet
	i.LightNode.DB.Model(&model.ContentWallet{}).Where("content = ?", content.ID).Find(&storageWalletAssignment)

	if storageWalletAssignment.ID != 0 {
		newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())
		if err != nil {
			fmt.Println("error on new wallet", err)
			return nil, err
		}

		// get the wallet entry
		var wallet model.Wallet
		i.LightNode.DB.Model(&model.Wallet{}).Where("id = ?", storageWalletAssignment.WalletId).Find(&wallet)
		decodedPkey, err := base64.StdEncoding.DecodeString(wallet.PrivateKey)
		if err != nil {
			fmt.Println("error on base64 decode", err)
			return nil, err
		}

		newWalletAddr, err := newWallet.WalletImport(context.Background(), &types.KeyInfo{
			Type:       types.KeyType(wallet.KeyType),
			PrivateKey: decodedPkey,
		})

		if err != nil {
			fmt.Println("error on wallet_estuary import", err)
			return nil, err
		}
		// new filclient just for this request
		filclient, err := fc.NewClient(i.LightNode.Node.Host, api, newWallet, newWalletAddr, i.LightNode.Node.Blockstore, i.LightNode.Node.Datastore, i.LightNode.Node.Config.DatastoreDir.Directory)
		if err != nil {
			fmt.Println("error on filclient", err)
			return nil, err
		}
		core.SetLibp2pManagerSubscribe(i.LightNode)
		return filclient, err
	}

	return i.LightNode.FilClient, nil
}

// GetStorageProviders Getting the storage providers.
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

// Sending a proposal to the peer.
func (i *StorageDealMakerProcessor) sendProposalV120(ctx context.Context,
	netprop network.Proposal,
	propCid cid.Cid,
	dealUUID uuid.UUID,
	dbid uint,
	dealProposal model.ContentDealProposalParameters) (bool, error) {

	// Create an auth token to be used in the request
	authToken, err := httptransport.GenerateAuthToken()

	if err != nil {
		return false, xerrors.Errorf("generating auth token for deal: %w", err)
	}

	netprop.FastRetrieval = !dealProposal.RemoveUnsealedCopy
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

	var transferParamsBoost boosttypes.HttpRequest
	json.Unmarshal([]byte(dealProposal.TransferParams), &transferParamsBoost)

	transferParams, err := json.Marshal(boosttypes.HttpRequest{
		URL: transferParamsBoost.URL,
		Headers: map[string]string{
			"Authorization": httptransport.BasicAuthHeader("", authToken),
		},
	})

	// Add an auth token for the data to the auth DB
	err = i.LightNode.FilClient.Libp2pTransferMgr.PrepareForDataRequest(ctx, dbid, authToken, propCid, rootCid, size)
	if err != nil {
		return false, xerrors.Errorf("preparing for data request: %w", err)
	}

	// Send the deal proposal to the storage provider
	var propPhase bool
	//var err error
	if i.Content.ConnectionMode == utils.CONNECTION_MODE_IMPORT && transferParamsBoost.URL == "" {
		propPhase, err = i.LightNode.FilClient.SendProposalV120WithOptions(
			ctx, netprop,
			fc.ProposalV120WithDealUUID(dealUUID),
			fc.ProposalV120WithLibp2pTransfer(announceAddr, authToken, dbid),
			fc.ProposalV120WithOffline(true),
			fc.ProposalV120WithSkipIPNIAnnounce(dealProposal.SkipIPNIAnnounce),
			fc.ProposalV120WithTransfer(smtypes.Transfer{
				Type: func() string {
					// starts with http
					if strings.Contains(dealProposal.TransferParams, "http") {
						return "http"
					}
					if strings.Contains(dealProposal.TransferParams, "https") {
						return "https"
					}
					if strings.Contains(dealProposal.TransferParams, "ftp") {
						return "ftp"
					}
					if strings.Contains(dealProposal.TransferParams, "ftps") {
						return "ftps"
					}
					return "libp2p"
				}(),
				ClientID: fmt.Sprintf("%d", dbid),
				Params:   transferParams,
				Size:     netprop.Piece.RawBlockSize,
			}),
		)
	} else {
		propPhase, err = i.LightNode.FilClient.SendProposalV120WithOptions(
			ctx, netprop,
			fc.ProposalV120WithDealUUID(dealUUID),
			fc.ProposalV120WithLibp2pTransfer(announceAddr, authToken, dbid),
			fc.ProposalV120WithSkipIPNIAnnounce(dealProposal.SkipIPNIAnnounce),
			fc.ProposalV120WithTransfer(smtypes.Transfer{
				Type: func() string {
					// starts with http
					if strings.Contains(dealProposal.TransferParams, "http") {
						return "http"
					}
					if strings.Contains(dealProposal.TransferParams, "https") {
						return "https"
					}
					if strings.Contains(dealProposal.TransferParams, "ftp") {
						return "ftp"
					}
					if strings.Contains(dealProposal.TransferParams, "ftps") {
						return "ftps"
					}
					return "libp2p"
				}(),
				ClientID: fmt.Sprintf("%d", dbid),
				Params:   transferParams,
				Size:     netprop.Piece.RawBlockSize,
			}),
		)
	}

	if err != nil {
		i.LightNode.FilClient.Libp2pTransferMgr.CleanupPreparedRequest(i.Context, dbid, authToken)
	}

	return propPhase, err
}
