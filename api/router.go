package api

import (
	"context"
	"delta/config"
	"delta/core"
	"delta/utils"
	"encoding/json"
	"fmt"
	"github.com/application-research/delta-db/messaging"
	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/xerrors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
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

// @host node.delta.store
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
	e.Use(middleware.RateLimiter(
		middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
			Rate: 50, Burst: 200, ExpiresIn: 5 * time.Minute,
		}),
	))
	e.Use(middleware.SecureWithConfig(
		middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			ContentSecurityPolicy: "default-src 'self'",
		}),
	)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, span := otel.Tracer("GlobalRouterRequest").Start(context.Background(), "GlobalRouterRequest")
			defer span.End()
			span.SetName("Request: " + c.Request().Method + " " + c.Path())
			span.SetAttributes(attribute.String("user-agent", c.Request().UserAgent()))
			span.SetAttributes(attribute.String("path", c.Path()))
			span.SetAttributes(attribute.String("method", c.Request().Method))
			span.SetAttributes(attribute.String("remote_ip", c.RealIP()))
			span.SetAttributes(attribute.String("host", c.Request().Host))
			span.SetAttributes(attribute.String("referer", c.Request().Referer()))
			span.SetAttributes(attribute.String("request_uri", c.Request().RequestURI))
			ip := DeltaNodeConfig.Node.AnnounceAddrIP

			fmt.Println("Request: " + c.Request().Method + " " + c.Path() + " " + c.Request().UserAgent() + " " + c.RealIP() + " " + c.Request().Host + " " + c.Request().Referer() + " " + c.Request().RequestURI)
			s := struct {
				RemoteIP string `json:"remote_ip"`
				PublicIP string `json:"public_ip"`
				Host     string `json:"host"`
				Referer  string `json:"referer"`
				Request  string `json:"request"`
				Path     string `json:"path"`
			}{
				RemoteIP: c.RealIP(),
				PublicIP: ip,
				Host:     c.Request().Host,
				Referer:  c.Request().Referer(),
				Request:  c.Request().RequestURI,
				Path:     c.Path(),
			}
			b, err := json.Marshal(s)
			if err != nil {
				log.Error(err)
			}

			// log all errors so we can pre-emptively fix them
			utils.GlobalDeltaDataReporter.TraceLog(
				messaging.LogEvent{
					LogEventType:   "Route: " + core.GetHostname() + " " + c.Request().Method + " " + c.Path(),
					SourceHost:     core.GetHostname(),
					SourceIP:       ip,
					LogEventObject: b,
					LogEvent:       c.Path(),
					DeltaUuid:      config.Node.InstanceUuid,
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				})
			return next(c)
		}
	})

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

	ip := DeltaNodeConfig.Node.AnnounceAddrIP

	// get the request body and log it

	s := struct {
		RemoteIP     string `json:"remote_ip"`
		PublicIP     string `json:"public_ip"`
		Host         string `json:"host"`
		Referer      string `json:"referer"`
		Request      string `json:"request"`
		Path         string `json:"path"`
		ErrorDetails string `json:"details"`
	}{
		RemoteIP:     c.RealIP(),
		PublicIP:     ip,
		Host:         c.Request().Host,
		Referer:      c.Request().Referer(),
		Request:      c.Request().RequestURI,
		Path:         c.Path(),
		ErrorDetails: err.Error(),
	}

	b, errM := json.Marshal(s)
	if errM != nil {
		log.Error(errM)
	}

	// It's sending the error to the log server.
	utils.GlobalDeltaDataReporter.TraceLog(
		messaging.LogEvent{
			LogEventType:   "Error: " + core.GetHostname() + " " + c.Request().Method + " " + c.Path(),
			SourceHost:     core.GetHostname(),
			SourceIP:       ip,
			LogEventObject: b,
			LogEvent:       c.Path(),
			DeltaUuid:      DeltaNodeConfig.Node.InstanceUuid,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})

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
