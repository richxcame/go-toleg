package handlers

import (
	"context"
	"fmt"
	"gotoleg/internal/db"
	"gotoleg/pkg/logger"
	"gotoleg/web/entities"
	"net/http"
	"os"
	"time"

	pb "gotoleg/rpc/gotoleg"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func SendTransactions(ctx *gin.Context) {
	rows, err := db.DB.Query(context.Background(), "SELECT * FROM transactions where status = '' AND result_status=''")
	if err != nil {
		logger.Error(err)
	}
	defer rows.Close()

	transactions := make([]entities.Transaction, 0)
	for rows.Next() {
		var trxn entities.Transaction

		if err := rows.Scan(&trxn.UUID, &trxn.CreatedAt, &trxn.UpdatedAt, &trxn.RequestLocalID, &trxn.RequestService, &trxn.RequestPhone, &trxn.RequestAmount, &trxn.Status, &trxn.ErrorCode, &trxn.ErrorMsg, &trxn.ResultStatus, &trxn.ResultRefNum, &trxn.ResultService, &trxn.ResultDestination, &trxn.ResultAmount, &trxn.ResultState, &trxn.IsChecked, &trxn.Client); err != nil {
			logger.Errorf("row scan error %v", err)
		}
		transactions = append(transactions, trxn)
	}
	fmt.Println(len(transactions))
}

// SendTransaction sends money to given phone number if the transaction didn't send
// For example: If couldn't get response from epoch time request, it might be the transaction didn't send to client
func SendTransaction(ctx *gin.Context) {
	// Get UUID from URL param
	uuid, ok := ctx.Params.Get("uuid")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "uuid is required",
			"message": "Coulnd't find UUID",
		})
		return
	}

	// Find the transaction with given UUID
	var trxn entities.Transaction
	err := db.DB.QueryRow(context.Background(), "SELECT * FROM transactions where uuid = $1", uuid).Scan(&trxn.UUID, &trxn.CreatedAt, &trxn.UpdatedAt, &trxn.RequestLocalID, &trxn.RequestService, &trxn.RequestPhone, &trxn.RequestAmount, &trxn.Status, &trxn.ErrorCode, &trxn.ErrorMsg, &trxn.ResultStatus, &trxn.ResultRefNum, &trxn.ResultService, &trxn.ResultDestination, &trxn.ResultAmount, &trxn.ResultState, &trxn.IsChecked, &trxn.Client)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Couldn't find the transaction",
		})
		return
	}

	// Check if transaction is already sent
	if trxn.Status != "" || trxn.ResultStatus != "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "transaction already sent",
			"message": "The transaction had been already sent",
		})
		return
	}

	// Send money
	result, hasError := send(trxn.RequestLocalID, trxn.RequestService, trxn.RequestPhone, trxn.RequestAmount)
	if hasError {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "grpc add error",
			"message": "Couldn't send money to given phone",
		})
		return
	}

	// Update after send money
	const sqlStr = `UPDATE transactions SET status=$1, error_code=$2, error_msg=$3, result_status=$4, result_ref_num=$5, result_service=$6, result_destination=$7, result_amount=$8, result_state=$9, updated_at=$10 WHERE uuid=$11`
	_, err = db.DB.Exec(context.Background(), sqlStr, result.Status, result.ErrorCode, result.ErrorMessage, result.Result.Status, result.Result.RefNum, result.Result.Service, result.Result.Destination, result.Result.Amount, result.Result.State, time.Now(), trxn.RequestLocalID)
	if err != nil {
		logger.Errorf("couldn't update database: %v, result: %v", err, result)
	}

	// Send success result
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaction has been sent",
	})
}

func send(id, service, phone, amount string) (result *pb.TransactionReply, hasError bool) {
	addr := os.Getenv("GOTOLEG_PORT")
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", addr), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("grpc connection error: %v", err)
		return nil, true
	}
	defer conn.Close()

	c := pb.NewTransactionClient(conn)
	// Contact the server and print out its response.
	context := metadata.AppendToOutgoingContext(context.Background(), "api_key", os.Getenv("GOTOLEG_SUPER_KEY"))
	result, err = c.Add(context, &pb.TransactionRequest{LocalID: id, Service: service, Phone: phone, Amount: amount})
	if err != nil {
		logger.Errorf("grpc add() error: %v", err)
		return nil, true
	}
	return result, false
}
