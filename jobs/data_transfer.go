package jobs

import (
	"delta/core"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"time"
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
		switch fst.Status {
		case datatransfer.Requested:
			d.LightNode.DB.Model(&d.ContentDeal).Where("id = ?", dbid).Updates(core.ContentDeal{
				TransferStarted: time.Now(),
			})
		case datatransfer.Ongoing:
		case datatransfer.TransferFinished, datatransfer.Completed:
			d.LightNode.DB.Model(&d.ContentDeal).Where("id = ?", dbid).Updates(core.ContentDeal{
				TransferFinished: time.Now(),
				SealedAt:         time.Now(),
			})
		case datatransfer.Failed:
			d.LightNode.DB.Model(&d.ContentDeal).Where("id = ?", dbid).Updates(core.ContentDeal{
				FailedAt: time.Now(),
			})
		default:

		}
	})
	return nil
}
