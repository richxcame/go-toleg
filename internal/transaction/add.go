package transaction

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gotoleg/internal/config"
	"gotoleg/internal/constants"
	"gotoleg/internal/db"
	"gotoleg/internal/utility"
	"gotoleg/pkg/arrs"
	"gotoleg/pkg/hmacsha1"
	"gotoleg/pkg/logger"
	pb "gotoleg/rpc/gotoleg"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

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
	// Get metadata and api_key from it
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Error("failed to get metadata")
		return nil, status.Errorf(codes.Unauthenticated, "failed to get metadata")
	}
	apiKey := make([]string, 0)
	apiKey = append(apiKey, md.Get("api_key")...)

	// Check the given "api_key" included in list of clients
	client, hasInList := arrs.HasMapWithKey(config.Clients, apiKey[0])
	if !hasInList {
		logger.Error("wrong api_key")
		return nil, status.Errorf(codes.InvalidArgument, "wrong api_key")
	}

	// Save the transaction to the database
	_uuid := uuid.New().String()
	sqlStatement := `
		INSERT INTO transactions (uuid, created_at, updated_at, client, request_local_id, request_service, request_phone, request_amount, note)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

	mnt, err := strconv.ParseFloat(in.Amount, 64)
	if err != nil {
		logger.Error(err, in)
		return nil, err
	}
	amount := strconv.Itoa(int(mnt))

	_, err = db.DB.Exec(context.Background(), sqlStatement, _uuid, time.Now(), time.Now(), client, in.LocalID, in.Service, in.Phone, amount, in.Note)
	if err != nil {
		logger.Error(err, in)
	}

	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		logger.Error(err, in)
		return nil, err
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", in.LocalID, in.Service, amount, in.Phone, ts, ts, constants.USERNAME)
	data := url.Values{
		"local-id":    {in.LocalID},
		"service":     {in.Service},
		"amount":      {amount},
		"destination": {in.Phone},
		"txn-ts":      {ts},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	resp, err := http.PostForm(constants.ADD_TRANSACTION_URL, data)
	if err != nil {
		logger.Error(err, in)
		return nil, err
	}
	defer resp.Body.Close()

	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err, in)
		return nil, err
	}

	var result GarynjaResponse
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		logger.Error(err, in)
		return nil, err
	}

	const sqlStr = `UPDATE transactions SET status=$1, error_code=$2, error_msg=$3, result_status=$4, result_ref_num=$5, result_service=$6, result_destination=$7, result_amount=$8, result_state=$9, updated_at=$10 WHERE uuid=$11`
	_, err = db.DB.Exec(context.Background(), sqlStr, result.Status, result.ErrorCode, result.ErrorMessage, result.Result.Status, result.Result.RefNum, result.Result.Service, result.Result.Destination, result.Result.Amount, result.Result.State, time.Now(), _uuid)
	if err != nil {
		logger.Error(err, in)
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
