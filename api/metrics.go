package api

import (
	"delta/metrics"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Configuring the metrics router.
func ConfigMetricsRouter(e *echo.Group) {
	// metrics
	phandle := promhttp.Handler()
	e.GET("/debug/metrics/prometheus", func(e echo.Context) error {
		phandle.ServeHTTP(e.Response().Writer, e.Request())

		return nil
	})

	e.GET("/debug/metrics", func(e echo.Context) error {
		return e.JSON(http.StatusOK, "Ok")
		//return nil
	})

	e.GET("/debug/metrics", func(e echo.Context) error {
		metrics.Exporter().ServeHTTP(e.Response().Writer, e.Request())
		return nil
	})
	e.GET("/debug/stack", func(e echo.Context) error {
		err := metrics.WriteAllGoroutineStacks(e.Response().Writer)
		if err != nil {
			log.Error(err)
		}
		return err
	})

	e.GET("/debug/pprof/:prof", metrics.ServeProfile) // Upload for testing
}
