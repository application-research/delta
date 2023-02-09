package api

import (
	"delta/core"
	"github.com/labstack/echo/v4"
)

//	 repair deals (re-create or re-try)
func ConfigureRepairRouter(e *echo.Group, node *core.LightNode) {
	e.GET("/repair", func(c echo.Context) error {
		return nil
	})

}
