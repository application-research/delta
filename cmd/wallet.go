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

type WalletRegisterResponse struct {
	Message    string `json:"message"`
	WalletAddr string `json:"wallet_addr"`
	WalletUUID string `json:"wallet_uuid"`
}

type WalletListResponse struct {
	Wallets []struct {
		ID         int       `json:"ID"`
		UUID       string    `json:"uuid"`
		Addr       string    `json:"addr"`
		Owner      string    `json:"owner"`
		KeyType    string    `json:"key_type"`
		PrivateKey string    `json:"private_key"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"wallets"`
}

// TODO: add a command to manage wallet via CLI
func WalletCmd(cfg *c.DeltaConfig) []*cli.Command {
	// add a command to run API node
	var walletCommands []*cli.Command

	walletCmd := &cli.Command{
		Name:  "wallet",
		Usage: "Run Delta wallet commands",
		Subcommands: []*cli.Command{
			{
				Name:  "register",
				Usage: "Register a new wallet",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "delta-host",
						Usage: "the delta host",
						Value: "http://localhost:1414",
					},
					&cli.StringFlag{
						Name:  "hex",
						Usage: "Hexed wallet from LOTUS/BOOSTD export",
					},
					&cli.StringFlag{
						Name:  "api-key",
						Usage: "The API key to use for the request",
					},
				},
				Action: func(context *cli.Context) error {
					deltaHostParam := context.String("delta-host")
					hexParam := context.String("hex")
					apiKeyParam := context.String("api-key")
					url := deltaHostParam + "/admin/wallet/register-hex"
					payload := map[string]string{
						"hex_key": hexParam,
					}
					data, err := json.Marshal(payload)
					if err != nil {
						panic(err)
					}

					req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
					req.Header.Set("Authorization", "Bearer "+apiKeyParam)
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					var response WalletRegisterResponse
					err = json.NewDecoder(resp.Body).Decode(&response)
					if err != nil {
						panic(err)
					}
					var buffer bytes.Buffer
					err = utils.PrettyEncode(response, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "List all wallets associated with the API key",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "delta-host",
						Usage: "the delta host",
						Value: "http://localhost:1414",
					},
					&cli.StringFlag{
						Name:  "api-key",
						Usage: "The API key to use for the request",
					},
				},
				Action: func(context *cli.Context) error {
					deltaHostParam := context.String("delta-host")
					apiKeyParam := context.String("api-key")
					url := deltaHostParam + "/admin/wallet/list"
					req, err := http.NewRequest("GET", url, nil)
					if err != nil {
						panic(err)
					}
					req.Header.Set("Authorization", "Bearer "+apiKeyParam)
					req.Header.Set("Content-Type", "application/json")

					var walletListResponse WalletListResponse
					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					err = json.NewDecoder(resp.Body).Decode(&walletListResponse)
					if err != nil {
						panic(err)
					}
					var buffer bytes.Buffer
					err = utils.PrettyEncode(walletListResponse, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
					return nil
				},
			},
		},
	}

	walletCommands = append(walletCommands, walletCmd)

	return walletCommands
}
