package utility

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"gotoleg/internal/constants"
	"gotoleg/pkg/hmacsha1"
)

type ServicesResult struct {
	Services []interface{} `json:"services"`
}
type ServicesResp struct {
	Status       string         `json:"status,omitempty"`
	ErrorCode    int            `json:"error-code,omitempty"`
	ErrorMessage string         `json:"error-msg,omitempty"`
	Result       ServicesResult `json:"result,omitempty"`
}

// GetServices fetches epoch time, generates hmac hash and makes post request to get services list
// Responses list of enabled service keys
//
// POST /api/<username>/<server>/dealer/services
// ts - epoch time at client
// hmac - hmac with access token
//
// msg = <ts>:<username>
func GetServices() ([]interface{}, error) {
	// Get epoch time
	epochTime, err := GetEpoch()
	if err != nil {
		return nil, err
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s", ts, constants.USERNAME)
	data := url.Values{
		"ts":   {ts},
		"hmac": {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.SERVICES_URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ServicesResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		return nil, err
	}
	return result.Result.Services, nil
}
