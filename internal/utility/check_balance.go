package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"gotoleg/internal/constants"
	"gotoleg/pkg/hmacsha1"
)

type CheckBalanceResult struct {
	AvailableFunds int64 `json:"available-funds"`
}
type CheckBalanceResp struct {
	Status       string             `json:"status,omitempty"`
	ErrorCode    int                `json:"error-code,omitempty"`
	ErrorMessage string             `json:"error-msg,omitempty"`
	Result       CheckBalanceResult `json:"result,omitempty"`
}

// CheckBalance fetches epoch time, generates hmac hash and makes post request to get your balance
// Responses balance in tenges
//
// POST /api/<username>/<server>/dealer/balance
// ts - epoch time at client
// hmac - hmac with access token
//
// msg = <ts>:<username>
func CheckBalance() (int64, error) {
	// Get epoch time
	epochTime, err := GetEpoch()
	if err != nil {
		return 0, err
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s", ts, constants.USERNAME)
	data := url.Values{
		"ts":   {ts},
		"hmac": {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.BALANCE_URL, data)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var result CheckBalanceResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		return 0, err
	}
	if result.Status != "SUCCESS" {
		return 0, errors.New(result.ErrorMessage)
	}
	return result.Result.AvailableFunds, nil
}
