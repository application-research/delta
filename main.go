// It creates a new Echo instance, adds some middleware, creates a new WhyPFS node, creates a new GatewayHandler, and then
// adds a route to the Echo instance
package main

import (
	"delta/cmd"
	c "delta/config"
	"delta/core"
	"delta/utils"
	"fmt"
	"github.com/application-research/delta-db/event_models"
	"github.com/application-research/delta-db/messaging"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	_ "net/http"
	"os"
	"runtime"
	"time"
)

var (
	log = logging.Logger("api")
)

func main() {

	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Hostname:", core.GetHostname())

	ip, err := core.GetPublicIP()
	if err != nil {
		fmt.Println("Error getting public IP:", err)
	}
	utils.GlobalDeltaDataReporter.Trace(messaging.DeltaMetricsBaseMessage{
		ObjectType: "DeltaStartupLogs",
		Object: event_models.DeltaStartupLogs{
			NodeInfo:  core.GetHostname(),
			OSDetails: runtime.GOARCH + " " + runtime.GOOS,
			IPAddress: ip,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	})

	// get the config
	cfg := c.InitConfig()

	// get all the commands
	var commands []*cli.Command

	// commands
	commands = append(commands, cmd.DaemonCmd(&cfg)...)
	//commands = append(commands, cmd.CommpCmd()...)
	app := &cli.App{
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
