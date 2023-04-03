package cmd

import (
	"delta/api"
	c "delta/config"
	"delta/core"
	"delta/jobs"
	"delta/utils"
	"fmt"
	"github.com/application-research/delta-db/db_models"
	_ "github.com/application-research/delta-db/db_models"
	"github.com/application-research/delta-db/messaging"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli/v2"
	"runtime"
	"time"
)

// DaemonCmd Creating a new command called `daemon` that will run the API node.
func DaemonCmd(cfg *c.DeltaConfig) []*cli.Command {

	// add a command to run API node
	var daemonCommands []*cli.Command

	daemonCmd := &cli.Command{
		Name:  "daemon",
		Usage: "Run the delta daemon",
		Description: "The delta daemon is the main process that runs the delta node. It is responsible for " +
			"handling all the incoming requests and processing them. It also runs the background jobs " +
			"that are required for the node to function properly.",

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
				Name:  "commp-mode",
				Usage: "standalone or cluster mode",
				Value: utils.COMPP_MODE_FILBOOST,
			},
			&cli.BoolFlag{
				Name:  "enable-websocket",
				Usage: "enable websocket or not",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "stats-collection",
				Usage: "enable stats collection or not",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "keep-copies",
				Usage: "keep copies of the data or not",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "miner-throttle",
				Usage: "enable miner throttle or not",
				Value: false,
			},
		},

		Action: func(c *cli.Context) error {

			fmt.Println("OS:", runtime.GOOS)
			fmt.Println("Architecture:", runtime.GOARCH)
			fmt.Println("Hostname:", core.GetHostname())

			ip, err := core.GetPublicIP()
			if err != nil {
				fmt.Println("Error getting public IP:", err)
			}

			fmt.Println(utils.Blue + "Starting Delta daemon..." + utils.Reset)
			repo := c.String("repo")
			walletDir := c.String("wallet-dir")
			mode := c.String("mode")
			enableWebsocket := c.Bool("enable-websocket")
			statsCollection := c.Bool("stats-collection")
			commpMode := c.String("commp-mode")

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

			cfg.Common.EnableWebsocket = enableWebsocket
			cfg.Common.StatsCollection = statsCollection
			cfg.Common.CommpMode = commpMode

			fmt.Println(utils.Blue + "Setting up the whypfs node... " + utils.Reset)
			fmt.Println("repo: ", utils.Purple+repo+utils.Reset)
			fmt.Println("walletDir: ", utils.Purple+walletDir+utils.Reset)
			fmt.Println("mode: ", utils.Purple+cfg.Common.Mode+utils.Reset)
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
			fmt.Println(utils.Blue + "Setting up the whypfs node... DONE" + utils.Reset)

			// set the node global meta
			fmt.Println(utils.Blue + "Computing the OS resources to use" + utils.Reset)
			core.ScanHostComputeResources(ln, repo)
			fmt.Println(utils.Blue + "Computing the OS resources to use... DONE" + utils.Reset)

			// run clean up
			fmt.Println(utils.Blue + "Running pre-start clean up" + utils.Reset)
			core.CleanUpContentAndPieceComm(ln)
			fmt.Println(utils.Blue + "Running pre-start clean up... DONE" + utils.Reset)

			// run the listeners
			fmt.Println(utils.Blue + "Subscribing the event listeners" + utils.Reset)
			core.SetLibp2pManagerSubscribe(ln)
			core.SetDataTransferEventsSubscribe(ln)
			fmt.Println(utils.Blue + "Subscribing the event listeners... DONE" + utils.Reset)

			// run the clean up every 30 minutes so we can retry and also remove the unecessary files on the blockstore.
			fmt.Println(utils.Blue + "Running the atomic cron jobs" + utils.Reset)
			RunScheduledCleanupAndRetryCron(ln)
			fmt.Println(utils.Blue + "Running the atomic cron jobs... DONE" + utils.Reset)

			// launch the API node
			fmt.Println(utils.Blue + "Starting Delta." + utils.Reset)
			fmt.Println(utils.Green + `

     %%%%%%%%/          %%%%%%%%%%%%%%% %%%%%     %%%%%%%%%%%%%%%%%     %%%%%%  
    @@@@@@@@@@@@@@@     @@@@@@@@@@@@@@ @@@@@      @@@@@@@@@@@@@@@@@   @@@@@@@@  
    @@@@@     @@@@@@@  @@@@@@          @@@@@           @@@@@         @@@@@@@@@@ 
   @@@@@@       @@@@@  @@@@@          @@@@@            @@@@@       @@@@@  @@@@@ 
   @@@@@        @@@@@ @@@@@@@@@@@@@@ (@@@@@           @@@@@       @@@@@   @@@@@ 
  @@@@@@       @@@@@@ @@@@@@@@@@@@@  @@@@@           /@@@@@      @@@@@    #@@@@,
  @@@@@       @@@@@@ @@@@@*         @@@@@@           @@@@@     @@@@@@@@@@@@@@@@@
 @@@@@@@@@@@@@@@@@   @@@@@@@@@@@@@@ @@@@@@@@@@@@@@  @@@@@@    @@@@@        @@@@@
 @@@@@@@@@@@@@@     @@@@@@@@@@@@@@ @@@@@@@@@@@@@@@  @@@@@    @@@@@         @@@@@

(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)(ᵔᴥᵔ)

By: Protocol Labs - Outercore Engineering
` + utils.Reset + utils.Red + "version: " + cfg.Common.Version + utils.Reset)

			fmt.Println(utils.Blue + "Reporting Delta startup logs" + utils.Reset)
			utils.GlobalDeltaDataReporter.Trace(messaging.DeltaMetricsBaseMessage{
				ObjectType: "DeltaStartupLogs",
				Object: db_models.DeltaStartupLogs{
					NodeInfo:      core.GetHostname(),
					OSDetails:     runtime.GOARCH + " " + runtime.GOOS,
					IPAddress:     ip,
					DeltaNodeUuid: cfg.Node.InstanceUuid,
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				},
			})
			fmt.Println(utils.Blue + "Reporting Delta startup logs... DONE" + utils.Reset)
			fmt.Println("----------------------------------")
			fmt.Println(utils.Green + "Welcome! Delta daemon is running..." + utils.Reset)
			fmt.Println(utils.Green + "You can check the documentation at: " + utils.Reset + utils.Purple + "https://github.com/application-research/delta/tree/main/docs" + utils.Reset)
			fmt.Println("----------------------------------")
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
	fmt.Println(utils.Purple + "Scheduling dispatchers and scanners..." + utils.Reset)
	maxCleanUpJobs := ln.Config.Dispatcher.MaxCleanupWorkers

	s := gocron.NewScheduler()
	s.Every(24).Hour().Do(func() { // let's clean and retry every 30 minutes. It'll only get the old data.
		dispatcher := core.CreateNewDispatcher()
		//dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln))
		dispatcher.AddJob(jobs.NewRetryProcessor(ln))
		dispatcher.Start(maxCleanUpJobs)

		core.CleanUpContentAndPieceComm(ln)
		core.ScanHostComputeResources(ln, ln.Node.Config.Blockstore)
	})

	s.Start()

}
