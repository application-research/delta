package jobs

import (
	"context"
	"delta/core"
	"delta/core/model"
	"delta/utils"
	"fmt"
	"github.com/ipfs/go-cid"
	"time"
)

type ItemContentCleanUpProcessor struct {
	Context   context.Context
	LightNode *core.DeltaNode
}

func NewItemContentCleanUpProcessor(ln *core.DeltaNode) IProcessor {
	return &ItemContentCleanUpProcessor{
		Context:   context.Background(),
		LightNode: ln,
	}
}

func (i ItemContentCleanUpProcessor) Run() error {

	// clear up finished CID deals.
	var contentsOnline []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status = ? and connection_mode = ?", "transfer-finished", "online").Find(&contentsOnline)

	for _, content := range contentsOnline {
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.Blockstore.DeleteBlock(i.Context, cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}

	var contentsOffline []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status = ? and connection_mode = ?", "deal-proposal-sent", "offline").Find(&contentsOffline)

	for _, content := range contentsOnline {
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.Blockstore.DeleteBlock(i.Context, cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}

	// clear up failed CID deals.
	var contentDeals []model.ContentDeal
	i.LightNode.DB.Model(&model.ContentDeal{}).Where("failed = ?", true).Find(&contentDeals)

	for _, contentDeal := range contentDeals {
		var content model.Content
		i.LightNode.DB.Model(&model.Content{}).Where("id = ?", contentDeal.Content).Find(&content)
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.Blockstore.DeleteBlock(i.Context, cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}

	// fail request that are older than 7 days. No point in retrying them.
	i.LightNode.DB.Model(&model.Content{}).Where("status <> ? and created_at < ?", "transfer-finished", time.Now().AddDate(0, 0, -7)).Updates(model.Content{
		Status:      utils.DEAL_STATUS_TRANSFER_FAILED,
		LastMessage: "Transfer failed. Record is older than 7 days.",
	})
	return nil
}
