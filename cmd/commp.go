package cmd

import (
	"bufio"
	"bytes"
	"delta/api"
	"delta/core"
	"delta/utils"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"os"
)

type CommpResult struct {
	Cid               string `json:"cid,omitempty"`
	Commp             string `json:"piece_commitment,omitempty"`
	PaddedPieceSize   uint64 `json:"padded_piece_size,omitempty"`
	UnPaddedPieceSize uint64 `json:"unpadded_piece_size,omitempty"`
}

// CommpCmd A CLI command that generates a piece commitment for a given file.
// `CommpCmd()` returns a slice of `*cli.Command`s
func CommpCmd() []*cli.Command {
	// add a command to run API node
	var commpCommands []*cli.Command
	var commpService = new(core.CommpService)
	commpFileCmd := &cli.Command{
		Name: "commp",
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
			var commpResult api.DealRequest
			file := c.String("file")
			forImport := c.Bool("for-import")

			miner := c.String("miner")
			openFile, err := os.Open(file)
			reader := bufio.NewReader(openFile)
			if err != nil {
				fmt.Println(err)
				return err
			}

			if miner != "" {
				commpResult.Miner = miner
			}

			if file != "" {
				pieceInfo, err := commpService.GenerateParallelCommp(reader)

				if err != nil {
					fmt.Println(err)
					return err
				}
				// return json to console.
				commpResult.PieceCommitment.Piece = pieceInfo.PieceCID.String()
				commpResult.PieceCommitment.PaddedPieceSize = uint64(pieceInfo.PieceSize)
				commpResult.PieceCommitment.UnPaddedPieceSize = uint64(pieceInfo.PieceSize.Unpadded())

				// if for offline, add connection mode offline
				if forImport {
					commpResult.ConnectionMode = "import"
				} else {
					commpResult.ConnectionMode = "e2e"
				}

				if err != nil {
					fmt.Println(err)
					return err
				}
				commpResult.Size = pieceInfo.PayloadSize

				var buffer bytes.Buffer
				err = utils.PrettyEncode(commpResult, &buffer)
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
					fmt.Println(deltaApiUrl + "/api/v1/deal/piece-commitment")
					req, err := http.NewRequest("POST", deltaApiUrl+"/api/v1/deal/piece-commitment", &buffer)
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

	return commpCommands

}
