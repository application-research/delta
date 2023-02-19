package cmd

import (
	"bufio"
	"context"
	"delta/core"
	"fmt"
	"github.com/application-research/filclient"
	"github.com/application-research/whypfs-core"
	"github.com/urfave/cli/v2"
	"os"
)

type CommpResult struct {
	Cid               string `json:"cid,omitempty"`
	Commp             string `json:"commp,omitempty"`
	PaddedPieceSize   uint64 `json:"padded-piece-size,omitempty"`
	UnPaddedPieceSize uint64 `json:"un-padded-piece-size,omitempty"`
}

func CommpCmd() []*cli.Command {
	// add a command to run API node
	var commpCommands []*cli.Command
	var commpService = new(core.CommpService)
	commpFileCmd := &cli.Command{
		Name: "commp",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "specify the file",
			},
		},
		Action: func(c *cli.Context) error {
			file := c.String("file")
			car := c.String("car")
			if file == "" {
				if car == "" {
					fmt.Println("file or car is required")
					return nil
				}
			}

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

			if file != "" {
				commp, payloadSize, unpadddedPiece, err := filclient.GeneratePieceCommitment(context.Background(), fileNode.Cid(), node.Blockstore)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("payloadcid: ", fileNode.Cid())
				fmt.Println("commp: ", commp)
				fmt.Println("payloadSize: ", payloadSize)
				fmt.Println("unpadddedPiece: ", unpadddedPiece)
				fmt.Println("paddedPiece: ", unpadddedPiece.Padded())
				return nil
			}

			return nil
		},
	}

	commpCarCmd := &cli.Command{
		Name: "commp-car",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "specify the car file",
			},
		},
		Action: func(c *cli.Context) error {
			var commpResult CommpResult
			car := c.String("file")
			openFile, err := os.Open(car)

			if err != nil {
				fmt.Println(err)
			}
			reader := bufio.NewReader(openFile)

			if car != "" {
				pieceInfo, err := commpService.GenerateCommPCarV2(reader)
				if err != nil {
					fmt.Println(err)
				}

				if err != nil {
					fmt.Println(err)
				}

				// return json to console.
				commpResult.Commp = pieceInfo.PieceCID.String()
				commpResult.PaddedPieceSize = uint64(pieceInfo.Size)
				commpResult.UnPaddedPieceSize = uint64(pieceInfo.Size.Unpadded())

			}
			return nil
		},
	}

	commpCommands = append(commpCommands, commpFileCmd)
	commpCommands = append(commpCommands, commpCarCmd)

	return commpCommands

}
