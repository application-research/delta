package cmd

import (
	"context"
	"delta/api"
	c "delta/config"
	"delta/core"
	"delta/jobs"
	"fmt"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli/v2"
)

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
				Usage: "specify the wallet directory path of the node. ",
			},
		},
		Action: func(c *cli.Context) error {

			repo := c.String("repo")
			walletDir := c.String("wallet-dir")

			if repo == "" {
				repo = ".whypfs"
			}

			if walletDir == "" {
				walletDir = "./wallet"
			}

			// create the node (with whypfs, db, filclient)
			nodeParams := core.NewLightNodeParams{
				Repo:             repo,
				DefaultWalletDir: walletDir,
				Config:           cfg,
			}
			ln, err := core.NewLightNode(context.Background(), nodeParams)
			if err != nil {
				return err
			}

			//	launch the dispatchers.
			go runProcessors(ln)

			//	launch clean up dispatch jobs
			//	any failures due to the node shutdown will be retried after 1 day
			go runCron(ln)

			// launch the API node
			api.InitializeEchoRouterConfig(ln)
			api.LoopForever()

			return nil
		},
	}

	daemonCommands = append(daemonCommands, daemonCmd)

	return daemonCommands

}

// Run the cron jobs.
// The cron jobs are run every 12 hours and are responsible for cleaning up the database and the blockstore.
// It also retries the failed tranfers.
func runCron(ln *core.DeltaNode) {

	maxCleanUpJobs := ln.Config.Dispatcher.MaxCleanupWorkers

	s := gocron.NewScheduler()
	s.Every(12).Hour().Do(func() {
		dispatcher := core.CreateNewDispatcher()
		dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln))
		dispatcher.AddJob(jobs.NewRetryProcessor(ln))
		dispatcher.AddJob(jobs.NewMinerCheckProcessor(ln))
		dispatcher.Start(maxCleanUpJobs) // fix 100 workers for now.
	})

	s.Start()

}

func runProcessors(ln *core.DeltaNode) {

	// run the job every 10 seconds.
	jobDispatch := ln.Config.Dispatcher.DispatchJobsEvery
	jobDispatchWorker := ln.Config.Dispatcher.MaxDispatchWorkers

	jobDispatchTick := time.NewTicker(time.Duration(jobDispatch) * time.Second)

	for {
		select {
		case <-jobDispatchTick.C:
			go func() {
				ln.Dispatcher.AddJob(jobs.NewDataTransferStatusListenerProcessor(ln))
				ln.Dispatcher.Start(jobDispatchWorker)
				for {
					if ln.Dispatcher.Finished() {
						fmt.Printf("All jobs finished.\n")
						break
					}
				}
			}()
		}
	}
}
