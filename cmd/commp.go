package cmd

import (
	"bufio"
	"bytes"
	"context"
	"delta/api"
	"delta/core"
	"encoding/json"
	"fmt"
	"github.com/application-research/filclient"
	"github.com/application-research/whypfs-core"
	"github.com/urfave/cli/v2"
	"io"
	"io/ioutil"
	"net/http"
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
			&cli.BoolFlag{
				Name:  "for-offline",
				Usage: "specify the car file",
				Value: true,
			},
			&cli.StringFlag{
				Name:  "delta-api-url",
				Usage: "specify the delta api url",
			},
			&cli.StringFlag{
				Name:  "delta-api-key",
				Usage: "Estuary API key",
			},
		},
		Action: func(c *cli.Context) error {
			var commpResult api.ContentMakeDealRequest
			car := c.String("file")
			forOffline := c.Bool("for-offline")
			miner := c.String("miner")
			//wallet := c.String("wallet")

			params := whypfs.NewNodeParams{
				Ctx:       context.Background(),
				Datastore: whypfs.NewInMemoryDatastore(),
			}

			node, err := whypfs.NewNode(params)
			openFile, err := os.Open(car)
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

			if miner != "" {
				commpResult.Miner = miner
			}

			if car != "" {
				fileNodeFromBs, err := node.GetFile(context.Background(), fileNode.Cid())
				pieceInfo, err := commpService.GenerateCommPCarV2(fileNodeFromBs)
				if err != nil {
					fmt.Println(err)
					return err
				}
				// return json to console.
				commpResult.Cid = fileNode.Cid().String()
				commpResult.PieceCommitment.Piece = pieceInfo.PieceCID.String()
				commpResult.PieceCommitment.PaddedPieceSize = uint64(pieceInfo.Size)
				commpResult.PieceCommitment.UnPaddedPieceSize = uint64(pieceInfo.Size.Unpadded())

				// if for offline, add connection mode offline
				if forOffline {
					commpResult.ConnectionMode = "offline"
				} else {
					commpResult.ConnectionMode = "online"
				}

				size, err := fileNode.Size()
				if err != nil {
					fmt.Println(err)
					return err
				}
				commpResult.Size = int64(size)

				var buffer bytes.Buffer
				err = PrettyEncode(commpResult, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())

				// if the delta api url and key is given, send the result to delta api.
				deltaApiUrl := c.String("delta-api-url")
				deltaApiKey := c.String("delta-api-key")

				if deltaApiUrl != "" && deltaApiKey != "" {
					// send the result to delta api.
					client := &http.Client{}
					fmt.Println(deltaApiUrl + "/api/v1/deal/commitment-piece")
					req, err := http.NewRequest("POST", deltaApiUrl+"/api/v1/deal/commitment-piece", &buffer)
					if err != nil {
						fmt.Println(err)
						return err
					}
					req.Header.Add("Content-Type", "application/json")
					req.Header.Add("Authorization", "Bearer "+deltaApiKey)
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println(err)
						return err
					}
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						fmt.Println(err)
						return err
					}
					fmt.Println(string(body))
				}
			}
			return nil
		},
	}

	commpCommands = append(commpCommands, commpFileCmd)
	commpCommands = append(commpCommands, commpCarCmd)

	return commpCommands

}

func PrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}
