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

func NewDeltaCmdNode(c *cli.Context) (*DeltaCmdNode, error) {
	deltaApi := c.String("delta-api")
	if deltaApi == "" {
		deltaApi = os.Getenv("DELTA_API")
		if deltaApi == "" {
			deltaApi = "http://localhost:1414"
		}
	}

	deltaAuth := c.String("delta-auth")
	if deltaAuth == "" {
		deltaAuth = os.Getenv("DELTA_AUTH")
		if deltaAuth == "" {
			return nil, fmt.Errorf("DELTA_AUTH env variable or --delta-auth flag is required")
		}
	}

	err := healthCheck(deltaApi, deltaAuth)

	if err != nil {
		return nil, fmt.Errorf("unable to communicate with delta daemon: %s", err)
	}

	return &DeltaCmdNode{
		DeltaApi:  deltaApi,
		DeltaAuth: deltaAuth,
	}, nil
}

// Verify that DDM API is reachable
func healthCheck(url string, authKey string) error {
	req, err := http.NewRequest("GET", url+"/open/node/info", nil)
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
