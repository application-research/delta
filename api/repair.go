package api

import (
	"delta/core"
	"github.com/labstack/echo/v4"
)

// repair deals (re-create or re-try)
func ConfigureRepairRouter(e *echo.Group, node *core.DeltaNode) {

	repair := e.Group("/repair")

	repair.GET("/deal", func(c echo.Context) error {
		return nil
	})

	repair.GET("/commp", func(c echo.Context) error {
		return nil
	})

}
