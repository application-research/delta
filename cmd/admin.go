package cmd

import (
	c "delta/config"
	"github.com/urfave/cli/v2"
)

// TODO: add a command to administer the node
func AdminCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var adminCommands []*cli.Command

	adminCmd := &cli.Command{
		Name:  "wallet",
		Usage: "Make a delta storage deal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "register-hex",
				Usage: "content to store",
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "The API key to use for the request",
			},
		},
	}
	adminCommands = append(adminCommands, adminCmd)

	return adminCommands
}
