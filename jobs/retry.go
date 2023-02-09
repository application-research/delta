package jobs

import (
	"delta/core"
)

type RetryProcessor struct {
	LightNode *core.LightNode
}

func NewRetryProcessor(ln *core.LightNode) IProcessor {
	return &RetryProcessor{
		LightNode: ln,
	}
}

func (i RetryProcessor) Run() error {

	// create the new logic again.
	// get all content with transfer status = "transfer-started" and "transfer-failed"
	var contents []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ?", "transfer-started").Find(&contents)

	for _, content := range contents {

		// get the piece
		var pieceCommp core.PieceCommitment
		i.LightNode.DB.Model(&core.PieceCommitment{}).Where("id = (select piece_commitment_id from contents c where c.id = ?)", content.ID).Find(&pieceCommp)

		// repair now
		i.LightNode.Dispatcher.AddJob(NewStorageDealMakerProcessor(i.LightNode, content, pieceCommp))
	}

	// get all the pending content jobs. we need to requeue them.

	var pinnedContents []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ?", "pinned").Find(&pinnedContents)

	for _, content := range contents {
		i.LightNode.Dispatcher.AddJob(NewPieceCommpProcessor(i.LightNode, content))
	}

	var contentsForCommp []core.Content
	i.LightNode.DB.Model(&core.Content{}).Where("status = ?", "piece-computing").Find(&contentsForCommp)

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

	return nil
}
