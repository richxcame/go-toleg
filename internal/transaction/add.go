package transaction

import (
	"context"
	"encoding/json"
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
func (s *Server) Add(ctx context.Context, in *pb.TransactionRequest) (*pb.TransactionReply, error) {
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		return nil, err
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", in.LocalID, in.Service, in.Amount, in.Phone, ts, ts, constants.USERNAME)
	data := url.Values{
		"local-id":    {in.LocalID},
		"service":     {in.Service},
		"amount":      {in.Amount},
		"destination": {in.Phone},
		"txn-ts":      {ts},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.ADD_TRANSACTION_URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result TransactionResp
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		return nil, err
	}

	return &pb.TransactionReply{
		Status:       result.Status,
		ErrorCode:    result.ErrorCode,
		ErrorMessage: result.ErrorMessage,
		Result: &pb.Result{
			Status:      result.Result.Status,
			RefNum:      result.Result.RefNum,
			Service:     result.Result.Service,
			Destination: result.Result.Destination,
			Amount:      result.Result.Amount,
			State:       result.Result.State,
		},
	}, nil
}
