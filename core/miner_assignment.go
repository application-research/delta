package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

type MinerAssignmentService struct {
	DeltaNode DeltaNode
}

// GetSPInfo
func NewMinerAssignmentService(node DeltaNode) *MinerAssignmentService {
	return &MinerAssignmentService{
		DeltaNode: node,
	}
}

// A function that takes in a parameter, byteSize, and returns a Provider and an error.
func (m MinerAssignmentService) GetSPWithGivenBytes(byteSize int64) (Provider, error) {
	bytSizeStr := fmt.Sprintf("%d", byteSize)
	fmt.Println("Getting SP with given bytes: ", m.DeltaNode.Config.ExternalApis.SpSelectionApi+"?size_bytes="+bytSizeStr)
	resp, err := http.Get(m.DeltaNode.Config.ExternalApis.SpSelectionApi + "?size_bytes=" + bytSizeStr)
	if err != nil {
		// handle error
		fmt.Println("Error making HTTP request:", err)
		return Provider{}, err
	}
	defer resp.Body.Close()
	var provider Provider
	err = json.NewDecoder(resp.Body).Decode(&provider)
	if err != nil {
		// handle error
		fmt.Println("Error decoding JSON:", err)
		return Provider{}, err
	}

	return provider, nil

}

// A function that takes in two parameters, byteSize and sourceIp, and returns a Provider and an error.
func (m MinerAssignmentService) GetSPWithGivenBytesAndIp(byteSize string, sourceIp string) (Provider, error) {
	resp, err := http.Get(m.DeltaNode.Config.ExternalApis.SpSelectionApi + "?size_bytes=" + byteSize + "&source_ip=" + sourceIp)
	if err != nil {
		// handle error
		fmt.Println("Error making HTTP request:", err)
		return Provider{}, err
	}
	defer resp.Body.Close()
	var provider Provider
	err = json.NewDecoder(resp.Body).Decode(&provider)
	if err != nil {
		// handle error
		fmt.Println("Error decoding JSON:", err)
		return Provider{}, err
	}

	return provider, nil
}
