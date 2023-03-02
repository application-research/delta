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

	repair.GET("/force-retry-all", func(c echo.Context) error {
		processor := jobs.NewRetryProcessor(node)
		node.Dispatcher.AddJobAndDispatch(processor, 1)

		return c.JSON(200, map[string]interface{}{
			"message": "retrying all deals",
		})

	})

	repair.GET("/deal/content/:contentId", func(c echo.Context) error {

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
	})

	repair.GET("/piece-commitment", func(c echo.Context) error {

		// retry the same piece-commitment
		return nil
	})

}
