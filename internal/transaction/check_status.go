package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"gotoleg/internal/constants"
	"gotoleg/internal/utility"
	"gotoleg/pkg/hmacsha1"
)

// CheckStatus checks status of transaction with given localID
//
// POST /api/<username>/<server>/txn/info
// local-id - client id of transaction (max 30 chars)
// ts - epoch time at client
// hmac - hmac with shared key
//
// msg = <local-id>:<ts>:<username>
func CheckStatus(localID string) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result GarynjaResponse
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Status != "SUCCESS" {
		log.Fatal(err)
	}
	fmt.Println(result)
}
