package utils

import (
	"fmt"
	"math"
	"time"
)

func DateToHeight(time time.Time) int64 {
	// time to string date only
	dateStr := time.Format("2006-01-02")
	return UnixToHeight(DateToUnixEpoch(dateStr, "00:00:00"))
}

func HeightToDate(height int64) time.Time {
	return time.Unix(HeightToUnix(height)*1000, 0)
}

func DateToUnixEpoch(dateStr string, timeStr string) int64 {
	dtStr := fmt.Sprintf("%sT%sZ", dateStr, timeStr)
	timeT, err := time.Parse(time.RFC3339, dtStr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	unixTimeSeconds := float64(timeT.Unix())
	return int64(unixTimeSeconds)
}

func UnixToHeight(unixEpoch int64) int64 {
	return int64(math.Floor(float64(unixEpoch-FILECOIN_GENESIS_UNIX_EPOCH) / 30))
}

func HeightToUnix(height int64) int64 {
	return (height * 30) + FILECOIN_GENESIS_UNIX_EPOCH
}
