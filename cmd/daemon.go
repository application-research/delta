package cmd

import (
	"context"
	"delta/api"
	c "delta/config"
	"delta/core"
	"delta/core/model"
	"delta/jobs"
	"delta/utils"
	"fmt"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli/v2"
	"runtime"
	"strconv"
	"syscall"
	"time"
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
				Usage: "specify the wallet_estuary directory path of the node. ",
			},
		},
		Action: func(c *cli.Context) error {

			repo := c.String("repo")
			walletDir := c.String("wallet-dir")

			if repo == "" {
				repo = ".whypfs"
			}

			if walletDir == "" {
				walletDir = "./wallet_estuary"
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

			// set the node global meta
			meta := setGlobalNodeMeta(ln, repo)
			ln.MetaInfo = meta

			runLib2pManagerListener(ln)
			runScheduledCron(ln)

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
// It also retries the failed transfers.
func runScheduledCron(ln *core.DeltaNode) {

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

func runLib2pManagerListener(d *core.DeltaNode) {
	d.FilClient.Libp2pTransferMgr.Subscribe(func(dbid uint, fst filclient.ChannelState) {
		switch fst.Status {
		case datatransfer.Requested:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			d.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				TransferStarted: time.Now(),
			})
		case datatransfer.TransferFinished, datatransfer.Completed:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			transferId, err := strconv.Atoi(fst.TransferID)
			if err != nil {
				fmt.Println(err)
			}
			d.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				DealID:           int64(transferId),
				TransferFinished: time.Now(),
				SealedAt:         time.Now(),
				LastMessage:      utils.DEAL_STATUS_TRANSFER_FINISHED,
			})
			d.DB.Model(&model.Content{}).Where("id = (select content from content_deals cd where cd.id = ?)", dbid).Updates(model.Content{
				Status: utils.DEAL_STATUS_TRANSFER_FINISHED,
			})
		case datatransfer.Failed:
			fmt.Println("Transfer status: ", fst.Status, " for transfer id: ", fst.TransferID, " for db id: ", dbid)
			var contentDeal model.ContentDeal
			d.DB.Model(&model.ContentDeal{}).Where("id = ?", dbid).Updates(model.ContentDeal{
				FailedAt: time.Now(),
			}).Find(&contentDeal)
			d.DB.Model(&model.Content{}).Joins("left join content_deals as cd on cd.content = c.id").Where("cd.id = ?", dbid).Updates(model.Content{
				Status: utils.DEAL_STATUS_TRANSFER_FAILED,
			})

			d.Dispatcher.AddJobAndDispatch(jobs.NewDataTransferRestartProcessor(d, contentDeal), 1)
		default:
		}
	})
}

func setGlobalNodeMeta(ln *core.DeltaNode, repo string) *model.InstanceMeta {

	// get the 80% of the total memory usage
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
