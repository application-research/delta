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

func runCron(ln *core.LightNode) {

	gocron.Every(1).Day().Do(func() {
		dispatcher := core.CreateNewDispatcher()
		dispatcher.AddJob(jobs.NewItemContentCleanUpProcessor(ln))
		dispatcher.AddJob(jobs.NewRetryProcessor(ln))
		dispatcher.Start(1000)
	})

	<-gocron.Start()

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
