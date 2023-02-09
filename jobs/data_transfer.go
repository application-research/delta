package jobs

import (
	"delta/core"
	"fmt"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
)

type DataTransferListenerProcessor struct {
	LightNode   *core.LightNode
	ContentDeal core.ContentDeal
}

func NewDataTransferListenerProcessor(ln *core.LightNode, contentDeal core.ContentDeal) IProcessor {
	return &DataTransferListenerProcessor{
		LightNode:   ln,
		ContentDeal: contentDeal,
	}
}

func (d DataTransferListenerProcessor) Run() error {
	d.LightNode.Filclient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst filclient.ChannelState) {
		fmt.Println("dbid", dbid)
		switch fst.Status {
		case datatransfer.Requested:
			fmt.Println("Requested")
		case datatransfer.TransferFinished, datatransfer.Completed:
			fmt.Println("TransferFinished, Completed")
		default:
			fmt.Println("default")
		}
	})
	return nil
}
