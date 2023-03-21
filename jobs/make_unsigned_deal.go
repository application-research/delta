package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	fc "github.com/application-research/filclient"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v9/market"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/ipfs/go-cid"
	"time"
)

type UnsignedDealMakerProcessor struct {
	// It creates a new `StorageDealMakerProcessor` object, which is a type of `IProcessor` object
	Context   context.Context
	LightNode *core.DeltaNode
	Content   *model.Content
	PieceComm *model.PieceCommitment
}

func NewUnsignedDealMakerProcessor(ln *core.DeltaNode, content model.Content, commitment model.PieceCommitment) IProcessor {
	return &UnsignedDealMakerProcessor{
		LightNode: ln,
		Content:   &content,
		PieceComm: &commitment,
		Context:   context.Background(),
	}
}

// Run The above code is a function that is part of the StorageDealMakerProcessor struct. It is a function that is called when
// the StorageDealMakerProcessor is run. It calls the makeStorageDeal function, which is defined in the same file.
func (i UnsignedDealMakerProcessor) Run() error {
	err := i.makeUnsignedStorageDeal(i.Content, i.PieceComm)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Making a deal with the miner.
func (i *UnsignedDealMakerProcessor) makeUnsignedStorageDeal(content *model.Content, pieceComm *model.PieceCommitment) error {

	i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
		Status: utils.CONTENT_DEAL_MAKING_PROPOSAL, //"making-deal-proposal",
	})

	// any error here, fail the content
	var minerAddress = i.GetAssignedMinerForContent(*content).Address
	var filClient, err = i.GetAssignedFilclientForContent(*content)
	var dealProposal = i.GetDealProposalForContent(*content)

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
	}

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
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
	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
	}

	pieceCid, err := cid.Decode(pieceComm.Piece)
	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
	}

	// label deal
	label, err := market.NewLabelFromString(dealProposal.Label)
	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
	}

	prop, rawUnsigned, err := filClient.MakeDealUnsigned(i.Context, minerAddress, payloadCid, priceBigInt, duration,
		fc.DealWithVerified(true),
		fc.DealWithFastRetrieval(!dealProposal.RemoveUnsealedCopy),
		fc.DealWithLabel(label),
		fc.DealWithPieceInfo(fc.DealPieceInfo{
			Cid:         pieceCid,
			Size:        abi.PaddedPieceSize(pieceComm.PaddedPieceSize),
			PayloadSize: uint64(pieceComm.Size),
		}),
	)

	if err != nil {
		i.LightNode.DB.Model(&content).Where("id = ?", content.ID).Updates(model.Content{
			Status:      utils.CONTENT_DEAL_PROPOSAL_FAILED, //"failed",
			LastMessage: err.Error(),
		})
		return err
	}

	// struct to json
	encodedDealProposal, err := json.Marshal(prop)
	encodedUnsignedDealProposal := hex.EncodeToString(rawUnsigned)
	// save the unsigned deal
	var unsignedDeal = model.ContentDealProposal{
		Content:   content.ID,
		Unsigned:  encodedUnsignedDealProposal,
		Meta:      string(encodedDealProposal),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	i.LightNode.DB.Create(&unsignedDeal)
	return nil

}

// GetAssignedMinerForContent Getting the miner address for the content.
func (i *UnsignedDealMakerProcessor) GetAssignedMinerForContent(content model.Content) MinerAddress {
	var storageMinerAssignment model.ContentMiner
	i.LightNode.DB.Model(&model.ContentMiner{}).Where("content = ?", content.ID).Find(&storageMinerAssignment)
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

// Getting the content deal proposal parameters for a given content.
func (i *UnsignedDealMakerProcessor) GetDealProposalForContent(content model.Content) model.ContentDealProposalParameters {
	var contentDealProposalParameters model.ContentDealProposalParameters
	i.LightNode.DB.Model(&model.ContentDealProposalParameters{}).Where("content = ?", content.ID).Find(&contentDealProposalParameters)
	return contentDealProposalParameters
}

// Getting the assigned filclient for the content.
func (i *UnsignedDealMakerProcessor) GetAssignedFilclientForContent(content model.Content) (*fc.FilClient, error) {
	api, _, err := core.LotusConnection(utils.LOTUS_API)
	if err != nil {
		return nil, err
	}

	var storageWalletAssignment model.ContentWallet
	i.LightNode.DB.Model(&model.ContentWallet{}).Where("content = ?", content.ID).Find(&storageWalletAssignment)

	if storageWalletAssignment.ID != 0 {
		newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())

		var wallet model.Wallet
		i.LightNode.DB.Model(&model.Wallet{}).Where("id = ?", storageWalletAssignment.WalletId).Find(&wallet)
		decodedPkey, err := base64.StdEncoding.DecodeString(wallet.PrivateKey)

		if err != nil {
			fmt.Println("error on unhex", err)
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

	return i.LightNode.FilClient, err
}

// GetStorageProviders Getting the storage providers.
func (i *UnsignedDealMakerProcessor) GetStorageProviders() []MinerAddress {
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
