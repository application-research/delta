package core

import (
	"fmt"
	model "github.com/application-research/delta-db/db_models"
	"gorm.io/gorm"
	"time"
)

type ReplicationService struct {
	// `NewReplicationService` is a function that returns a `ReplicationService` struct
	LightNode *DeltaNode
}

// NewReplicationService Creating a new `ReplicationService` struct.
func NewReplicationService(ln *DeltaNode) *ReplicationService {
	return &ReplicationService{
		LightNode: ln,
	}
}

type DealReplication struct {
	Content                      model.Content                       `json:"content"`
	ContentDealProposalParameter model.ContentDealProposalParameters `json:"deal_proposal_parameter"`
}

func (r ReplicationService) ReplicateContent(contentSource DealReplication, numberOfReplication int, txn *gorm.DB) []model.Content {
	var replicatedContents []model.Content
	for i := 0; i < numberOfReplication; i++ {
		var newContent model.Content
		var newContentDealProposalParameter model.ContentDealProposalParameters
		newContent = contentSource.Content
		newContentDealProposalParameter = contentSource.ContentDealProposalParameter
		newContent.ID = 0

		err := txn.Create(&newContent).Error
		if err != nil {
			//tx.Rollback()
			fmt.Println(err)
			return nil
		}

		newContentDealProposalParameter.ID = 0
		newContentDealProposalParameter.Content = newContent.ID
		err = txn.Create(&newContentDealProposalParameter).Error
		if err != nil {
			//tx.Rollback()
			fmt.Println(err)
			return nil
		}
		//	assign a miner
		minerAssignService := NewMinerAssignmentService()
		provider, errOnPv := minerAssignService.GetSPWithGivenBytes(newContent.Size)
		if errOnPv != nil {
			fmt.Println(errOnPv)
			return nil
		}

		contentMinerAssignment := model.ContentMiner{
			Miner:     provider.Address,
			Content:   newContent.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = txn.Create(&contentMinerAssignment).Error
		if err != nil {
			//tx.Rollback()
			fmt.Println(err)
			return nil
		}

		replicatedContents = append(replicatedContents, newContent)

	}
	return replicatedContents
}
