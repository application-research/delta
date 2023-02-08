package cmd

import (
	"context"
	"fc-deal-making-service/api"
	"fc-deal-making-service/core"
	"fc-deal-making-service/jobs"
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"strconv"
	"time"
)

func DaemonCmd() []*cli.Command {
	// add a command to run API node
	var daemonCommands []*cli.Command

	daemonCmd := &cli.Command{
		Name:  "daemon",
		Usage: "A light version of Estuary that allows users to upload and download data from the Filecoin network.",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "repo",
			},
		},
		Action: func(c *cli.Context) error {

			repo := c.String("repo")

			if repo == "" {
				repo = ".whypfs"
			}

			ln, err := core.NewLightNode(context.Background(), repo)
			if err != nil {
				return err
			}

			//	launch the jobs
			go runProcessors(ln)
			go runRequeue(ln)
			go runMinerCheck(ln)
			// launch the API node
			api.InitializeEchoRouterConfig(ln)
			api.LoopForever()

			return nil
		},
	}

	// add commands.
	daemonCommands = append(daemonCommands, daemonCmd)

	return daemonCommands

}

func runProcessors(ln *core.LightNode) {

	// run the job every 10 seconds.
	jobDispatch, err := strconv.Atoi(viper.Get("DISPATCH_JOBS_EVERY").(string))
	jobDispatchWorker, err := strconv.Atoi(viper.Get("MAX_DISPATCH_WORKERS").(string))
	//pieceCommpJobFreq, err := strconv.Atoi(viper.Get("PIECE_COMMP_JOB_FREQ").(string))
	//replicationJobFreq, err := strconv.Atoi(viper.Get("REPLICATION_JOB_FREQ").(string))
	//minerCheckJobFreq, err := strconv.Atoi(viper.Get("MINER_INFO_UPDATE_JOB_FREQ").(string))

	if err != nil {
		jobDispatch = 10
	}

	jobDispatchTick := time.NewTicker(time.Duration(jobDispatch) * time.Second)
	//pieceCommpJobFreqTick := time.NewTicker(time.Duration(pieceCommpJobFreq) * time.Second)
	//replicationJobFreqTick := time.NewTicker(time.Duration(replicationJobFreq) * time.Second)
	//minerCheckJobFreqTick := time.NewTicker(time.Duration(minerCheckJobFreq) * time.Second)

	for {
		select {
		case <-jobDispatchTick.C:
			go func() {
				ln.Dispatcher.Start(jobDispatchWorker)
				for {
					if ln.Dispatcher.Finished() {
						fmt.Printf("All jobs finished.\n")
						break
					}
				}
			}()
			//case <-replicationJobFreqTick.C:
			//	go func() {
			//		replicationRun := jobs.NewStorageDealMakerProcessor(ln)
			//		go replicationRun.Run()
			//
			//	}()
			//case <-minerCheckJobFreqTick.C:
			//	go func() {
			//		minerCheckRun := jobs.NewMinerCheckProcessor(ln)
			//		go minerCheckRun.Run()
			//
			//	}()
		}

	}
}
func runRequeue(ln *core.LightNode) {
	// get all the pending content jobs. we need to requeue them.

	var contents []core.Content
	ln.DB.Model(&core.Content{}).Where("status = ?", "pinned").Find(&contents)

	for _, content := range contents {
		ln.Dispatcher.AddJob(jobs.NewItemContentProcessor(ln, content))
	}

	var pieceCommps []core.PieceCommitment
	ln.DB.Model(&core.PieceCommitment{}).Where("status = ?", "open").Find(&pieceCommps)

	for _, pieceCommp := range pieceCommps {
		var content core.Content
		ln.DB.Model(&core.Content{}).Where("piece_commitment_id = ?", pieceCommp.ID).Find(&content)
		ln.Dispatcher.AddJob(jobs.NewItemReplicationProcessor(ln, content, pieceCommp))
	}

	var contentsForDeletion []core.Content
	ln.DB.Model(&core.Content{}).Where("status = ?", "replication-complete").Find(&contentsForDeletion)

	for _, content := range contentsForDeletion {
		ln.Dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln, content)) // just delete it
	}

}

func runMinerCheck(ln *core.LightNode) {
	jobs.NewMinerCheckProcessor(ln).Run()
}
