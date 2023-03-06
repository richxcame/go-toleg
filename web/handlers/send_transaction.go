package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"gotoleg/internal/constants"
	"gotoleg/internal/db"
	"gotoleg/internal/transaction"
	"gotoleg/internal/utility"
	"gotoleg/pkg/hmacsha1"
	"gotoleg/pkg/logger"
	"gotoleg/web/entities"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	sucessCount := 0
	errorCount := 0
	for _, v := range transactions {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = []gin.Param{{Key: "uuid", Value: v.UUID.String()}}
		SendTransaction(ctx)
		if w.Result().StatusCode == 200 {
			sucessCount++
		} else {
			errorCount++
		}
	}
	ctx.JSON(200, gin.H{
		"sucess_count": sucessCount,
		"error_count":  errorCount,
	})
}

// SendTransaction resends money if status or result_status is equal to empty string
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
	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		logger.Errorf("epoch time get error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't get epoch time",
		})
		return
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", trxn.RequestLocalID, trxn.RequestService, trxn.RequestAmount, trxn.RequestPhone, ts, ts, constants.USERNAME)
	data := url.Values{
		"local-id":    {trxn.RequestLocalID},
		"service":     {trxn.RequestService},
		"amount":      {trxn.RequestAmount},
		"destination": {trxn.RequestPhone},
		"txn-ts":      {ts},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	// Send post request to send money
	resp, err := http.PostForm(constants.ADD_TRANSACTION_URL, data)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't post transaction's data",
		})
		return
	}
	defer resp.Body.Close()

	// Read response bytes
	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't parse response bytes",
		})
		return
	}

	// Parse response bytes
	var result transaction.GarynjaResponse
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't parse from bytes to struct",
		})
		return
	}

	// Update after send money
	const sqlStr = `UPDATE transactions SET status=$1, error_code=$2, error_msg=$3, result_status=$4, result_ref_num=$5, result_service=$6, result_destination=$7, result_amount=$8, result_state=$9, updated_at=$10 WHERE uuid=$11`
	_, err = db.DB.Exec(context.Background(), sqlStr, result.Status, result.ErrorCode, result.ErrorMessage, result.Result.Status, result.Result.RefNum, result.Result.Service, result.Result.Destination, result.Result.Amount, result.Result.State, time.Now(), trxn.UUID)
	if err != nil {
		logger.Errorf("couldn't update database: %v, result: %v", err, result)
	}

	// Send success result
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaction has been sent",
	})
}
