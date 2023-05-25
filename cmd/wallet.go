package cmd

import (
	"bytes"
	c "delta/config"
	"delta/core"
	"delta/utils"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/urfave/cli/v2"
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

type WalletResponse struct {
	PublicKey  string `json:"public_key,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	KeyType    string `json:"key_type,omitempty"`
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
				Name:  "generate",
				Usage: "Generate a new wallet",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "dir",
						Usage: "Wallet directory where the wallet file will be stored",
					},
					&cli.BoolFlag{
						Name:  "show-private-key",
						Usage: "Show private key",
						Value: false,
					},
				},
				Action: func(context *cli.Context) error {
					walletDir := context.String("dir")
					showPrivateKey := context.Bool("show-private-key")

					wallet, err := core.SetupWallet(walletDir)
					if err != nil {
						return err
					}
					walletAddr, err := wallet.GetDefault()

					walletResponse := WalletResponse{
						PublicKey: walletAddr.String(),
					}

					if showPrivateKey {
						// import the new wallet
						hexedKey := hex.EncodeToString(walletAddr.Payload())
						decodeKey := base64.StdEncoding.EncodeToString([]byte(hexedKey))
						walletResponse.PrivateKey = decodeKey
					}

					if err != nil {
						panic(err)
					}
					var buffer bytes.Buffer
					err = utils.PrettyEncode(walletResponse, &buffer)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(buffer.String())
					fmt.Println(utils.Purple + "Wallet generated successfully. Make sure to backup your wallet file." + utils.Purple)

					return nil
				},
			},
			{
				Name:  "register",
				Usage: "Register a new wallet",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "hex",
						Usage: "Hexed wallet from LOTUS/BOOSTD export",
					},
				},
				Action: func(context *cli.Context) error {
					cmd, err := NewDeltaCmdNode(context)
					if err != nil {
						return err
					}

					hexParam := context.String("hex")
					url := cmd.DeltaApi + "/admin/wallet/register-hex"
					payload := map[string]string{
						"hex_key": hexParam,
					}
					data, err := json.Marshal(payload)
					if err != nil {
						panic(err)
					}

					req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
					req.Header.Set("Authorization", "Bearer "+cmd.DeltaAuth)
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
				Action: func(context *cli.Context) error {
					cmd, err := NewDeltaCmdNode(context)
					if err != nil {
						return err
					}

					url := cmd.DeltaApi + "/admin/wallet/list"
					req, err := http.NewRequest("GET", url, nil)
					if err != nil {
						panic(err)
					}
					req.Header.Set("Authorization", "Bearer "+cmd.DeltaAuth)
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
