package jobs

import (
	"context"
	"delta/core"
	"fmt"
	"github.com/ipfs/go-cid"
)

type ItemContentCleanUpProcessor struct {
	Context   context.Context
	LightNode *core.LightNode
}

func NewItemContentCleanUpProcessor(ln *core.LightNode) IProcessor {
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

	return nil
}
