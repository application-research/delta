package api

import (
	"context"
	"delta/core"
	"github.com/application-research/delta-db/db_models"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// ConfigureNodeInfoRouter It configures the router to handle requests for node information
func ConfigureNodeInfoRouter(e *echo.Group, node *core.DeltaNode) {

	nodeGroup := e.Group("/node")
	nodeGroup.GET("/info", handleNodeInfo(node))
	nodeGroup.GET("/uuids", handleNodeUuidInfo(node))
	nodeGroup.GET("/addr", handleNodeAddr(node))
	nodeGroup.GET("/peers", handleNodePeers(node))
	nodeGroup.GET("/host", handleNodeHost(node))
	nodeGroup.GET("/api-key", handleNodeHostApiKey(node))
}

// If the node is in standalone mode, return the API key
func handleNodeHostApiKey(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, span := otel.Tracer("handleNodeHostApiKey").Start(context.Background(), "handleNodeHostApiKey")
		defer span.End()

		span.SetName("ConfigureNodeInfoRouter")
		span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
		span.SetAttributes(attribute.String("path", c.Path()))
		span.SetAttributes(attribute.String("method", c.Request().Method))

		if node.Config.Common.Mode != "standalone" {
			return c.JSON(200, "This is not a standalone node")
		}
		return c.JSON(200, map[string]string{"standalone_api_key": node.Config.Standalone.APIKey})
	}

}

// It returns a function that takes a `DeltaNode` and returns a function that takes an `echo.Context` and returns an
// `error`
func handleNodeHost(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, span := otel.Tracer("handleNodeHost").Start(context.Background(), "handleNodeHostApiKey")
		defer span.End()

		span.SetName("ConfigureNodeInfoRouter")
		span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
		span.SetAttributes(attribute.String("path", c.Path()))
		span.SetAttributes(attribute.String("method", c.Request().Method))

		return c.JSON(200, node.Node.Host.ID())
	}
}

// It returns a function that takes a `DeltaNode` and returns a function that takes a `Context` and returns an `error`
func handleNodePeers(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, span := otel.Tracer("handleNodePeers").Start(context.Background(), "handleNodeHostApiKey")
		defer span.End()

		span.SetName("ConfigureNodeInfoRouter")
		span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
		span.SetAttributes(attribute.String("path", c.Path()))
		span.SetAttributes(attribute.String("method", c.Request().Method))
		return c.JSON(200, node.Node.Host.Network().Peers())
	}
}

// It returns a function that takes a `DeltaNode` and returns a function that takes a `Context` and returns an `error`
func handleNodeAddr(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, span := otel.Tracer("handleNodePeers").Start(context.Background(), "handleNodeHostApiKey")
		defer span.End()

		span.SetName("handleNodeAddr")
		span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
		span.SetAttributes(attribute.String("path", c.Path()))
		span.SetAttributes(attribute.String("method", c.Request().Method))
		return c.JSON(200, node.Node.Host.Addrs())
	}
}

type Uuids struct {
	ID           uint   `json:"id"`
	InstanceUuid string `json:"instance_uuid"`
	CreatedAt    string `json:"created_at"`
}

// It returns a function that returns a JSON response with the node's name, description, and type
func handleNodeUuidInfo(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, span := otel.Tracer("handleNodeUuidInfo").Start(context.Background(), "handleNodeUuidInfo")
		defer span.End()

		span.SetName("handleNodeUuidInfo")
		span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
		span.SetAttributes(attribute.String("path", c.Path()))
		span.SetAttributes(attribute.String("method", c.Request().Method))

		// instance meta
		var uuids []Uuids
		node.DB.Model(&db_models.InstanceMeta{}).Scan(&uuids)

		return c.JSON(200, uuids)
	}
}

// It returns a function that returns a JSON response with the node's name, description, and type
func handleNodeInfo(node *core.DeltaNode) func(c echo.Context) error {
	return func(c echo.Context) error {
		_, span := otel.Tracer("handleNodePeers").Start(context.Background(), "handleNodeHostApiKey")
		defer span.End()

		span.SetName("handleNodeInfo")
		span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
		span.SetAttributes(attribute.String("path", c.Path()))
		span.SetAttributes(attribute.String("method", c.Request().Method))

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
