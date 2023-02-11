package jobs

import (
	"delta/core"
	"fmt"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"strconv"
	"time"
)

type DataTransferStatusListenerProcessor struct {
	LightNode *core.DeltaNode
}

func NewDataTransferStatusListenerProcessor(ln *core.DeltaNode) IProcessor {
	return &DataTransferStatusListenerProcessor{
		LightNode: ln,
	}
}

func (d DataTransferStatusListenerProcessor) Run() error {
	d.LightNode.FilClient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst filclient.ChannelState) {
		switch fst.Status {
		case datatransfer.Requested:
			d.LightNode.DB.Model(&core.ContentDeal{}).Where("id = ?", dbid).Updates(core.ContentDeal{
				TransferStarted: time.Now(),
			})
		case datatransfer.TransferFinished, datatransfer.Completed:
			transferId, err := strconv.Atoi(fst.TransferID)
			if err != nil {
				fmt.Println(err)
			}
			d.LightNode.DB.Model(&core.ContentDeal{}).Where("id = ?", dbid).Updates(core.ContentDeal{
				DealID:           int64(transferId),
				TransferFinished: time.Now(),
				SealedAt:         time.Now(),
				LastMessage:      "transfer-finished",
			})
			d.LightNode.DB.Model(&core.Content{}).Where("id = (select content from content_deals cd where cd.id = ?)", dbid).Updates(core.Content{
				Status: "transfer-finished",
			})
		case datatransfer.Failed:
			var contentDeal core.ContentDeal
			d.LightNode.DB.Model(&core.ContentDeal{}).Where("id = ?", dbid).Updates(core.ContentDeal{
				FailedAt: time.Now(),
			}).Find(&contentDeal)
			d.LightNode.DB.Model(&core.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(core.Content{
				Status: "transfer-failed",
			})

			d.LightNode.Dispatcher.AddJob(NewDataTransferRestartListenerProcessor(d.LightNode, contentDeal))
		default:

		}
	})
	return nil
}
