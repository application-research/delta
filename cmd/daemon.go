package cmd

import (
	"context"
	"fc-deal-making-service/api"
	"fc-deal-making-service/core"
	"fc-deal-making-service/jobs"
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
	pieceCommpJobFreq, err := strconv.Atoi(viper.Get("PIECE_COMMP_JOB_FREQ").(string))
	replicationJobFreq, err := strconv.Atoi(viper.Get("REPLICATION_JOB_FREQ").(string))

	if err != nil {
		pieceCommpJobFreq = 10
		replicationJobFreq = 10
	}

	pieceCommpJobFreqTick := time.NewTicker(time.Duration(pieceCommpJobFreq) * time.Second)
	replicationJobFreqTick := time.NewTicker(time.Duration(replicationJobFreq) * time.Second)

	for {
		select {
		case <-pieceCommpJobFreqTick.C:
			go func() {
				pieceCommpRun := jobs.NewPieceCommpProcessor(ln)
				go pieceCommpRun.Run()
			}()
		case <-replicationJobFreqTick.C:
			go func() {
				replicationRun := jobs.NewStorageDealMakerProcessor(ln)
				go replicationRun.Run()

			}()
		}
	}
}
