package handlers

import (
	"gotoleg/internal/db"
	"gotoleg/internal/transaction"
	"gotoleg/pkg/logger"
	"gotoleg/web/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckTrxnStatus(c *gin.Context) {
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
	err := db.DB.QueryRow(ctx, "SELECT uuid, created_at, updated_at, request_local_id, request_service, request_phone, request_amount, status, error_code, error_msg, result_status, result_ref_num, result_service, result_destination, result_amount, result_state, result_reason, is_checked, client, note FROM transactions where uuid = $1", uuid).Scan(&trxn.UUID, &trxn.CreatedAt, &trxn.UpdatedAt, &trxn.RequestLocalID, &trxn.RequestService, &trxn.RequestPhone, &trxn.RequestAmount, &trxn.Status, &trxn.ErrorCode, &trxn.ErrorMsg, &trxn.ResultStatus, &trxn.ResultRefNum, &trxn.ResultService, &trxn.ResultDestination, &trxn.ResultAmount, &trxn.ResultState, &trxn.ResultReason, &trxn.IsChecked, &trxn.Client, &trxn.Note)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Couldn't find the transaction",
		})
		return
	}

	result, err := transaction.CheckStatus(trxn.RequestLocalID)
	if err != nil {
		logger.Errorf("couldn't check status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't check status of the transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transaction": result,
	})
}
