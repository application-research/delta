package jobs

import (
	"context"
	"delta/core"
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

	// get all content with transfer status = "transfer-finished"
	var contents []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ?", "transfer-finished").Find(&contents)

	for _, content := range contents {
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
	i.LightNode.DB.Model(&core.Content{}).Where("status <> ? and created_at < ?", "transfer-finished", time.Now().AddDate(0, 0, -7)).Updates(core.Content{
		Status: "transfer-failed",
	})
	return nil
}
