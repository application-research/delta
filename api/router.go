package api

import (
	"delta/config"
	"delta/core"
	_ "delta/docs/swagger"
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
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/xerrors"
)

var (
	OsSignal        chan os.Signal
	log             = logging.Logger("router")
	DeltaNode       *core.DeltaNode
	DeltaNodeConfig *config.DeltaConfig
)

// HttpError A struct that contains a struct that contains a bool and a string.
// @property {int} Code - The HTTP status code
// @property {string} Reason - The reason for the error.
// @property {string} Details - The details of the error.
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

// InitializeEchoRouterConfig Initializing the router.
// It's initializing the Echo router, and configuring the routes for the API
// @title Delta API
// @description This is the API for the Delta application.
// @termsOfService http://delta.store

// @contact.name API Support

// @license.name Apache 2.0 Apache-2.0 OR MIT

// @host localhost:1414
// @BasePath  /
// @securityDefinitions.Bearer
// @securityDefinitions.Bearer.type apiKey
// @securityDefinitions.Bearer.in header
// @securityDefinitions.Bearer.name Authorization
func InitializeEchoRouterConfig(ln *core.DeltaNode, config config.DeltaConfig) {

	DeltaNode = ln
	DeltaNodeConfig = &config

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.SecureWithConfig(
		middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			ContentSecurityPolicy: "default-src 'self'",
		}),
	)

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.Pre(middleware.RemoveTrailingSlash())
	e.HTTPErrorHandler = ErrorHandler
	apiGroup := e.Group("/api/v1")
	openApiGroup := e.Group("/open")
	adminApiGroup := e.Group("/admin")

	// health check api group
	healthCheckApiGroup := e.Group("/health")

	// debug profiler group
	debugGroup := e.Group("/debug")
	ConfigureDebugProfileRouter(debugGroup, ln)

	// Authentication
	apiGroup.Use(Authenticate(config))
	adminApiGroup.Use(Authenticate(config))

	// health check
	ConfigureHealthCheckRouter(healthCheckApiGroup, ln)

	// admin api
	ConfigureAdminRouter(adminApiGroup, ln)

	// protected
	ConfigureDealRouter(apiGroup, ln)
	ConfigureStatsCheckRouter(apiGroup, ln)
	ConfigureRepairRouter(apiGroup, ln)

	// open api
	ConfigureNodeInfoRouter(openApiGroup, ln)
	ConfigureOpenStatsCheckRouter(openApiGroup, ln)
	ConfigureOpenInfoCheckRouter(openApiGroup, ln)

	// metrics
	ConfigMetricsRouter(openApiGroup)

	// It's checking if the websocket is enabled.
	if config.Common.EnableWebsocket {
		// websocket
		fmt.Println("Websocket enabled")
		ws := core.NewWebsocketService(ln)
		go ws.HandlePieceCommitmentMessages()
		go ws.HandleContentDealMessages()
		go ws.HandleContentMessages()
		ConfigureWebsocketRouter(openApiGroup, ln)
	}

	// swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:1414")) // configuration
}

func Authenticate(config config.DeltaConfig) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// check if the authorization header is present
			authorizationString := c.Request().Header.Get("Authorization")
			authParts := strings.Split(authorizationString, " ")
			// validate the auth parts
			if len(authParts) != 2 {
				return c.JSON(http.StatusUnauthorized, HttpErrorResponse{
					Error: HttpError{
						Code:    http.StatusUnauthorized,
						Reason:  http.StatusText(http.StatusUnauthorized),
						Details: "Invalid authorization header",
					},
				})
			}
			if authParts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, HttpErrorResponse{
					Error: HttpError{
						Code:    http.StatusUnauthorized,
						Reason:  http.StatusText(http.StatusUnauthorized),
						Details: "Invalid authorization header",
					},
				})
			}
			if authParts[1] == "" {
				return c.JSON(http.StatusUnauthorized, HttpErrorResponse{
					Error: HttpError{
						Code:    http.StatusUnauthorized,
						Reason:  http.StatusText(http.StatusUnauthorized),
						Details: "Invalid authorization header",
					},
				})
			}
			if config.Common.Mode == "standalone" {
				// check if the API key is present and is valid
				if config.Standalone.APIKey != authParts[1] {
					return c.JSON(http.StatusUnauthorized, HttpErrorResponse{
						Error: HttpError{
							Code:    http.StatusUnauthorized,
							Reason:  http.StatusText(http.StatusUnauthorized),
							Details: "Invalid API key provided for standalone mode",
						},
					})
				}
			}
			// if everything is good. we can check the token against estuary-auth.
			response, err := http.Post(
				config.ExternalApis.AuthSvcApi+"/check-api-key",
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
	}
}

// GetAuthResponse It's making a request to the auth API to check if the API key is valid.
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

// Usage `Echo#Pre(RemoveTrailingSlash())`
func ValidateRequestBody() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return nil
	}
}

// ErrorHandler It's a function that is called when an error occurs.
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
// It's a function that is called when an error occurs.
func LoopForever() {
	fmt.Printf("Entering infinite loop\n")
	signal.Notify(OsSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	_ = <-OsSignal
	fmt.Printf("Exiting infinite loop received OsSignal\n")
}
