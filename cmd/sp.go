package cmd

import (
	"bytes"
	c "delta/config"
	"delta/utils"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type StorageProvider struct {
	Providers []Provider `json:"storageProviders"`
}
type Provider struct {
	ID                             string `json:"id"`
	Address                        string `json:"address"`
	AddressOfOwner                 string `json:"address_of_owner"`
	AddressOfWorker                string `json:"address_of_worker"`
	AddressOfBeneficiary           string `json:"address_of_beneficiary"`
	SectorSizeBytes                string `json:"sector_size_bytes"`
	MaxPieceSizeBytes              string `json:"max_piece_size_bytes"`
	MinPieceSizeBytes              string `json:"min_piece_size_bytes"`
	PriceAttofil                   string `json:"price_attofil"`
	PriceVerifiedAttofil           string `json:"price_verified_attofil"`
	BalanceAttofil                 string `json:"balance_attofil"`
	LockedFundsAttofil             string `json:"locked_funds_attofil"`
	InitialPledgeAttofil           string `json:"initial_pledge_attofil"`
	RawPowerBytes                  string `json:"raw_power_bytes"`
	QualityAdjustedPowerBytes      string `json:"quality_adjusted_power_bytes"`
	TotalRawPowerBytes             string `json:"total_raw_power_bytes"`
	TotalQualityAdjustedPowerBytes string `json:"total_quality_adjusted_power_bytes"`
	TotalStorageDealCount          string `json:"total_storage_deal_count"`
	TotalSectorsSealedByPostCount  string `json:"total_sectors_sealed_by_post_count"`
	PeerID                         string `json:"peer_id"`
	Height                         string `json:"height"`
	LotusVersion                   string `json:"lotus_version"`
	Multiaddrs                     struct {
		Addresses []string `json:"addresses"`
	} `json:"multiaddrs"`
	Metadata             interface{} `json:"metadata"`
	AddressOfControllers struct {
		Addresses []string `json:"addresses"`
	} `json:"address_of_controllers"`
	Tipset struct {
		Cids []struct {
			NAMING_FAILED string `json:"/"`
		} `json:"cids"`
	} `json:"tipset"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pricing struct {
	StoragePrice string `json:"storagePrice"`
}

func SpCmd(cfg *c.DeltaConfig) []*cli.Command {
	var spCommands []*cli.Command
	spCmd := &cli.Command{
		Name:  "sp",
		Usage: "SP CLI using data.storage.market API",
		Subcommands: []*cli.Command{
			{
				Name:  "info",
				Usage: "Get storage provider info",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Usage: "The address of the storage provider",
					},
				}, Action: func(context *cli.Context) error {

				addr := context.String("addr")
				if addr == "" {
					return fmt.Errorf("addr is required")
				}

				provider, err := fetchProviderByAddr(addr)
				if err != nil {
					fmt.Println(err)
					return err
				}
				var buffer bytes.Buffer
				err = utils.PrettyEncode(provider, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())
				return nil
			},
			},
			{
				Name:  "selection",
				Usage: "Get a random storage provider",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:  "min-piece-size",
						Usage: "The minimum piece size to consider in bytes",
					},
					&cli.Int64Flag{
						Name:  "max-piece-size",
						Usage: "The maximum piece size to consider in bytes",
					},
				}, Action: func(context *cli.Context) error {
				// Parse query parameters
				min := context.Int64("min-piece-size")
				max := context.Int64("max-piece-size")

				if min == 0 {
					min = 1
				}
				if max == 0 {
					max = 1 << 62
				}

				providers, err := fetchProviders(min, max)
				if err != nil {
					fmt.Println(err)
					return err
				}

				// Select one random provider
				rand.Seed(time.Now().UnixNano())
				randomIndex := rand.Intn(len(providers))
				randomProvider := providers[randomIndex]

				if err != nil {
					panic(err)
				}
				var buffer bytes.Buffer
				err = utils.PrettyEncode(randomProvider, &buffer)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(buffer.String())

				return nil
			},
			},
		},
	}
	spCommands = append(spCommands, spCmd)

	return spCommands
}

func fetchProviderByAddr(addr string) (Provider, error) {
	response, err := http.Get("https://data.storage.market/api/providers/" + addr)
	if err != nil {
		return Provider{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Provider{}, err
	}

	var providers Provider
	err = json.Unmarshal(body, &providers)

	if err != nil {
		return Provider{}, err
	}

	return providers, nil
}

func fetchProviders(minPieceSize, maxPieceSize int64) ([]Provider, error) {
	response, err := http.Get("https://data.storage.market/api/providers")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var providers StorageProvider
	err = json.Unmarshal(body, &providers)
	if err != nil {
		return nil, err
	}

	// Filter providers by piece size
	var filteredProviders []Provider
	for _, provider := range providers.Providers {

		pvMinPieceSize, _ := strconv.Atoi(provider.MinPieceSizeBytes)
		pvMaxPieceSize, _ := strconv.Atoi(provider.MaxPieceSizeBytes)
		if int64(pvMinPieceSize) >= minPieceSize && int64(pvMaxPieceSize) <= maxPieceSize {
			filteredProviders = append(filteredProviders, provider)
		}
	}

	return filteredProviders, nil
}
