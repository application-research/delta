package jobs

import (
	"context"
	"fc-deal-making-service/core"
	"fmt"
	"github.com/ipfs/go-cid"
)

type ItemContentCleanUpProcessor struct {
	ContentProcessor
}

func NewItemContentCleanUpProcessor(ln *core.LightNode, content core.Content) IProcessor {
	return &ItemContentCleanUpProcessor{
		ContentProcessor{
			LightNode: ln,
			Content:   content,
			Context:   context.Background(),
		},
	}
}

func (i ItemContentCleanUpProcessor) Run() error {

	cidD, err := cid.Decode(i.Content.Cid)
	fmt.Println("cleaning up" + i.Content.Cid)
	if err != nil {
		fmt.Println("error on cid")
	}
	i.LightNode.Node.DAGService.Remove(i.Context, cidD)
	return nil
}
