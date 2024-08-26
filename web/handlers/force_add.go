package handlers

import (
	"gotoleg/internal/db"
	"gotoleg/internal/transaction"
	"gotoleg/pkg/logger"
	"gotoleg/web/entities"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
)

// ForceAddDeclinedTransactions resends all transactions with result_status='DECLINED'
func ForceAddDeclinedTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	rows, err := db.DB.Query(ctx, "SELECT * FROM transactions where result_status='DECLINED'")
	if err != nil {
		logger.Error(err)
	}
	defer rows.Close()

	transactions := make([]entities.Transaction, 0)
	for rows.Next() {
		var trxn entities.Transaction

		if err := rows.Scan(&trxn.UUID, &trxn.CreatedAt, &trxn.UpdatedAt, &trxn.RequestLocalID, &trxn.RequestService, &trxn.RequestPhone, &trxn.RequestAmount, &trxn.Status, &trxn.ErrorCode, &trxn.ErrorMsg, &trxn.ResultStatus, &trxn.ResultRefNum, &trxn.ResultService, &trxn.ResultDestination, &trxn.ResultAmount, &trxn.ResultState, &trxn.ResultReason, &trxn.IsChecked, &trxn.Client); err != nil {
			logger.Errorf("row scan error %v", err)
		}
		transactions = append(transactions, trxn)
	}

	sucessCount := 0
	errorCount := 0
	for _, v := range transactions {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: v.UUID.String()}}
		ForceAddDeclinedTransaction(c)
		if w.Result().StatusCode == 200 {
			sucessCount++
		} else {
			errorCount++
		}
	}
	c.JSON(200, gin.H{
		"sucess_count": sucessCount,
		"error_count":  errorCount,
	})
}

// ForceAddTransaction sends or resends money to client
// Add transaction or resend declined transaction, all in one place. Note that, if transaction is DECLINED, parameters will be used from information already at gateway.
func ForceAddDeclinedTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	// Get UUID from URL param
	uuid, ok := c.Params.Get("uuid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "uuid is required",
			"message": "Coulnd't find UUID",
		})
		return
	}

	// Find the transaction with given UUID
	var trxn entities.Transaction
	err := db.DB.QueryRow(ctx, "SELECT uuid, created_at, updated_at, request_local_id, request_service, request_phone, request_amount, status, error_code, error_msg, result_status, result_ref_num, result_service, result_destination, result_amount, result_state, result_reason, is_checked, client FROM transactions where uuid = $1", uuid).
		Scan(&trxn.UUID,
			&trxn.CreatedAt,
			&trxn.UpdatedAt,
			&trxn.RequestLocalID,
			&trxn.RequestService,
			&trxn.RequestPhone,
			&trxn.RequestAmount,
			&trxn.Status,
			&trxn.ErrorCode,
			&trxn.ErrorMsg,
			&trxn.ResultStatus,
			&trxn.ResultRefNum,
			&trxn.ResultService,
			&trxn.ResultDestination,
			&trxn.ResultAmount,
			&trxn.ResultState,
			&trxn.ResultReason,
			&trxn.IsChecked,
			&trxn.Client)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Couldn't find the transaction",
		})
		return
	}

	// Check transaction is declined or not
	if trxn.ResultStatus != "DECLINED" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "transaction is not declined",
			"message": "The transaction is not decline",
		})
		return
	}

	// Force add transaction
	result, err := transaction.ForceAdd(trxn.RequestAmount, trxn.RequestPhone, trxn.RequestLocalID, trxn.RequestService)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't send or resend money",
		})
	}

	// Update after force add
	const sqlStr = `UPDATE transactions SET status=$1, error_code=$2, error_msg=$3, result_status=$4, result_ref_num=$5, result_service=$6, result_destination=$7, result_amount=$8, result_state=$9, updated_at=$10 WHERE uuid=$11`
	_, err = db.DB.Exec(ctx, sqlStr, result.Status, result.ErrorCode, result.ErrorMessage, result.Result.Status, result.Result.RefNum, result.Result.Service, result.Result.Destination, result.Result.Amount, result.Result.State, time.Now(), trxn.UUID)
	if err != nil {
		logger.Errorf("couldn't update database: %v, result: %v", err, result)
	}

	// Send success result
	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "The declined transaction is resend",
		"transaction": result,
	})
}
