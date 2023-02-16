// It creates a new Echo instance, adds some middleware, creates a new WhyPFS node, creates a new GatewayHandler, and then
// adds a route to the Echo instance
package main

import (
	"delta/cmd"
	_ "net/http"
	"os"

	c "delta/config"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
)

var (
	log = logging.Logger("api")
)

func main() {
	cfg := c.InitConfig()

	// get all the commands
	var commands []*cli.Command

	// commands
	commands = append(commands, cmd.DaemonCmd(&cfg)...)
	commands = append(commands, cmd.CommpCmd()...)
	app := &cli.App{
		Commands: commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
