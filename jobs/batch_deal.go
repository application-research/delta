package jobs

import (
	"context"
	"delta/api/models"
	"delta/core"
	"delta/utils"
	"fmt"
	"github.com/application-research/delta-db/db_models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type BatchDealRequestProcessor struct {
	Context     context.Context
	DealRequest []models.DealRequest
	LightNode   *core.DeltaNode
}

func (b BatchDealRequestProcessor) Run() error {
	//TODO implement me
	panic("implement me")
}

func NewBatchDealRequestProcessor(ln *core.DeltaNode, dealRequest []models.DealRequest) IProcessor {
	return &BatchDealRequestProcessor{
		Context:     context.Background(),
		LightNode:   ln,
		DealRequest: dealRequest,
	}
}

type BatchJob struct {
	DealRequests []models.DealRequest
	Node         *core.DeltaNode
}

type BatchJobResult struct {
	DealResponses []models.DealResponse
	BatchImportID int64
	Error         error
}

func processBatchJob(job BatchJob) BatchJobResult {
	node := job.Node
	dealRequests := job.DealRequests

	var dealResponses []models.DealResponse
	for _, dealRequest := range dealRequests {
		// Process each deal request
		// ...
		fmt.Println(dealRequest)
		// Append the deal response to the list
		dealResponses = append(dealResponses, models.DealResponse{
			// Populate the DealResponse fields accordingly
		})
	}

	// Create a batch import object
	batchImportUuid := uuid.New().String()
	batchImport := db_models.BatchImport{
		Uuid:      batchImportUuid,
		Status:    utils.BATCH_IMPORT_STATUS_STARTED,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := node.DB.Create(&batchImport).Error
	if err != nil {
		return BatchJobResult{
			Error: errors.New("Error creating a batch import object"),
		}
	}

	// Update the batch import status
	var batchImportToBeUpdate db_models.BatchImport
	node.DB.Raw("SELECT * FROM batch_imports WHERE id = ?", batchImport.ID).Scan(&batchImportToBeUpdate)
	batchImportToBeUpdate.Status = utils.BATCH_IMPORT_STATUS_COMPLETED
	batchImportToBeUpdate.UpdatedAt = time.Now()
	node.DB.Save(&batchImportToBeUpdate)

	return BatchJobResult{
		DealResponses: dealResponses,
		BatchImportID: batchImport.ID,
	}
}
