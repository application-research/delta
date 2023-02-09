package api

import (
	"delta/core"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// node info
// multi addr
// peer id
// host

func ConfigureNodeInfoRouter(e *echo.Group, node *core.LightNode) {
	nodeGroup := e.Group("/node")
	nodeGroup.GET("/info", func(c echo.Context) error {
		nodeName := viper.Get("NODE_NAME").(string)
		nodeDescription := viper.Get("NODE_DESCRIPTION").(string)
		nodeType := viper.Get("NODE_TYPE").(string)

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
