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

// ForceAdd add transaction or resend declined transaction, all in one place. Note that, if transaction is DECLINED, parameters will be used from information already at gateway.
//
// POST /api/<username>/<server>/txn/force/add
// local-id - client id of transaction (max 30 chars)
// service - service key (can be received from 0.3. Services Request)
// amount - transaction amount in cents (100 for 1 manat)
// destination - msisdn (no country code for mts and tmcell)
// txn-ts - epoch time of transaction
// ts - epoch time at client
// hmac - hmac with access token
//
// msg = <local-id>:<service>:<amount>:<destination>:<txn-ts>:<ts>:<username>
func ForceAdd(amount, phone, localID, service string) (*GarynjaResponse, error) {

	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		return nil, err
	}

	// Prepare request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", localID, service, amount, phone, ts, ts, constants.USERNAME)
	data := url.Values{
		"local-id":    {localID},
		"service":     {service},
		"amount":      {amount},
		"destination": {phone},
		"txn-ts":      {ts},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.FORCE_ADD_URL, data)
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
