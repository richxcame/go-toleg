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

type AddTransactionResult struct {
	Status      string `json:"status"`
	RefNum      int64  `json:"ref-num"`
	Service     string `json:"service"`
	Destination string `json:"destination"`
	Amount      int    `json:"amount"`
	State       string `json:"state"`
}
type AddTransactionResp struct {
	Status       string               `json:"status,omitempty"`
	ErrorCode    int                  `json:"error-code,omitempty"`
	ErrorMessage string               `json:"error-msg,omitempty"`
	Result       AddTransactionResult `json:"result,omitempty"`
}

// Add sends money to give destination
//
// POST /api/<username>/<server>/txn/add
// local-id - client id of transaction (max 30 chars)
// service - service key (can be received from 0.3. Services Request)
// amount - transaction amount in cents (100 for 1 manat)
// destination - msisdn (no country code for mts and tmcell)
// txn-ts - epoch time of transaction
// ts - epoch time at client
// hmac - hmac with access token
//
// msg = <local-id>:<service>:<amount>:<destination>:<txn-ts>:<ts>:<username>
func Add(phone, amount string) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	localID := "2"
	service := ""
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

	resp, err := http.PostForm(constants.ADD_TRANSACTION_URL, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result AddTransactionResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result.Status != "SUCCESS" {
		log.Fatal(err)
	}
	fmt.Println(result)
}
