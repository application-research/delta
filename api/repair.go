package api

import (
	"delta/core"
	"delta/jobs"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
)

// ConfigureRepairRouter repair deals (re-create or re-try)
// It's a function that configures the repair router
func ConfigureRepairRouter(e *echo.Group, node *core.DeltaNode) {
	repair := e.Group("/repair")
	repair.GET("/deal/force-retry-all", handleForceRetryPendingContents(node))
	repair.GET("/deal/content/:contentId", handleRepairDealContent(node))
	repair.GET("/deal/piece-commitment/:pieceCommitmentId", handleRepairPieceCommitment(node))
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
