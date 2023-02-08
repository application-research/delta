package jobs

import (
	"encoding/json"
	"fc-deal-making-service/core"
	"fmt"
	"gorm.io/gorm"
	"net/http"
)

type MinerCheckProcessor struct {
	Processor
}

// job run to get updated miner list and their piece sizes
func NewMinerCheckProcessor(ln *core.LightNode) IProcessor {
	return &MinerCheckProcessor{
		Processor{
			LightNode: ln,
		},
	}
}

func (m MinerCheckProcessor) Run() error {

	// remove any record of the miner on the list
	m.LightNode.DB.Transaction(func(tx *gorm.DB) error {
		// delete all
		tx.Delete(core.MinerInfo{})
		tx.Delete(core.MinerPrice{})
		return nil
	})

	// refresh list
	// get the miner list
	var minerInfos []core.MinerInfo
	req, err := http.Get("https://api.estuary.tech/public/miners/")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer req.Body.Close()
	json.NewDecoder(req.Body).Decode(&minerInfos)
	fmt.Println(minerInfos)

	// for each miner info, get miner price
	var minerPrices []core.MinerPrice
	for _, minerInfo := range minerInfos {
		var minerPrice core.MinerPrice
		reqMinerPrice, err := http.Get("https://api.estuary.tech/public/miners/storage/query/" + minerInfo.Addr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer reqMinerPrice.Body.Close()
		json.NewDecoder(reqMinerPrice.Body).Decode(&minerPrice)
		minerPrices = append(minerPrices, minerPrice)
	}

	return nil
}
