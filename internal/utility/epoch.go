package utility

import (
	"encoding/json"
	"io"
	"net/http"

	"gotoleg/internal/constants"
)

// GetEpoch fetches unix time in seconds from "garynja"
// Returns epoch time
//
// GET /api/epoch
func GetEpoch() (int64, error) {
	resp, err := http.Get(constants.EPOCH_URL)
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
