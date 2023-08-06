package utils

import (
	"fmt"
	"delta/models"
	"github.com/application-research/filclient"
	datatransfer "github.com/filecoin-project/go-data-transfer"
)

// GetChannelID `GetChannelID` takes a string and returns a `datatransfer.ChannelID` and an error
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
