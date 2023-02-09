package jobs

import (
	"delta/core"
	"fmt"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"time"
)

type DataTransferStatusListenerProcessor struct {
	LightNode *core.LightNode
}

func NewDataTransferStatusListenerProcessor(ln *core.LightNode) IProcessor {
	return &DataTransferStatusListenerProcessor{
		LightNode: ln,
	}
}

func (d DataTransferStatusListenerProcessor) Run() error {
	d.LightNode.FilClient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst filclient.ChannelState) {
		fmt.Println("data transfer status listener", dbid, fst.Status)
		switch fst.Status {
		case datatransfer.Requested:
			d.LightNode.DB.Model(&core.ContentDeal{}).Where("id = ?", dbid).Updates(core.ContentDeal{
				TransferStarted: time.Now(),
			})
			d.LightNode.DB.Model(&core.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(core.Content{
				Status: "transfer-requested",
			})
		case datatransfer.Ongoing:
			d.LightNode.DB.Model(&core.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(core.Content{
				Status: "transfer-ongoing",
			})
		case datatransfer.TransferFinished, datatransfer.Completed:
			d.LightNode.DB.Model(&core.ContentDeal{}).Where("id = ?", dbid).Updates(core.ContentDeal{
				TransferFinished: time.Now(),
				SealedAt:         time.Now(),
			})
			d.LightNode.DB.Model(&core.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(core.Content{
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

			// retry
			d.LightNode.Dispatcher.AddJob(NewDataTransferRestartListenerProcessor(d.LightNode, contentDeal))
		default:

		}
	})
	return nil
}
