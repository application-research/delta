package api

import (
	"delta/core"
	httpprof "net/http/pprof"
	"runtime/pprof"
	"time"

	"github.com/labstack/echo/v4"
)

// ConfigureDebugRouter It configures the router to handle requests for node profiling
func ConfigureDebugProfileRouter(profilerGroup *echo.Group, node *core.DeltaNode) {

	profilerGroup.GET("/pprof/:prof", func(c echo.Context) error {
		httpprof.Handler(c.Param("prof")).ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	profilerGroup.GET("/cpuprofile", func(c echo.Context) error {
		if err := pprof.StartCPUProfile(c.Response()); err != nil {
			return err
		}

		defer pprof.StopCPUProfile()

		select {
		case <-c.Request().Context().Done():
			return c.Request().Context().Err()
		case <-time.After(time.Second * 30):
		}
		return nil
	})
}
