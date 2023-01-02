package transaction

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/richxcame/gotoleg/internal/constants"
	"github.com/richxcame/gotoleg/internal/utility"
	"github.com/richxcame/gotoleg/pkg/hmacsha1"
)

// CorrectDeclined resends declined transactions by correcting service key
//
// POST /api/<username>/<server>/txn/change-service
// local-id - client id of transaction (max 30 chars)
// service - new service key
// ts - epoch time at client
// hmac - hmac with shared key
//
// msg = <local-id>:<service>:<ts>:<username>
func CorrectDeclined(localID, service string) {

	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s", localID, service, ts, constants.USERNAME)
	data := url.Values{
		"local-id": {localID},
		"service":  {service},
		"ts":       {ts},
		"hmac":     {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.CORRECT_DECLINED_URL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result TransactionResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Status != "SUCCESS" {
		log.Fatal(err)
	}
	fmt.Println(result)
}
