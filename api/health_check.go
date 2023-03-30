package api

import (
	"delta/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

// > ConfigureHealthCheckRouter is a function that takes a pointer to an echo.Group and a pointer to a DeltaNode and
// returns nothing
func ConfigureHealthCheckRouter(healthCheckApiGroup *echo.Group, node *core.DeltaNode) {

	healthCheckAuthApiGroup := healthCheckApiGroup.Group("/check/auth")

	//	health check api withouth auth
	healthCheckApiGroup.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	//	health check api
	healthCheckAuthApiGroup.Use(Authenticate(*DeltaNodeConfig))
	healthCheckAuthApiGroup.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
}
