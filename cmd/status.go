package cmd

import (
	"bytes"
	c "delta/config"
	"delta/utils"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"net/http"
	"time"
)

type StatusResponse struct {
	Content struct {
		ID                int       `json:"ID"`
		Name              string    `json:"name"`
		Size              int       `json:"size"`
		Cid               string    `json:"cid"`
		PieceCommitmentID int       `json:"piece_commitment_id"`
		Status            string    `json:"status"`
		RequestType       string    `json:"request_type"`
		ConnectionMode    string    `json:"connection_mode"`
		LastMessage       string    `json:"last_message"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	} `json:"content"`
	DealProposalParameters []struct {
		ID        int       `json:"ID"`
		Content   int       `json:"content"`
		Label     string    `json:"label"`
		Duration  int       `json:"duration"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"deal_proposal_parameters"`
	DealProposals    interface{} `json:"deal_proposals"`
	Deals            interface{} `json:"deals"`
	PieceCommitments []struct {
		ID                 int       `json:"ID"`
		Cid                string    `json:"cid"`
		Piece              string    `json:"piece"`
		Size               int       `json:"size"`
		PaddedPieceSize    int       `json:"padded_piece_size"`
		UnnpaddedPieceSize int       `json:"unnpadded_piece_size"`
		Status             string    `json:"status"`
		LastMessage        string    `json:"last_message"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
	} `json:"piece_commitments"`
}

func StatusCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var statusCommands []*cli.Command

	statusCmd := &cli.Command{
		Name:        "status",
		Usage:       "Get the status of a content, deal, or piece commitment",
		Description: "Get the status of a content, deal, or piece commitment. The type of status can be either content, deal, or piece-commitment.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "delta-host",
				Usage: "the delta host",
				Value: "http://localhost:1414",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "content, deal, or piece-commitment",
			},
			&cli.StringFlag{
				Name:  "id",
				Usage: "the id of the content, deal, or piece-commitment",
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "The API key to use for the request",
			},
		},
		Action: func(context *cli.Context) error {
			deltaHostParam := context.String("delta-host")
			typeParam := context.String("type")
			apiKeyParam := context.String("api-key")
			idParam := context.String("id")

			fmt.Println(deltaHostParam)
			fmt.Println(typeParam)
			fmt.Println(apiKeyParam)
			fmt.Println(idParam)
			var dealStatusResponse StatusResponse
			url := deltaHostParam + "/open/stats/" + typeParam + "/" + idParam
			if typeParam == "content" {
				// Create a new HTTP request with the desired method and URL.
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					panic(err)
				}

				// Set the Authorization header.
				req.Header.Set("Authorization", "Bearer "+apiKeyParam)

				// Send the HTTP request and print the response.
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				// Print the response status code.
				fmt.Println(resp.Status)
				err = json.NewDecoder(resp.Body).Decode(&dealStatusResponse)
				if err != nil {
					panic(err)
				}
				var buffer bytes.Buffer
				err = utils.PrettyEncode(dealStatusResponse, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())
			}

			if typeParam == "deal" {
				// TODO: implement
				fmt.Println("Not implemented yet")
			}

			if typeParam == "piece-commitment" {
				// TODO: implement
				fmt.Println("Not implemented yet")
			}

			return nil
		},
	}
	statusCommands = append(statusCommands, statusCmd)

	return statusCommands
}
