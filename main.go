// It creates a new Echo instance, adds some middleware, creates a new WhyPFS node, creates a new GatewayHandler, and then
// adds a route to the Echo instance
package main

import (
	"delta/cmd"
	c "delta/config"
	_ "embed"
	"fmt"
	_ "net/http"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var (
	log = logging.Logger("api")
)

var Commit string
var Version string

// It initializes the config, gets all the commands, and runs the app.
func main() {

	// get the config
	cfg := c.InitConfig()
	cfg.Common.Commit = Commit
	cfg.Common.Version = Version

	// get all the commands
	var commands []*cli.Command

	// commands
	commands = append(commands, cmd.DaemonCmd(&cfg)...)
	commands = append(commands, cmd.CarCmd(&cfg)...)
	commands = append(commands, cmd.CommpCmd(&cfg)...)
	commands = append(commands, cmd.DealCmd(&cfg)...)
	commands = append(commands, cmd.StatusCmd(&cfg)...)
	commands = append(commands, cmd.WalletCmd(&cfg)...)

	app := &cli.App{
		Commands:    commands,
		Name:        "delta",
		Description: "A deal making engine microservice for the filecoin network",
		Version:     fmt.Sprintf("%s+git.%s\n", cfg.Common.Version, cfg.Common.Commit),
		Usage:       "delta [command] [arguments]",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
