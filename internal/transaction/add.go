package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	pb "github.com/richxcame/gotoleg/gotoleg"
	"github.com/richxcame/gotoleg/internal/constants"
	"github.com/richxcame/gotoleg/internal/utility"
	"github.com/richxcame/gotoleg/pkg/hmacsha1"
)

// TODO: make standart struct with types
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

type Server struct {
	pb.UnimplementedTransactionServer
}

// Add sends money to given destination
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
func (s *Server) Add(ctx context.Context, in *pb.AddTransactionRequest) (*pb.AddTransactionReply, error) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		fmt.Println(err, "epoch")
		return nil, err
	}
	fmt.Println(epochTime)

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", in.LocalID, in.Service, in.Amount, in.Phone, ts, ts, constants.USERNAME)
	data := url.Values{
		"local-id":    {in.LocalID},
		"service":     {in.Service},
		"amount":      {in.Amount},
		"destination": {in.Amount},
		"txn-ts":      {ts},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.ADD_TRANSACTION_URL, data)
	if err != nil {
		fmt.Println(err, "post")
		return nil, errors.New("")
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err, "respInBytes")
		return nil, errors.New("")
	}

	var result AddTransactionResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		fmt.Println(err, "marshall")
		return nil, errors.New("")
	}

	if result.Status != "SUCCESS" {
		// TODO: add more types
		fmt.Println(err, "is not success", result)
		return nil, errors.New("")
	}
	return &pb.AddTransactionReply{Status: "SUCCESS"}, nil
}
