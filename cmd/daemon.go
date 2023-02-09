package cmd

import (
	"context"
	"delta/api"
	"delta/core"
	"delta/jobs"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"strconv"
	"sync"
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

func runCron(ln *core.LightNode) {

	gocron.Every(1).Day().Do(func() {
		dispatcher := core.CreateNewDispatcher()
		dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln))
		dispatcher.AddJob(jobs.NewRetryProcessor(ln))
		dispatcher.Start(1000)
	})

	<-gocron.Start()

}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

func runProcessors(ln *core.LightNode) {

	// run the job every 10 seconds.
	jobDispatch, err := strconv.Atoi(viper.Get("DISPATCH_JOBS_EVERY").(string))
	jobDispatchWorker, err := strconv.Atoi(viper.Get("MAX_DISPATCH_WORKERS").(string))

	// 	cron job
	//	- retry
	//	- clean up

	if err != nil {
		jobDispatch = 10
	}

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
func runRequeue(ln *core.LightNode) {
	// get all the pending content jobs. we need to requeue them.

	var contents []core.Content
	ln.DB.Model(&core.Content{}).Where("status = ?", "pinned").Find(&contents)

	for _, content := range contents {
		ln.Dispatcher.AddJob(jobs.NewPieceCommpProcessor(ln, content))
	}

	var contentsForCommp []core.Content
	ln.DB.Model(&core.Content{}).Where("status = ?", "piece-computing").Find(&contentsForCommp)

	for _, content := range contentsForCommp {
		ln.Dispatcher.AddJob(jobs.NewPieceCommpProcessor(ln, content))
	}

	var pieceCommps []core.PieceCommitment
	ln.DB.Model(&core.PieceCommitment{}).Where("status = ?", "open").Find(&pieceCommps)

	for _, pieceCommp := range pieceCommps {
		var content core.Content
		ln.DB.Model(&core.Content{}).Where("piece_commitment_id = ?", pieceCommp.ID).Find(&content)
		ln.Dispatcher.AddJob(jobs.NewStorageDealMakerProcessor(ln, content, pieceCommp))
	}

	var contentsForDeletion []core.Content
	ln.DB.Model(&core.Content{}).Where("status = ?", "replication-complete").Find(&contentsForDeletion)

}
