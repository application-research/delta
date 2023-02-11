package jobs

import (
	"delta/core"
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

func (i RetryProcessor) Run() error {

	// create the new logic again.
	// if transfer-started but older than 3 days, then requeue the job.

	var contents []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ? and created_at < ?", "transfer-started", time.Now().AddDate(0, 0, -3)).Find(&contents)

	for _, content := range contents {

		// get the piece
		var pieceCommp core.PieceCommitment
		i.LightNode.DB.Model(&core.PieceCommitment{}).Where("id = (select piece_commitment_id from contents c where c.id = ?)", content.ID).Find(&pieceCommp)

		// repair now
		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, content, pieceCommp))
	}

	// get all the pinned content jobs that are stuck for more than 1 day
	var pinnedContents []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ? and created_at < ?", "pinned", time.Now().AddDate(0, 0, -1)).Find(&pinnedContents)

	for _, content := range contents {
		i.LightNode.Dispatcher.AddJob(NewPieceCommpProcessor(i.LightNode, content))
	}

	// get all piece computing jobs that are stuck for more than 1 day
	var contentsForCommp []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ? and created_at < ?", "piece-computing", time.Now().AddDate(0, 0, -1)).Find(&contentsForCommp)

	for _, content := range contentsForCommp {
		i.LightNode.Dispatcher.AddJob(NewPieceCommpProcessor(i.LightNode, content))
	}

	var pieceCommps []core.PieceCommitment
	i.LightNode.DB.Model(&core.PieceCommitment{}).Where("status = ?", "open").Find(&pieceCommps)

	for _, pieceCommp := range pieceCommps {
		var content core.Content
		i.LightNode.DB.Model(&core.Content{}).Where("piece_commitment_id = ?", pieceCommp.ID).Find(&content)
		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, content, pieceCommp))
	}

	// we also want to retry those request that are not failed but not "transfer-finished" that are within the 48 hour period.

	return nil
}
