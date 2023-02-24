package api

import (
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
)

// ConfigureAdminRouter configures the admin router
// This is the router that is used to administer the node
func ConfigureAdminRouter(e *echo.Group, node *core.DeltaNode) {

	adminRepair := e.Group("/repair")
	adminWallet := e.Group("/wallet")
	adminDashboard := e.Group("/dashboard")
	adminStats := e.Group("/stats")

	adminStats.GET("/miner/:minerId", func(c echo.Context) error {

		var contents []model.Content
		node.DB.Raw("select c.* from content_deals cd, contents c where cd.content = c.id and cd.miner = ?", c.Param("minerId")).Scan(&contents)

		var contentMinerAssignment []model.ContentMiner
		node.DB.Raw("select cma.* from content_miners cma, contents c where cma.content = c.id and cma.miner = ?", c.Param("minerId")).Scan(&contentMinerAssignment)

		return c.JSON(200, map[string]interface{}{
			"content": contents,
			"cmas":    contentMinerAssignment,
		})
	})

	// repair endpoints
	adminRepair.GET("/deal", func(c echo.Context) error {
		return nil
	})

	adminRepair.GET("/commp", func(c echo.Context) error {
		return nil
	})

	adminRepair.GET("/run-cleanup", func(c echo.Context) error {
		return nil
	})

	adminRepair.GET("/retry-deal-making-content", func(c echo.Context) error {
		return nil
	})

	// add wallet_estuary endpoint
	adminWallet.POST("/add", func(c echo.Context) error {
		return nil
	})

	adminWallet.POST("/import", func(c echo.Context) error {
		return nil
	})

	// list wallet_estuary endpoint
	adminWallet.GET("/list", func(c echo.Context) error {
		return nil
	})

	adminWallet.GET("/info", func(c echo.Context) error {
		return nil
	})

	adminDashboard.GET("/index", func(c echo.Context) error {
		return nil
	})
}
