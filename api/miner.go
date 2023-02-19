package api

import (
	"delta/core"
	"delta/core/model"
	"github.com/labstack/echo/v4"
)

// https://api.estuary.tech/public/miners/storage/query/f01963614
// https://api.estuary.tech/public/miners/

// ConfigureAdminRouter configures the admin router
// This is the router that is used to administer the node
func ConfigureMinerRouter(e *echo.Group, node *core.DeltaNode) {

	configureMiner := e.Group("/miner")

	//	get stats of miner
	configureMiner.GET("/stats/:minerId", func(c echo.Context) error {
		minerId := c.Param("minerId")
		var dealsListMiner []model.ContentDeal
		node.DB.Model(&model.ContentDeal{}).Where("miner = ?", minerId).Order("created_at desc").Find(&dealsListMiner)
		c.JSON(200, dealsListMiner)
		return nil
	})

	// 	get the stats of miner and content
	//	get the stats of miner and content id
}
