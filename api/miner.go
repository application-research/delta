package api

import (
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
)

// https://api.estuary.tech/public/miners/storage/query/f01963614
// https://api.estuary.tech/public/miners/

// ConfigureMinerRouter ConfigureAdminRouter configures the admin router
// This is the router that is used to administer the node
func ConfigureMinerRouter(e *echo.Group, node *core.DeltaNode) {

	configureMiner := e.Group("/miner")

	//	get stats of miner
	configureMiner.GET("/stats/:minerId", handleMinerStats(node))

	// 	get the stats of miner and content
	//	get the stats of miner and content id
}

// A function that returns a function that takes a context and returns an error.
func handleMinerStats(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		minerId := c.Param("minerId")
		var dealsListMiner []model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("miner = ?", minerId).Order("created_at desc").Find(&dealsListMiner)
		c.JSON(200, dealsListMiner)
		return nil
	}
}
