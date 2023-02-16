package api

import (
	"delta/core"

	"github.com/labstack/echo/v4"
)

func ConfigureNodeInfoRouter(e *echo.Group, node *core.DeltaNode) {
	nodeGroup := e.Group("/node")
	nodeGroup.GET("/info", func(c echo.Context) error {
		nodeName := node.Config.Node.Name
		nodeDescription := node.Config.Node.Description
		nodeType := node.Config.Node.Type

		return c.JSON(200, map[string]string{
			"name":        nodeName,
			"description": nodeDescription,
			"type":        nodeType,
		})
	})

	nodeGroup.GET("/addr", func(c echo.Context) error {
		return c.JSON(200, node.Node.Host.Addrs())
	})

	nodeGroup.GET("/peers", func(c echo.Context) error {
		return c.JSON(200, node.Node.Host.Network().Peers())
	})

	nodeGroup.GET("/host", func(c echo.Context) error {
		return c.JSON(200, node.Node.Host.ID())
	})
}
