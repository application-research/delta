package utils

import (
	"fmt"
	"github.com/application-research/delta-db/db_models"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
)

func GetChannelID(dtChannel string) (datatransfer.ChannelID, error) {
	if dtChannel == "" {
		return datatransfer.ChannelID{}, db_models.ErrNoChannelID
	}

	chid, err := filclient.ChannelIDFromString(dtChannel)
	if err != nil {
		err = fmt.Errorf("incorrectly formatted data transfer channel ID in contentDeal record: %w", err)
		return datatransfer.ChannelID{}, err
	}
	return *chid, nil
}
