package core

import (
	"context"
	"delta/utils"
	"fmt"
	model "delta/models"
	fc "github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/ipfs/go-cid"
	"time"
)

// SetDataTransferEventsSubscribe `filc.SubscribeToDataTransferEvents(func(event datatransfer.Event, channelState datatransfer.ChannelState) {`
//
// The above function is a callback function that is called whenever a data transfer event occurs. The callback function
// takes two arguments: `event` and `channelState`. The `event` argument is of type `datatransfer.Event` and the
// `channelState` argument is of type `datatransfer.ChannelState`
func SetDataTransferEventsSubscribe(i *DeltaNode) {
	fmt.Println(utils.Purple + "Subscribing to transfer channel events..." + utils.Reset)
	i.FilClient.SubscribeToDataTransferEvents(func(event datatransfer.Event, channelState datatransfer.ChannelState) {
		switch event.Code {
		case datatransfer.DataQueued:
			fmt.Println("Data Transfer Queued event: ", event, " for transfer id: ", channelState.TransferID(), " for db id: ", channelState.BaseCID())
		case datatransfer.Complete:
			fmt.Println("Data Transfer Complete event: ", event, " for transfer id: ", channelState.TransferID(), " for db id: ", channelState.BaseCID())
		case datatransfer.Error, datatransfer.Disconnected, datatransfer.ReceiveDataError, datatransfer.Cancel, datatransfer.RequestTimedOut, datatransfer.SendDataError:
			fmt.Println("Data Transfer Error event: ", event, " for transfer id: ", channelState.TransferID(), " for db id: ", channelState.BaseCID())
		}
	})
}

// SetLibp2pManagerSubscribe It subscribes to the libp2p transfer manager and updates the database with the status of the transfer
func SetLibp2pManagerSubscribe(i *DeltaNode) {

	fmt.Println(utils.Purple + "Subscribing to transfer channel states..." + utils.Reset)
	i.FilClient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst fc.ChannelState) {
		//fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
		switch fst.Status {
		case datatransfer.Requested:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			var contentDeal model.ContentDeal
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Find(&contentDeal)
			// save the content deal
			contentDeal.TransferStarted = time.Now()
			contentDeal.UpdatedAt = time.Now()
			contentDeal.LastMessage = utils.DEAL_STATUS_TRANSFER_STARTED
			i.DB.Save(&contentDeal)

		case datatransfer.TransferFinished, datatransfer.Completed:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			//transferId, err := strconv.Atoi(fst.TransferID)
			//if err != nil {
			//	fmt.Println(err)
			//}

			// save the content deal
			var contentDeal model.ContentDeal
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Find(&contentDeal)
			contentDeal.TransferFinished = time.Now()
			contentDeal.SealedAt = time.Now()
			contentDeal.UpdatedAt = time.Now()
			contentDeal.OnChainAt = time.Now()
			contentDeal.LastMessage = utils.DEAL_STATUS_TRANSFER_FINISHED
			i.DB.Save(&contentDeal)

			// save the content status
			var content model.Content
			i.DB.Model(&model.Content{}).Where("id in (select cd.content from content_deals cd where cd.id = ?)", dbid).Find(&content)
			content.Status = utils.DEAL_STATUS_TRANSFER_FINISHED
			content.LastMessage = utils.DEAL_STATUS_TRANSFER_FINISHED
			content.UpdatedAt = time.Now()
			i.DB.Save(&content)

			// remove from the blockstore
			cidToDelete, err := cid.Decode(content.Cid)
			if err != nil {
				fmt.Println(err)
			}
			if i.Config.Node.KeepCopies {
				fmt.Println("Keeping a copy of the content - not removing from the blockstore - CID: ", cidToDelete, "")
			} else {
				go i.Node.Blockservice.DeleteBlock(context.Background(), cidToDelete)
			}
		case datatransfer.Failed, datatransfer.Failing, datatransfer.Cancelled, datatransfer.InitiatorPaused, datatransfer.ResponderPaused, datatransfer.ChannelNotFoundError:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			var contentDeal model.ContentDeal
			i.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Find(&contentDeal)
			contentDeal.LastMessage = fst.Message
			contentDeal.UpdatedAt = time.Now()
			contentDeal.FailedAt = time.Now()
			i.DB.Save(&contentDeal)

			var content model.Content
			i.DB.Model(&model.Content{}).Where("id in (select cd.content from content_deals cd where cd.id = ?)", dbid).Find(&content)
			content.Status = utils.DEAL_STATUS_TRANSFER_FAILED
			content.LastMessage = fst.Message
			content.UpdatedAt = time.Now()
			i.DB.Save(&content)

			// remove from the blockstore
			cidToDelete, err := cid.Decode(content.Cid)
			if err != nil {
				fmt.Println(err)
			}
			if i.Config.Node.KeepCopies {
				fmt.Println("Keeping a copy of the content - not removing from the blockstore - CID: ", cidToDelete, "")
			} else {
				go i.Node.Blockservice.DeleteBlock(context.Background(), cidToDelete)
			}
		default:
		}
	})
}
