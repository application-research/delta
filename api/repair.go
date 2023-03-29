package api

import (
	"delta/core"
	"delta/jobs"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
)

type RetryDealResponse struct {
	Status       string      `json:"status"`
	Message      string      `json:"message"`
	NewContentId int64       `json:"new_content_id,omitempty"`
	OldContentId interface{} `json:"old_content_id,omitempty"`
}

// ConfigureRepairRouter repair deals (re-create or re-try)
// It's a function that configures the repair router
func ConfigureRepairRouter(e *echo.Group, node *core.DeltaNode) {
	repair := e.Group("/repair")
	repair.GET("/deal/force-retry-all", handleForceRetryPendingContents(node))
	repair.GET("/deal/content/:contentId", handleRepairDealContent(node))
	repair.GET("/deal/piece-commitment/:pieceCommitmentId", handleRepairPieceCommitment(node))

	retry := e.Group("/retry")
	retry.GET("/deal/:contentId", handleRetryContent(node))
}

func handleRetryContent(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		var contentId = c.Param("contentId")

		// look up the content
		var content model.Content
		node.DB.Model(&model.Content{}).Where("id = ?", contentId).First(&content)

		// re-create the same content
		var newContent = new(model.Content)
		newContent = &content
		newContent.ID = 0
		node.DB.Model(&model.Content{}).Create(&newContent) // create new

		// re-queue the content
		processor := jobs.NewPieceCommpProcessor(node, *newContent)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		err := c.JSON(200, RetryDealResponse{
			Status:       "success",
			Message:      "Deal request received. Please take note of the new_content_id. You can use the content_id to check the status of the deal.",
			NewContentId: newContent.ID,
		})
		if err != nil {
			return err
		}

		return nil
	}
}

// It takes a piece commitment id, finds the piece commitment, and re-queues the job
func handleRepairPieceCommitment(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		// get the piece commitment'
		var pieceCommitment model.PieceCommitment
		node.DB.Model(&model.PieceCommitment{}).Where("id = ?", c.Param("pieceCommitmentId")).First(&pieceCommitment)

		// if the piece commitment is not found, throw an error.
		if pieceCommitment.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "piece commitment not found",
			})
		}

		// if the piece commitment is found, re-queue the job.
		var content model.Content
		node.DB.Model(&model.Content{}).Where("piece_commitment_id = ?", pieceCommitment.ID).First(&content)

		processor := jobs.NewPieceCommpProcessor(node, content)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "re-queued piece commitment",
		})
	}
}

// It creates a new job processor, adds it to the dispatcher, and returns a JSON response
func handleForceRetryPendingContents(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		processor := jobs.NewRetryProcessor(node)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying all deals",
		})

	}
}

// It takes a content ID, finds the content deal entry for that content, and then retries the deal
func handleRepairDealContent(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		paramContentId := c.Param("contentId")

		// get the content deal entry
		var contentDeal model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("content = ?", paramContentId).First(&contentDeal)

		// if not content deal entry, throw an error.
		if contentDeal.ID == 0 {
			return c.JSON(200, map[string]interface{}{
				"message": "content deal not found",
			})
		}

		// if the deal is not in the right state, throw an error.
		var content model.Content
		node.DB.Model(&model.Content{}).Where("id = ?", paramContentId).First(&content)

		// retry it.
		processor := jobs.NewPieceCommpProcessor(node, content)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying deal",
			"content": content,
		})
	}
}
