package jobs

import (
	"delta/core"
	"encoding/json"
	"fmt"
	model "delta/models"
	"gorm.io/gorm"
	"net/http"
)

type MinerCheckProcessor struct {
	Processor
}

// NewMinerCheckProcessor job run to get updated miner list and their piece sizes
// `NewMinerCheckProcessor` creates a new `MinerCheckProcessor` object and returns it
func NewMinerCheckProcessor(ln *core.DeltaNode) IProcessor {
	return &MinerCheckProcessor{
		Processor{
			LightNode: ln,
		},
	}
}

// Run Getting the list of miners and their prices from the API and storing them in the database.
func (m MinerCheckProcessor) Run() error {

	// remove any record of the miner on the list
	m.LightNode.DB.Transaction(func(tx *gorm.DB) error {
		// delete all
		tx.Delete(model.MinerInfo{}).Where("id > 0")
		return nil
	})

	// refresh list
	// get the miner list
	var minerInfos []model.MinerInfo
	req, err := http.Get("https://api.estuary.tech/public/miners/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer req.Body.Close()
	json.NewDecoder(req.Body).Decode(&minerInfos)
	fmt.Println(minerInfos)

	// for each miner info, get miner price
	var minerPrices []model.MinerPrice
	for _, minerInfo := range minerInfos {
		fmt.Println(minerInfo.Addr)
		var minerPrice model.MinerPrice
		reqMinerPrice, err := http.Get("https://api.estuary.tech/public/miners/storage/query/" + minerInfo.Addr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer reqMinerPrice.Body.Close()
		json.NewDecoder(reqMinerPrice.Body).Decode(&minerPrice)
		minerPrices = append(minerPrices, minerPrice)
	}

	fmt.Println(minerPrices)

	m.LightNode.DB.Transaction(func(tx *gorm.DB) error {
		// insert all
		tx.Create(&minerInfos)
		tx.Create(&minerPrices)
		return nil
	})
	return nil
}
