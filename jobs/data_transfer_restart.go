package jobs

import (
	"context"
	"delta/core"
	"fmt"
	"github.com/application-research/filclient"
)

type DataTransferRestartListenerProcessor struct {
	LightNode   *core.LightNode
	ContentDeal core.ContentDeal
}

func NewDataTransferRestartListenerProcessor(ln *core.LightNode, contentDeal core.ContentDeal) IProcessor {
	return &DataTransferRestartListenerProcessor{
		LightNode:   ln,
		ContentDeal: contentDeal,
	}
}

func (d DataTransferRestartListenerProcessor) Run() error {
	// get the deal data transfer state pull deals
	channelId, err := d.ContentDeal.ChannelID()
	st, err := d.LightNode.FilClient.TransferStatus(context.Background(), &channelId)
	if err != nil && err != filclient.ErrNoTransferFound {
		return err
	}

	if st == nil {
		return fmt.Errorf("no data transfer state was found")
	}

	err = d.LightNode.FilClient.RestartTransfer(context.Background(), &channelId)
	if err != nil {
		return err
	}
	// subscribe to data transfer events
	d.LightNode.Dispatcher.AddJob(NewDataTransferStatusListenerProcessor(d.LightNode))
	return nil
}
