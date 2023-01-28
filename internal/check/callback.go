package check

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

// Check callback
//
// POST <callback>
// username - dealer's user name
// server-label - server label
// action - 'cd-state'
// service - service_key
// destination - destination
// state - current state of transaction
// ts - epoch time at server
// hmac - hmac with shared key
//
// msg = <username>:<server-label>:<action>:<service>:<destinaiton>:<state>:<ts>
func Callback(phone, service, state string) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", constants.USERNAME, constants.SERVER, "cd-state", service, phone, state, ts)
	data := url.Values{
		"username":     {constants.USERNAME},
		"server-label": {constants.SERVER},
		"action":       {"cd-state"},
		"service":      {service},
		"destination":  {phone},
		"state":        {state},
		"ts":           {ts},
		"hmac":         {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.CHECK_CALLBACK_URL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result CheckDestinationResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Status != "SUCCESS" {
		log.Fatal(err)
	}
	fmt.Println(result)

}
