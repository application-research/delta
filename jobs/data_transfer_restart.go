package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"github.com/application-research/filclient"
)

type DataTransferRestartListenerProcessor struct {
	LightNode   *core.DeltaNode
	ContentDeal model.ContentDeal
}

func NewDataTransferRestartProcessor(ln *core.DeltaNode, contentDeal model.ContentDeal) IProcessor {
	return &DataTransferRestartListenerProcessor{
		LightNode:   ln,
		ContentDeal: contentDeal,
	}
}

func (d DataTransferRestartListenerProcessor) Run() error {
	// get the deal data transfer state pull deals
	dtChan, err := utils.GetChannelID(d.ContentDeal.DTChan)
	if err != nil {
		fmt.Println(err)
		return err
	}
	channelId := dtChan
	st, err := d.LightNode.FilClient.TransferStatus(context.Background(), &channelId)
	if err != nil && err != filclient.ErrNoTransferFound {
		fmt.Println(err)
		return err
	}

	if st == nil {
		return fmt.Errorf("no data transfer state was found")
	}

	err = d.LightNode.FilClient.RestartTransfer(context.Background(), &channelId)
	if err != nil {
		return err
	}
	return nil
}
