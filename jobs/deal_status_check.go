package jobs

import (
	"context"
	"delta/core"
	"encoding/base64"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	fc "github.com/application-research/filclient"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/wallet"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
)

type DealStatusCheck struct {
	LightNode *core.DeltaNode
	Content   *model.Content
}

func (d DealStatusCheck) Run() error {
	var contentDeals []model.ContentDeal
	// get the latest content deal of the content
	d.LightNode.DB.Where("content = ?", d.Content.ID).Order("created_at desc").Find(&contentDeals)

	filcOfContent, errFilc := d.GetAssignedFilclientForContent(*d.Content)
	if errFilc != nil {
		return errFilc
	}

	for _, contentDeal := range contentDeals {

		if contentDeal.DealUUID == "" {
			return nil
		}

		miner, err := address.NewFromString(contentDeal.Miner)
		if err != nil {
			return err
		}

		cidProp, err := cid.Decode(contentDeal.PropCid)
		if err != nil {
			return err

		}
		dealUuid, err := uuid.Parse(contentDeal.DealUUID)
		if err != nil {
			return err
		}

		// get the status
		status, err := filcOfContent.DealStatus(context.Background(), miner, cidProp, &dealUuid)
		if err != nil {
			return err
		}
		contentDeal.DealID = int64(status.DealID)

		if status.State != storagemarket.StorageDealUnknown {
			d.Content.Status = storagemarket.DealStates[status.State]
			d.Content.LastMessage = storagemarket.DealStatesDescriptions[status.State]
			contentDeal.LastMessage = storagemarket.DealStatesDescriptions[status.State]
			contentDeal.DealID = int64(status.DealID)
		}
		d.LightNode.DB.Save(&contentDeal)
		d.LightNode.DB.Save(&d.Content)
	}
	return nil
}

func NewDealStatusCheck(ln *core.DeltaNode, content *model.Content) IProcessor {
	return &DealStatusCheck{
		LightNode: ln,
		Content:   content,
	}
}

func (d DealStatusCheck) GetAssignedFilclientForContent(content model.Content) (*fc.FilClient, error) {
	api := d.LightNode.LotusApiNode
	var storageWalletAssignment model.ContentWallet
	d.LightNode.DB.Model(&model.ContentWallet{}).Where("content = ?", content.ID).Find(&storageWalletAssignment)

	if storageWalletAssignment.ID != 0 {
		newWallet, err := wallet.NewWallet(wallet.NewMemKeyStore())
		if err != nil {
			fmt.Println("error on new wallet", err)
			return nil, err
		}

		// get the wallet entry
		var wallet model.Wallet
		d.LightNode.DB.Model(&model.Wallet{}).Where("id = ?", storageWalletAssignment.WalletId).Find(&wallet)
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
		filclient, err := fc.NewClient(d.LightNode.Node.Host, api, newWallet, newWalletAddr, d.LightNode.Node.Blockstore, d.LightNode.Node.Datastore, d.LightNode.Node.Config.DatastoreDir.Directory)
		if err != nil {
			fmt.Println("error on filclient", err)
			return nil, err
		}
		core.SetLibp2pManagerSubscribe(d.LightNode)
		return filclient, err
	}

	return d.LightNode.FilClient, nil
}
