package cmd

import (
	"context"
	"fc-deal-making-service/api"
	"fc-deal-making-service/core"
	"fc-deal-making-service/jobs"
	"fmt"
	"github.com/urfave/cli/v2"
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
				Name: "enable-api",
			},
		},
		Action: func(c *cli.Context) error {

			ln, err := core.NewLightNode(context.Background())
			if err != nil {
				return err
			}

			//	launch the jobs
			go runJobs(ln)

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

func runJobs(ln *core.LightNode) {

	fmt.Println("run jobs")
	// run the job every 10 seconds.
	tick := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tick.C:
			go func() {
				fmt.Println("piece commp processor")
				pieceNewCommpProcessor := jobs.NewPieceCommpProcessor(ln)
				pieceNewCommpProcessor.Run()
			}()

			go func() {
				fmt.Println("replication deal maker")
				dealMaker := jobs.NewReplicationProcessor(ln)
				dealMaker.Run()
			}()
		}
	}

}
