package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"gotoleg/internal/constants"
	"gotoleg/internal/utility"
	"gotoleg/pkg/hmacsha1"
)

// ResendDeclined resend declined transactions with given localID
//
// POST /api/<username>/<server>/txn/retry
// local-id - client id of transaction (max 30 chars)
// ts - epoch time at client
// hmac - hmac with shared key
func ResendDeclined(localID string) (*GarynjaResponse, error) {

	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		return nil, err
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s", localID, ts, constants.USERNAME)
	data := url.Values{
		"local-id": {localID},
		"ts":       {ts},
		"hmac":     {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.TRANSACTION_STATUS_URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result GarynjaResponse
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil

}
