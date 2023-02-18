package jobs

import (
	"context"
	"delta/core"
	"delta/core/model"
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
	var contents []model.Content
	i.LightNode.DB.Model(&model.Content{}).Where("status <> ? and created_at < ?", "transfer-finished", time.Now().AddDate(0, 0, -1)).Find(&contents)

	for _, content := range contents {

		// get the piece
		if content.Status == utils.CONTENT_PINNED || content.Status == utils.CONTENT_PIECE_COMPUTING || content.Status == utils.CONTENT_PIECE_COMPUTED || content.Status == utils.CONTENT_DEAL_PROPOSAL_SENT || content.Status == utils.CONTENT_DEAL_SENDING_PROPOSAL || content.Status == utils.CONTENT_DEAL_MAKING_PROPOSAL {
			var pieceCommp model.PieceCommitment
			i.LightNode.DB.Model(&model.PieceCommitment{}).Where("id = (select piece_commitment_id from contents c where c.id = ?)", content.ID).Find(&pieceCommp)
			i.LightNode.Dispatcher.AddJobAndDispatch(NewStorageDealMakerProcessor(i.LightNode, content, pieceCommp), 1)
		} else if content.Status == utils.CONTENT_FAILED_TO_PIN || content.Status == utils.DEAL_STATUS_TRANSFER_FAILED || content.Status == utils.CONTENT_DEAL_PROPOSAL_FAILED || content.Status == utils.CONTENT_PIECE_COMPUTING_FAILED {
			// delete/ignore
			cidToDelete, err := cid.Decode(content.Cid)
			if err != nil {
				fmt.Println("error in decoding cid", err)
				continue
			}
			i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidToDelete)
		} else {
			// fail it entirely
			content.Status = utils.CONTENT_FAILED_TO_PROCESS
			content.LastMessage = "failed to process even after retrying."
			i.LightNode.DB.Model(&model.Content{}).Where("id = ?", content.ID).Updates(content)
			cidToDelete, err := cid.Decode(content.Cid)
			if err != nil {
				fmt.Println("error in decoding cid", err)
				continue
			}
			i.LightNode.Node.Blockservice.DeleteBlock(context.Background(), cidToDelete)
		}

	}

	return nil
}
