package utils

import (
	"delta/core"
	"delta/core/model"
	"delta/jobs"
	"fmt"
	fc "github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"strconv"
	"time"
)

func SetFilclientLibp2pSubscribe(filc *fc.FilClient, i *core.DeltaNode) {
	filc.Libp2pTransferMgr.Subscribe(func(dbid uint, fst fc.ChannelState) {
		switch fst.Status {
		case datatransfer.Requested:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				TransferStarted: time.Now(),
				UpdatedAt:       time.Now(),
			})
		case datatransfer.TransferFinished, datatransfer.Completed:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			transferId, err := strconv.Atoi(fst.TransferID)
			if err != nil {
				fmt.Println(err)
			}
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				DealID:           int64(transferId),
				TransferFinished: time.Now(),
				SealedAt:         time.Now(),
				UpdatedAt:        time.Now(),
				LastMessage:      DEAL_STATUS_TRANSFER_FINISHED,
			})
			i.DB.Model(&model.Content{}).Where("id in (select cd.content from content_deals cd where cd.id = ?)", dbid).Updates(model.Content{
				Status:    DEAL_STATUS_TRANSFER_FINISHED,
				UpdatedAt: time.Now(),
			})
		case datatransfer.Failed:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			var contentDeal model.ContentDeal
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				FailedAt:  time.Now(),
				UpdatedAt: time.Now(),
			}).Find(&contentDeal)
			i.DB.Model(&model.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(model.Content{
				Status:    DEAL_STATUS_TRANSFER_FAILED,
				UpdatedAt: time.Now(),
			})

			i.Dispatcher.AddJobAndDispatch(jobs.NewDataTransferRestartProcessor(i, contentDeal), 1)
		default:
		}
	})
}
