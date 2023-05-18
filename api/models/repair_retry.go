package models

import model "github.com/application-research/delta-db/db_models"

type RetryDealResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	NewContentId int64       `json:"new_content_id,omitempty"`
	OldContentId interface{} `json:"old_content_id,omitempty"`
}

type MultipleImportRequest struct {
	ContentID   string      `json:"content_id"`
	DealRequest DealRequest `json:"metadata"`
}

type ImportRetryRequest struct {
	ContentIds []string `json:"content_ids"`
}
type ImportRetryResponse struct {
	Message     string        `json:"message"`
	Content     model.Content `json:"content"`
	DealRequest DealRequest   `json:"metadata"`
}
