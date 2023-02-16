package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/application-research/filclient"
	"github.com/application-research/whypfs-core"
	"github.com/urfave/cli/v2"
	"os"
)

func CommpCmd() []*cli.Command {
	// add a command to run API node
	var commpCommands []*cli.Command

	commpCmd := &cli.Command{
		Name: "commp",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "specify the repo blockstore path of the node. ",
			},
		},
		Action: func(c *cli.Context) error {
			file := c.String("file")

			params := whypfs.NewNodeParams{
				Ctx:       context.Background(),
				Datastore: whypfs.NewInMemoryDatastore(),
			}
			node, err := whypfs.NewNode(params)
			openFile, err := os.Open(file)
			reader := bufio.NewReader(openFile)
			if err != nil {
				fmt.Println(err)
				return err
			}

			fileNode, err := node.AddPinFile(context.Background(), reader, nil)
			if err != nil {
				fmt.Println(err)
				return err
			}

			commp, payloadSize, unpadddedPiece, err := filclient.GeneratePieceCommitment(context.Background(), fileNode.Cid(), node.Blockstore)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("commp: ", commp)
			fmt.Println("payloadSize: ", payloadSize)
			fmt.Println("unpadddedPiece: ", unpadddedPiece)

			return nil
		},
	}

	commpCommands = append(commpCommands, commpCmd)

	return commpCommands

}
