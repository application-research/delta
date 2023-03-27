package cmd

import (
	"bytes"
	"delta/api"
	c "delta/config"
	"delta/utils"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type DealMetadata api.DealRequest
type DealResponse api.DealResponse

func DealCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var dealCommands []*cli.Command

	dealCmd := &cli.Command{
		Name:        "deal",
		Usage:       "Make a Network Storage Deal on Delta",
		Description: "Make a delta storage deal. The type of deal can be either e2e (online) or import (offline).",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "delta-host",
				Usage: "the delta host",
				Value: "http://localhost:1414",
			},
			&cli.StringFlag{
				Name:  "type",
				Usage: "e2e (online) or import (offline)",
			},
			&cli.StringFlag{
				Name:  "file",
				Usage: "file to make a deal with. Required only for e2e deals.",
			},
			&cli.StringFlag{
				Name:  "metadata",
				Usage: "metadata to store",
			},
			&cli.StringFlag{
				Name:  "api-key",
				Usage: "The API key to use for the request",
			},
		},
		Action: func(context *cli.Context) error {
			deltaHostParam := context.String("delta-host")
			fileParam := context.String("file")
			typeParam := context.String("type")
			apiKeyParam := context.String("api-key")
			metadataParam := context.String("metadata")
			var metadata DealMetadata
			var metadataArr []DealMetadata
			var e2eResponse DealResponse
			var importResponse []DealResponse

			// validate
			if typeParam == "e2e" {
				if fileParam == "" {
					fmt.Println("file is required for e2e deals")
					os.Exit(1)
				}
			}
			if apiKeyParam == "" {
				fmt.Println("api-key is required")
				os.Exit(1)
			}
			if metadataParam == "" {
				fmt.Println("metadata is required")
				os.Exit(1)
			}

			if typeParam == "import" {
				if fileParam != "" {
					fmt.Println("file is not required for import deals")
					os.Exit(1)
				}
			}

			var endpoint string
			url := deltaHostParam + "/api/v1"
			if typeParam == "e2e" {

				err := json.Unmarshal([]byte(metadataParam), &metadata)
				if err != nil {
					var buffer bytes.Buffer
					err = utils.PrettyEncode(DealResponse{
						Status:  "error",
						Message: "Error unmarshalling metadata",
					}, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
				}

				endpoint = url + "/deal/end-to-end"

				// Create a new HTTP request with the desired method and URL.
				req, err := http.NewRequest("POST", endpoint, nil)
				if err != nil {
					panic(err)
				}

				// Set the Authorization header.
				req.Header.Set("Authorization", "Bearer "+apiKeyParam)

				// Create a new multipart writer for the request body.
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				// Add the data file to the multipart writer.
				file, err := os.Open(fileParam)
				if err != nil {
					panic(err)
				}

				part, err := writer.CreateFormFile("data", file.Name())
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(part, file)
				if err != nil {
					panic(err)
				}

				err = writer.WriteField("metadata", metadataParam)
				if err != nil {
					panic(err)
				}

				// Close the multipart writer.
				err = writer.Close()
				if err != nil {
					panic(err)
				}

				// Set the content type header for the request.
				req.Header.Set("Content-Type", writer.FormDataContentType())

				// Set the request body to the multipart writer's buffer.
				req.Body = io.NopCloser(body)

				// Send the HTTP request and print the response.
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				fmt.Println(resp.Status)
				if resp.StatusCode != 200 {

					errorResponse := struct {
						Error struct {
							Code    int    `json:"code"`
							Reason  string `json:"reason"`
							Details string `json:"details"`
						}
					}{}

					err = json.NewDecoder(resp.Body).Decode(&errorResponse)
					if err != nil {
						fmt.Println(err)
					}

					var buffer bytes.Buffer
					err = utils.PrettyEncode(errorResponse, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
					return nil
				}
				err = json.NewDecoder(resp.Body).Decode(&e2eResponse)
				if err != nil {
					panic(err)
				}
				var buffer bytes.Buffer
				err = utils.PrettyEncode(e2eResponse, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())
				return nil
			}
			if typeParam == "import" {

				err := json.Unmarshal([]byte(metadataParam), &metadataArr)
				if err != nil {
					var buffer bytes.Buffer
					err = utils.PrettyEncode(DealResponse{
						Status:  "error",
						Message: "Error unmarshalling metadata",
					}, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
				}
				endpoint = url + "/deal/imports"

				// Create a new HTTP request with the desired method and URL.
				req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(metadataParam)))
				if err != nil {
					panic(err)
				}

				// Set the Authorization and Content-Type headers.
				req.Header.Set("Authorization", "Bearer "+apiKeyParam)
				req.Header.Set("Content-Type", "application/json")

				// Send the HTTP request and print the response.
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					errorResponse := struct {
						Error struct {
							Code    int    `json:"code"`
							Reason  string `json:"reason"`
							Details string `json:"details"`
						}
					}{}

					err = json.NewDecoder(resp.Body).Decode(&errorResponse)
					if err != nil {
						fmt.Println(err)
					}

					var buffer bytes.Buffer
					err = utils.PrettyEncode(errorResponse, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
					return nil
				}
				err = json.NewDecoder(resp.Body).Decode(&importResponse)
				if err != nil {
					panic(err)
				}

				// print output

				var buffer bytes.Buffer
				err = utils.PrettyEncode(importResponse, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())
				return nil
			}
			return nil
		},
	}
	dealCommands = append(dealCommands, dealCmd)

	return dealCommands
}
