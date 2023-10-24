package utility

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gotoleg/internal/constants"
)

// GetEpoch fetches unix time in seconds from "garynja"
// Returns epoch time
//
// GET /api/epoch
func GetEpoch() (int64, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(constants.EPOCH_URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var epochTime int64
	err = json.Unmarshal(respInBytes, &epochTime)
	if err != nil {
		return 0, err
	}
	return epochTime, nil
}
