package jobs

import (
	"context"
	"delta/core"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/ipfs/go-cid"
)

// ItemContentCleanUpProcessor ContentCleanUpProcessor `ContentCleanUpProcessor` is a struct that has a `Context` and a `LightNode` field.
// @property Context - The context of the current request.
// @property LightNode - This is the node that we want to clean up.
type ItemContentCleanUpProcessor struct {
	Content   model.Content
	LightNode *core.DeltaNode
}

// NewItemContentCleanUpProcessor `NewItemContentCleanUpProcessor` creates a new `ContentCleanUpProcessor` instance
func NewItemContentCleanUpProcessor(ln *core.DeltaNode, content model.Content) IProcessor {
	return &ItemContentCleanUpProcessor{
		Content:   content,
		LightNode: ln,
	}
}

// Run Cleaning up the database.
func (i ItemContentCleanUpProcessor) Run() error {

	// clear up finished CID deals.
	var contentsOnline []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status = ? and connection_mode = ? and id = ?", "transfer-finished", "e2e", i.Content.ID).Find(&contentsOnline)

	for _, content := range contentsOnline {
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.DAGService.Remove(context.Background(), cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}

	// clear up failed CID deals.
	var contentDeals []model.ContentDeal
	i.LightNode.DB.Model(&model.ContentDeal{}).Where("failed = ? and content = ?", true, i.Content.ID).Find(&contentDeals)

	for _, contentDeal := range contentDeals {
		var content model.Content
		i.LightNode.DB.Model(&model.Content{}).Where("id = ?", contentDeal.Content).Find(&content)
		cidD, err := cid.Decode(content.Cid)
		if err != nil {
			fmt.Println("error in decoding cid", err)
			continue
		}
		err = i.LightNode.Node.DAGService.Remove(context.Background(), cidD)
		if err != nil {
			fmt.Println("error in deleting block", err)
			continue
		}
	}
	return nil
}
