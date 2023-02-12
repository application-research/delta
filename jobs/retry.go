package jobs

import (
	"context"
	"delta/core"
	"delta/utils"
	"fmt"
	"github.com/ipfs/go-cid"
	"time"
)

type RetryProcessor struct {
	LightNode *core.DeltaNode
}

func NewRetryProcessor(ln *core.DeltaNode) IProcessor {
	return &RetryProcessor{
		LightNode: ln,
	}
}

// Run DB heavy process. We need to check the status of the content and requeue the job if needed.
func (i RetryProcessor) Run() error {

	// create the new logic again.
	// if transfer-started but older than 3 days, then requeue the job.
	var contents []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status <> ? and created_at < ?", "transfer-finished", time.Now().AddDate(0, 0, -1)).Find(&contents)

	for _, content := range contents {

		// get the piece
		if content.Status == utils.CONTENT_PINNED || content.Status == utils.CONTENT_PIECE_COMPUTING || content.Status == utils.CONTENT_PIECE_COMPUTED || content.Status == utils.CONTENT_DEAL_PROPOSAL_SENT || content.Status == utils.CONTENT_DEAL_SENDING_PROPOSAL || content.Status == utils.CONTENT_DEAL_MAKING_PROPOSAL {
			var pieceCommp core.PieceCommitment
			i.LightNode.DB.Model(&core.PieceCommitment{}).Where("id = (select piece_commitment_id from contents c where c.id = ?)", content.ID).Find(&pieceCommp)
			i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, content, pieceCommp))
		} else if content.Status == utils.DEAL_STATUS_TRANSFER_FAILED {
			// delete/ignore
			cidToDelete, err := cid.Decode(content.Cid)
			if err != nil {
				fmt.Println("error in decoding cid", err)
				continue
			}
			i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidToDelete)
		}

	}

	// blockstore check. if the cid is not on the blockstore, don't try, just delete/fail it.
	// check content status. if it's failed, don't even.
	return nil
}
