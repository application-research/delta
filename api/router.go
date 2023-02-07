package api

import (
	"fc-deal-making-service/core"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	OsSignal chan os.Signal
	log      = logging.Logger("router")
)

type HttpError struct {
	Code    int    `json:"code,omitempty"`
	Reason  string `json:"reason"`
	Details string `json:"details"`
}

func (he HttpError) Error() string {
	if he.Details == "" {
		return he.Reason
	}
	return he.Reason + ": " + he.Details
}

type HttpErrorResponse struct {
	Error HttpError `json:"error"`
}

// RouterConfig configures the API node
func InitializeEchoRouterConfig(ln *core.LightNode) {
	// Echo instance
	e := echo.New()
	e.File("/", "templates/index.html")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.HTTPErrorHandler = ErrorHandler

	defaultGatewayRoute := e.Group("")
	ConfigureGatewayRouter(defaultGatewayRoute, ln) // access to light node

	apiGroup := e.Group("/api/v1")       // no protection for now
	ConfigureGatewayRouter(apiGroup, ln) // access to light node
	ConfigurePinningRouter(apiGroup, ln) // store

	ConfigMetricsRouter(apiGroup) // metrics

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:1313")) // configuration
}

func ErrorHandler(err error, c echo.Context) {
	var httpRespErr *HttpError
	if xerrors.As(err, &httpRespErr) {
		log.Errorf("handler error: %s", err)
		if err := c.JSON(httpRespErr.Code, HttpErrorResponse{Error: *httpRespErr}); err != nil {
			log.Errorf("handler error: %s", err)
			return
		}
		return
	}

	var echoErr *echo.HTTPError
	if xerrors.As(err, &echoErr) {
		if err := c.JSON(echoErr.Code, HttpErrorResponse{
			Error: HttpError{
				Code:    echoErr.Code,
				Reason:  http.StatusText(echoErr.Code),
				Details: echoErr.Message.(string),
			},
		}); err != nil {
			log.Errorf("handler error: %s", err)
			return
		}
		return
	}

	log.Errorf("handler error: %s", err)
	if err := c.JSON(http.StatusInternalServerError, HttpErrorResponse{
		Error: HttpError{
			Code:    http.StatusInternalServerError,
			Reason:  http.StatusText(http.StatusInternalServerError),
			Details: err.Error(),
		},
	}); err != nil {
		log.Errorf("handler error: %s", err)
		return
	}
}

// LoopForever on signal processing
func LoopForever() {
	fmt.Printf("Entering infinite loop\n")

	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	_ = <-OsSignal

	fmt.Printf("Exiting infinite loop received OsSignal\n")
}
