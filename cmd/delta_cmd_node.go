package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

var CLIConnectFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "delta-api",
		Usage:   "API connection info",
		EnvVars: []string{"DELTA_API"},
		Hidden:  true,
	},
	&cli.StringFlag{
		Name:    "delta-auth",
		Usage:   "delta auth token",
		EnvVars: []string{"DELTA_AUTH"},
		Hidden:  true,
	},
}

type DeltaCmdNode struct {
	DeltaApi  string
	DeltaAuth string
}

// NewDeltaCmdNode It creates a new DeltaCmdNode struct, which is a struct that contains the Delta API URL and the Delta API auth token
func NewDeltaCmdNode(c *cli.Context) (*DeltaCmdNode, error) {
	deltaApi := getFlagOrEnvVar(c, "delta-api", "DELTA_API", "http://localhost:1414")
	deltaAuth := getFlagOrEnvVar(c, "delta-auth", "DELTA_AUTH", "")

	if deltaAuth == "" {
		return nil, fmt.Errorf("DELTA_AUTH env variable or --delta-auth flag is required")
	}

	if err := healthCheck(deltaApi, deltaAuth); err != nil {
		return nil, fmt.Errorf("unable to communicate with delta daemon: %s", err)
	}

	return &DeltaCmdNode{
		DeltaApi:  deltaApi,
		DeltaAuth: deltaAuth,
	}, nil
}

// If the flag is set, use it. If not, check the environment variable. If that's not set, use the default value
func getFlagOrEnvVar(c *cli.Context, flagName, envVarName, defaultValue string) string {
	value := c.String(flagName)
	if value == "" {
		value = os.Getenv(envVarName)
		if value == "" {
			value = defaultValue
		}
	}
	return value
}

// It constructs an HTTP request to the `/open/node/info` endpoint, sets the `Authorization` header to the value of the
// `authKey` parameter, and then makes the request. If the response status code is not 200, it returns an error
func healthCheck(url string, authKey string) error {
	req, err := http.NewRequest("GET", url+"/health/check/auth/ping", nil)
	if err != nil {
		return fmt.Errorf("could not construct http request %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+authKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("could not make http request %s", err)
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf(string(body))
	}

	return err
}
