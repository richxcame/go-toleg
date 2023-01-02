package check

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

// CheckDestination checks the given destination
//
// POST /api/<username>/<server>/cd/add
// destination - msisdn (no country code for mts and tmcell)
// service - service key
// ts - epoch time at client
// hmac - hmac with shared key
//
// msg = <destination>:<service>:<ts>:<username>
func CheckDestination(phone, service string) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s", phone, service, ts, constants.USERNAME)
	data := url.Values{
		"destination": {phone},
		"service":     {service},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.CHECK_DESTINATION_URL, data)
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
