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

// Callback is an address on dealer's server which will be used to inform sender of transaction about it's state change in order to make system asynchronous. Dealer no more needs to POLL the processing gateway for transaction state.
//
// POST <callback>
// username - dealer's user name
// server-label - server label
// action - 'txn-state'
// local-id - client id of transaction
// state - current state of transaction
// ts - epoch time at server
// hmac - hmac with shared key
//
// msg = <username>:<server-label>:<action>:<local-id>:<state>:<ts>
func Callback(localID, state string) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s", constants.USERNAME, constants.SERVER, "txn-state", localID, state, ts)
	data := url.Values{
		"username":     {constants.USERNAME},
		"server-label": {constants.SERVER},
		"action":       {"txn-state"},
		"local-id":     {localID},
		"state":        {state},
		"ts":           {ts},
		"hmac":         {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	// TODO: ask for callback url
	resp, err := http.PostForm(constants.CALLBACK_URL, data)
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
