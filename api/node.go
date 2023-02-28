package api

import (
	"delta/core"

	"github.com/labstack/echo/v4"
)

// ConfigureNodeInfoRouter It configures the router to handle requests for node information
func ConfigureNodeInfoRouter(e *echo.Group, node *core.DeltaNode) {
	nodeGroup := e.Group("/node")
	nodeGroup.GET("/info", handleNodeInfo(node))
	nodeGroup.GET("/addr", handleNodeAddr(node))
	nodeGroup.GET("/peers", handleNodePeers(node))
	nodeGroup.GET("/host", handleNodeHost(node))
	nodeGroup.GET("/api-key", handleNodeHostApiKey(node))
}

// If the node is in standalone mode, return the API key
func handleNodeHostApiKey(node *core.DeltaNode) func(c echo.Context) error {

	if node.Config.Common.Mode != "standalone" {
		return func(c echo.Context) error {
			return c.JSON(200, "This is not a standalone node")
		}
	}

	// return the api key if standalone mode.
	return func(c echo.Context) error {
		return c.JSON(200, map[string]string{"standalone_api_key": node.Config.Standalone.APIKey})
	}
}

// It returns a function that takes a `DeltaNode` and returns a function that takes an `echo.Context` and returns an
// `error`
func handleNodeHost(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(200, node.Node.Host.ID())
	}
}

// It returns a function that takes a `DeltaNode` and returns a function that takes a `Context` and returns an `error`
func handleNodePeers(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(200, node.Node.Host.Network().Peers())
	}
}

// It returns a function that takes a `DeltaNode` and returns a function that takes a `Context` and returns an `error`
func handleNodeAddr(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(200, node.Node.Host.Addrs())
	}
}

// It returns a function that returns a JSON response with the node's name, description, and type
func handleNodeInfo(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {

		ws := core.NewWebsocketService(node)
		ws.SendMessage("Hello, Client!")

		nodeName := node.Config.Node.Name
		nodeDescription := node.Config.Node.Description
		nodeType := node.Config.Node.Type

		return c.JSON(200, map[string]string{
			"name":        nodeName,
			"description": nodeDescription,
			"type":        nodeType,
		})
	}
}
