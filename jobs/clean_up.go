package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/ipfs/go-cid"
	"time"
)

// ContentCleanUpProcessor `ContentCleanUpProcessor` is a struct that has a `Context` and a `LightNode` field.
// @property Context - The context of the current request.
// @property LightNode - This is the node that we want to clean up.
type ContentCleanUpProcessor struct {
	Context   context.Context
	LightNode *core.DeltaNode
}

// NewItemContentCleanUpProcessor `NewItemContentCleanUpProcessor` creates a new `ContentCleanUpProcessor` instance
func NewContentCleanUpProcessor(ln *core.DeltaNode) IProcessor {
	return &ContentCleanUpProcessor{
		Context:   context.Background(),
		LightNode: ln,
	}
}

// Run Cleaning up the database.
func (i ContentCleanUpProcessor) Run() error {

	// clear up finished CID deals.
	var contentsOnline []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status = ? and connection_mode = ?", "transfer-finished", "e2e").Find(&contentsOnline)

	for _, content := range contentsOnline {
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}

	var contentsOffline []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status = ? and connection_mode = ?", "deal-proposal-sent", "import").Find(&contentsOffline)

	for _, content := range contentsOffline {
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidD)
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
		err = i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}

	// clear up cids that are older than 3 days.
	var oldContents []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status not in(?,?,?,?) and created_at < ?", "transfer-failed", "deal-proposal-failed", "transfer-finished", "deal-proposal-sent", time.Now().AddDate(0, 0, -3)).Find(&oldContents)

	for _, content := range oldContents {

		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}

		content.Status = utils.DEAL_STATUS_TRANSFER_FAILED
		content.LastMessage = "Transfer failed. Record is older than 3 days."
		i.LightNode.DB.Save(&content)

	}
	return nil
}
