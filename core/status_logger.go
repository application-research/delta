package core

import model "delta/models"

// StatusLogger is used to change the status of each system objects in a async manner.
// Each state content, commp and transfers needs to be updated and logged in the database.

type StatusLogger struct {
	LightNode DeltaNode
}

func NewStatusLogger(node DeltaNode) *StatusLogger {
	return &StatusLogger{
		LightNode: node,
	}
}

// UpdateContentStatus updates the status of a content object.
func (s *StatusLogger) UpdateContentStatus(content model.Content, status string) error {
	tx := s.LightNode.DB.Model(&content).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// UpdatePieceCommStatus Updating the status of a piece commitment object.
func (s *StatusLogger) UpdatePieceCommStatus(pieceCommp model.PieceCommitment, status string) error {
	tx := s.LightNode.DB.Model(&pieceCommp).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// UpdateContentDealStatus Updating the status of a content deal object.
func (s *StatusLogger) UpdateContentDealStatus(pieceCommp model.ContentDeal, status string) error {
	tx := s.LightNode.DB.Model(&pieceCommp).Update("status", status)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
