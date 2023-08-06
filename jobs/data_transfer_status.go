package jobs

import (
	"delta/core"
	"delta/utils"
	"fmt"
	model "delta/models"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"strconv"
	"time"
)

// DataTransferStatusListenerProcessor It's a struct that contains a pointer to a DeltaNode.
// @property LightNode - This is the DeltaNode that will be used to process the data transfer.
type DataTransferStatusListenerProcessor struct {
	LightNode *core.DeltaNode
}

// NewDataTransferStatusListenerProcessor It creates a new instance of the `DataTransferStatusListenerProcessor` struct, and returns a pointer to it
func NewDataTransferStatusListenerProcessor(ln *core.DeltaNode) IProcessor {
	return &DataTransferStatusListenerProcessor{
		LightNode: ln,
	}
}

// Run It's a function that is called when the data transfer status changes.
func (d DataTransferStatusListenerProcessor) Run() error {
	d.LightNode.FilClient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst filclient.ChannelState) {
		fmt.Println("Data Transfer Status Listener: ", fst.Status)
		switch fst.Status {
		case datatransfer.Requested:
			d.LightNode.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				TransferStarted: time.Now(),
			})
		case datatransfer.TransferFinished, datatransfer.Completed:
			transferId, err := strconv.Atoi(fst.TransferID)
			if err != nil {
				fmt.Println(err)
			}
			d.LightNode.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				DealID:           int64(transferId),
				TransferFinished: time.Now(),
				SealedAt:         time.Now(),
				LastMessage:      utils.DEAL_STATUS_TRANSFER_FINISHED,
			})
			d.LightNode.DB.Model(&model.Content{}).Where("id = (select content from content_deals cd where cd.id = ?)", dbid).Updates(model.Content{
				Status: utils.DEAL_STATUS_TRANSFER_FINISHED,
			})
		case datatransfer.Failed:
			var contentDeal model.ContentDeal
			d.LightNode.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				FailedAt: time.Now(),
			}).Find(&contentDeal)
			d.LightNode.DB.Model(&model.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(model.Content{
				Status: utils.DEAL_STATUS_TRANSFER_FAILED,
			})

			d.LightNode.Dispatcher.AddJobAndDispatch(NewDataTransferRestartProcessor(d.LightNode, contentDeal), 1)
		default:
		}
	})
	fmt.Println("Data Transfer Status Listener Ended")

	return nil
}
