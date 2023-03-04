package cmd

import (
	"delta/api"
	c "delta/config"
	"delta/core"
	"github.com/urfave/cli/v2"
)

// DaemonCmd Creating a new command called `daemon` that will run the API node.
func DaemonCmd(cfg *c.DeltaConfig) []*cli.Command {

	// add a command to run API node
	var daemonCommands []*cli.Command

	daemonCmd := &cli.Command{
		Name:  "daemon",
		Usage: "A light version of Estuary that allows users to upload and download data from the Filecoin network.",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "repo",
				Usage: "specify the repo blockstore path of the node. ",
			},
			&cli.StringFlag{
				Name:  "wallet-dir",
				Usage: "specify the wallet_estuary directory path of the node. ",
			},

			&cli.StringFlag{
				Name:  "mode",
				Usage: "standalone or cluster mode",
			},

			&cli.StringFlag{
				Name:  "enable-websocket",
				Usage: "enable websocket or not",
			},
		},
		Action: func(c *cli.Context) error {

			repo := c.String("repo")
			walletDir := c.String("wallet-dir")
			mode := c.String("mode")
			enableWebsocket := c.String("enable-websocket")

			if repo == "" {
				repo = ".whypfs"
			}

			if walletDir == "" {
				walletDir = "./wallet_estuary"
			}

			if mode == "" {
				cfg.Common.Mode = "cluster"
			} else {
				cfg.Common.Mode = mode
			}

			if enableWebsocket == "" {
				cfg.Common.EnableWebsocket = false
			} else {
				cfg.Common.EnableWebsocket = true
			}

			// create the node (with whypfs, db, filclient)
			nodeParams := core.NewLightNodeParams{
				Repo:             repo,
				DefaultWalletDir: walletDir,
				Config:           cfg,
			}

			ln, err := core.NewLightNode(nodeParams)
			if err != nil {
				return err
			}

			// set the node global meta
			core.ScanHostComputeResources(ln, repo)

			// run clean up
			core.CleanUpContentAndPieceComm(ln)

			// run the listeners
			core.SetLibp2pManagerSubscribe(ln)
			core.SetDataTransferEventsSubscribe(ln)

			// run the clean up every 30 minutes so we can retry and also remove the unecessary files on the blockstore.
			core.RunScheduledCleanupAndRetryCron(ln)

			// launch the API node
			api.InitializeEchoRouterConfig(ln, *cfg)
			api.LoopForever()

			return nil
		},
	}

	daemonCommands = append(daemonCommands, daemonCmd)

	return daemonCommands

}
