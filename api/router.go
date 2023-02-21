package api

import (
	"delta/core"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/xerrors"
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

type AuthResponse struct {
	Result struct {
		Validated bool   `json:"validated"`
		Details   string `json:"details"`
	} `json:"result"`
}

func InitializeEchoRouterConfig(ln *core.DeltaNode) {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.HTTPErrorHandler = ErrorHandler

	apiGroup := e.Group("/api/v1")
	openApiGroup := e.Group("/open")
	adminApiGroup := e.Group("/admin")
	ConfigureAdminRouter(adminApiGroup, ln)

	//apiGroup.Use(func(echo.HandlerFunc) echo.HandlerFunc {
	//	// check if upload of this node is disabled.
	//	if ln.MetaInfo.DisableRequest {
	//		return func(c echo.Context) error {
	//			return c.JSON(http.StatusForbidden, HttpErrorResponse{
	//				Error: HttpError{
	//					Code:    http.StatusForbidden,
	//					Reason:  http.StatusText(http.StatusForbidden),
	//					Details: "upload is disabled - due to memory or disk space limits",
	//				},
	//			})
	//		}
	//	}
	//	return nil
	//})

	apiGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationString := c.Request().Header.Get("Authorization")
			authParts := strings.Split(authorizationString, " ")

			response, err := http.Post(
				"https://estuary-auth-api.onrender.com/check-api-key",
				"application/json",
				strings.NewReader(fmt.Sprintf(`{"token": "%s"}`, authParts[1])),
			)

			if err != nil {
				log.Errorf("handler error: %s", err)
				return c.JSON(http.StatusInternalServerError, HttpErrorResponse{
					Error: HttpError{
						Code:    http.StatusInternalServerError,
						Reason:  http.StatusText(http.StatusInternalServerError),
						Details: err.Error(),
					},
				})
			}

			authResp, err := GetAuthResponse(response)
			if err != nil {
				log.Errorf("handler error: %s", err)
				return c.JSON(http.StatusInternalServerError, HttpErrorResponse{
					Error: HttpError{
						Code:    http.StatusInternalServerError,
						Reason:  http.StatusText(http.StatusInternalServerError),
						Details: err.Error(),
					},
				})
			}

			if authResp.Result.Validated == false {
				return c.JSON(http.StatusUnauthorized, HttpErrorResponse{
					Error: HttpError{
						Code:    http.StatusUnauthorized,
						Reason:  http.StatusText(http.StatusUnauthorized),
						Details: authResp.Result.Details,
					},
				})
			}
			if authResp.Result.Validated == true {
				return next(c)
			}
			return next(c)
		}
	})

	ConfigMetricsRouter(apiGroup)
	DealRouter(apiGroup, ln)
	ConfigureStatsCheckRouter(apiGroup, ln)
	ConfigureNodeInfoRouter(openApiGroup, ln)
	ConfigureOpenStatsCheckRouter(openApiGroup, ln)
	ConfigureRepairRouter(apiGroup, ln)
	ConfigureMinerRouter(apiGroup, ln)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:1414")) // configuration
}

func GetAuthResponse(resp *http.Response) (AuthResponse, error) {

	jsonBody := AuthResponse{}
	err := json.NewDecoder(resp.Body).Decode(&jsonBody)
	if err != nil {

		log.Error("empty json body")
		return AuthResponse{
			Result: struct {
				Validated bool   `json:"validated"`
				Details   string `json:"details"`
			}{
				Validated: false,
				Details:   "empty json body",
			},
		}, nil
	}

	return jsonBody, nil
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
