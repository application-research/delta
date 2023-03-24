package cmd

import (
	c "delta/config"
	"github.com/urfave/cli/v2"
)

// TODO: add a command to make deals via CLI
func DealCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var dealCommands []*cli.Command

	dealCmd := &cli.Command{
		Name:  "deal",
		Usage: "Make a delta storage deal",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "content",
				Usage: "content to store",
			},
			&cli.StringFlag{
				Name:  "metadata",
				Usage: "metadata to store",
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "The API key to use for the request",
			},
		},
	}
	dealCommands = append(dealCommands, dealCmd)

	return dealCommands
}
