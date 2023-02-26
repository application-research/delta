package api

import (
	"delta/core"
	"github.com/labstack/echo/v4"
)

// repair deals (re-create or re-try)
func ConfigureRepairRouter(e *echo.Group, node *core.DeltaNode) {

	repair := e.Group("/repair")

	repair.GET("/deal/content", func(c echo.Context) error {

		// retry the same content id
		return nil
	})

	repair.GET("/piece-commitment", func(c echo.Context) error {

		// retry the same piece-commitment
		return nil
	})

}
