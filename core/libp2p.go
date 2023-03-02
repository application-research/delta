package core

import (
	"context"
	"delta/utils"
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	fc "github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/ipfs/go-cid"
	"strconv"
	"time"
)

// `filc.SubscribeToDataTransferEvents(func(event datatransfer.Event, channelState datatransfer.ChannelState) {`
//
// The above function is a callback function that is called whenever a data transfer event occurs. The callback function
// takes two arguments: `event` and `channelState`. The `event` argument is of type `datatransfer.Event` and the
// `channelState` argument is of type `datatransfer.ChannelState`
func SetDataTransferEventsSubscribe(i *DeltaNode) {
	i.FilClient.SubscribeToDataTransferEvents(func(event datatransfer.Event, channelState datatransfer.ChannelState) {
		fmt.Println("Event: ", event, " for transfer id: ", channelState.TransferID(), " for db id: ", channelState.BaseCID())
		switch event.Code {
		case datatransfer.Error, datatransfer.Disconnected, datatransfer.ReceiveDataError, datatransfer.Cancel, datatransfer.RequestTimedOut, datatransfer.SendDataError:
			fmt.Println("Error event: ", event, " for transfer id: ", channelState.TransferID(), " for db id: ", channelState.BaseCID())
		}
	})
}

// It subscribes to the libp2p transfer manager and updates the database with the status of the transfer
func SetLibp2pManagerSubscribe(i *DeltaNode) {
	i.FilClient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst fc.ChannelState) {
		fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
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
				OnChainAt:        time.Now(),
				LastMessage:      utils.DEAL_STATUS_TRANSFER_FINISHED,
			})
			var content model.Content
			i.DB.Model(&model.Content{}).Where("id in (select cd.content from content_deals cd where cd.id = ?)", dbid).Find(&content)
			content.Status = utils.DEAL_STATUS_TRANSFER_FINISHED
			content.UpdatedAt = time.Now()
			i.DB.Save(&content)

			// remove from the blockstore
			cidToDelete, err := cid.Decode(content.Cid)
			go i.Node.DAGService.Remove(context.Background(), cidToDelete)
		case datatransfer.Failed, datatransfer.Failing, datatransfer.Cancelled, datatransfer.InitiatorPaused, datatransfer.ResponderPaused, datatransfer.ChannelNotFoundError:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			var contentDeal model.ContentDeal
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				FailedAt:  time.Now(),
				UpdatedAt: time.Now(),
			}).Find(&contentDeal)
			i.DB.Model(&model.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(model.Content{
				Status:    utils.DEAL_STATUS_TRANSFER_FAILED,
				UpdatedAt: time.Now(),
			})
		default:
		}
	})
}
