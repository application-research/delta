package cmd
//
//import (
//	"delta/api"
//	"delta/core"
//	"github.com/urfave/cli/v2"
//)
//
//func CommpCmd() []*cli.Command {
//	// add a command to run API node
//	var commpCommands []*cli.Command
//
//	commpCmd := &cli.Command{
//		Name: "commp",
//
//		Flags: []cli.Flag{
//			&cli.StringFlag{
//				Name:  "file",
//				Usage: "specify the repo blockstore path of the node. ",
//			},
//		},
//		Action: func(c *cli.Context) error {
//
//			repo := c.String("file")
//
//			if repo == "" {
//				repo = ".whypfs"
//			}
//
//			if walletDir == "" {
//				walletDir = "./wallet"
//			}
//
//			// create the node (with whypfs, db, filclient)
//			nodeParams := core.NewLightNodeParams{
//				Repo:             repo,
//				DefaultWalletDir: walletDir,
//			}
//			ln, err := core.NewLightNode(context.Background(), nodeParams)
//			if err != nil {
//				return err
//			}
//
//			//	launch the dispatchers.
//			go runProcessors(ln)
//
//			//	launch clean up dispatch jobs
//			//	any failures due to the node shutdown will be retried after 1 day
//			go runCron(ln)
//
//			// launch the API node
//			api.InitializeEchoRouterConfig(ln)
//			api.LoopForever()
//
//			return nil
//		},
//	}
//
//	daemonCommands = append(daemonCommands, daemonCmd)
//
//	return daemonCommands
//
//}
