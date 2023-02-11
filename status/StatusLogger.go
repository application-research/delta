package status

import "delta/core"

// StatusLogger is used to change the status of each system objects in a async manner.
// Each state content, commp and transfers needs to be updated and logged in the database.

type StatusLogger struct {
	LightNode core.DeltaNode
}

func NewStatusLogger(node core.DeltaNode) *StatusLogger {
	return &StatusLogger{
		LightNode: node,
	}
}

// UpdateContentStatus updates the status of a content object.
func (s *StatusLogger) UpdateContentStatus(content core.Content, status string) error {
	tx := s.LightNode.DB.Model(&content).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *StatusLogger) UpdatePieceCommStatus(pieceCommp core.PieceCommitment, status string) error {
	tx := s.LightNode.DB.Model(&pieceCommp).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s *StatusLogger) UpdateContentDealStatus(pieceCommp core.ContentDeal, status string) error {
	tx := s.LightNode.DB.Model(&pieceCommp).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
