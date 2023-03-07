package cmd

import (
	"delta/api"
	c "delta/config"
	"delta/core"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli/v2"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

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
			&cli.StringFlag{
				Name:  "stats-collection",
				Usage: "enable stats collection or not",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(Blue + "Starting Delta daemon..." + Reset)
			repo := c.String("repo")
			walletDir := c.String("wallet-dir")
			mode := c.String("mode")
			enableWebsocket := c.String("enable-websocket")
			statsCollection := c.String("stats-collection")

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

			if statsCollection == "" {
				cfg.Common.StatsCollection = true
			} else {
				cfg.Common.StatsCollection = false
			}

			fmt.Println("Setting up the whypfs node... ")
			fmt.Println("repo: ", Purple+repo+Reset)
			fmt.Println("walletDir: ", Purple+walletDir+Reset)
			fmt.Println("mode: ", Purple+cfg.Common.Mode+Reset)
			fmt.Println("enableWebsocket: ", cfg.Common.EnableWebsocket)
			fmt.Println("statsCollection: ", cfg.Common.StatsCollection)

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
			fmt.Println("Setting up the whypfs node... DONE")

			// set the node global meta
			fmt.Println("Computing the OS resources to use")
			core.ScanHostComputeResources(ln, repo)
			fmt.Println("Computing the OS resources to use... DONE")

			// run clean up
			fmt.Println("Running pre-start clean up")
			core.CleanUpContentAndPieceComm(ln)
			fmt.Println("Running pre-start clean up... DONE")

			// run the listeners
			fmt.Println("Subscribing the event listeners")
			core.SetLibp2pManagerSubscribe(ln)
			core.SetDataTransferEventsSubscribe(ln)
			fmt.Println("Subscribing the event listeners... DONE")

			// run the clean up every 30 minutes so we can retry and also remove the unecessary files on the blockstore.
			fmt.Println("Running the atomatic cron jobs")
			RunScheduledCleanupAndRetryCron(ln)
			fmt.Println("Running the atomatic cron jobs... DONE" + Reset)

			// launch the API node
			fmt.Println("----------------------------------")
			fmt.Println(Green + "Welcome! Delta daemon is running..." + Reset)
			fmt.Println("----------------------------------")
			fmt.Println(Purple + "Thank you for enabling stats collection. This helps us improve the product! If you don't want to share your stats, you can disable it by running the daemon with --stats-collection=false" + Reset)
			api.InitializeEchoRouterConfig(ln, *cfg)
			api.LoopForever()

			return nil
		},
	}

	daemonCommands = append(daemonCommands, daemonCmd)

	return daemonCommands

}

// RunScheduledCleanupAndRetryCron Run the cron jobs.
// The cron jobs are run every 12 hours and are responsible for cleaning up the database and the blockstore.
// It also retries the failed transfers.
// `RunScheduledCleanupAndRetryCron` is a function that runs a cron job on a node
func RunScheduledCleanupAndRetryCron(ln *core.DeltaNode) {

	maxCleanUpJobs := ln.Config.Dispatcher.MaxCleanupWorkers

	s := gocron.NewScheduler()
	s.Every(30).Minutes().Do(func() { // let's clean and retry every 30 minutes. It'll only get the old data.
		dispatcher := core.CreateNewDispatcher()
		//dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln))
		//dispatcher.AddJob(jobs.NewRetryProcessor(ln))
		dispatcher.Start(maxCleanUpJobs)

		core.ScanHostComputeResources(ln, ln.Node.Config.Blockstore)
	})

	s.Start()

}
