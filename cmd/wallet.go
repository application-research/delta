package cmd

import (
	c "delta/config"
	"github.com/urfave/cli/v2"
)

// TODO: add a command to manage wallet via CLI
func WalletCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var walletCommands []*cli.Command

	walletCmd := &cli.Command{
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
	walletCommands = append(walletCommands, walletCmd)

	return walletCommands
}
