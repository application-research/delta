package cmd

import (
	"delta/api"
	c "delta/config"
	"delta/core"
	model "github.com/application-research/delta-db/db_models"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli/v2"
	"runtime"
	"syscall"
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
			meta := setGlobalNodeMeta(ln, repo)
			ln.MetaInfo = meta

			core.SetLibp2pManagerSubscribe(ln)
			core.SetDataTransferEventsSubscribe(ln)
			runScheduledCron(ln)

			// launch the API node
			api.InitializeEchoRouterConfig(ln, *cfg)
			api.LoopForever()

			return nil
		},
	}

	daemonCommands = append(daemonCommands, daemonCmd)

	return daemonCommands

}

// Run the cron jobs.
// The cron jobs are run every 12 hours and are responsible for cleaning up the database and the blockstore.
// It also retries the failed transfers.
// `runScheduledCron` is a function that runs a cron job on a node
func runScheduledCron(ln *core.DeltaNode) {

	maxCleanUpJobs := ln.Config.Dispatcher.MaxCleanupWorkers

	s := gocron.NewScheduler()
	s.Every(12).Hour().Do(func() {
		dispatcher := core.CreateNewDispatcher()
		//dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln))
		//dispatcher.AddJob(jobs.NewRetryProcessor(ln))
		dispatcher.Start(maxCleanUpJobs)
	})

	s.Start()

}

// Setting the global node meta.
// > This function sets the global node metadata for the given node
func setGlobalNodeMeta(ln *core.DeltaNode, repo string) *model.InstanceMeta {

	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	totalMemory := memStats.Sys
	totalMemory80 := totalMemory * 90 / 100

	// get the 80% of the total disk usage
	var stat syscall.Statfs_t
	syscall.Statfs(repo, &stat)
	totalStorage := stat.Blocks * uint64(stat.Bsize)
	totalStorage90 := totalStorage * 90 / 100

	// delete all data from the instance meta table
	ln.DB.Model(&model.InstanceMeta{}).Delete(&model.InstanceMeta{}, "id > ?", 0)
	// re-create
	instanceMeta := &model.InstanceMeta{
		MemoryLimit:                      totalMemory80,
		StorageLimit:                     totalStorage90,
		DisableRequest:                   false,
		DisableCommitmentPieceGeneration: false,
		DisableStorageDeal:               false,
		DisableOnlineDeals:               false,
		DisableOfflineDeals:              false,
	}
	ln.DB.Model(&model.InstanceMeta{}).Create(instanceMeta)

	return instanceMeta

}
