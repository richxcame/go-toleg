package handlers

import (
	"context"
	"gotoleg/internal/db"
	"gotoleg/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const GetTransactionsQuery = `select uuid, created_at, updated_at, request_local_id, request_service, request_phone, request_amount, status, error_code, error_msg, result_status, result_ref_num, result_service, result_destination, result_amount, result_state, is_checked, client from transactions offset $1 limit $2`

type Transaction struct {
	UUID              uuid.UUID `json:"uuid"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	RequestLocalID    string    `json:"request_local_id"`
	RequestService    string    `json:"request_service"`
	RequestPhone      string    `json:"request_phone"`
	RequestAmount     string    `json:"request_amount"`
	Status            string    `json:"status"`
	ErrorCode         int       `json:"error_code"`
	ErrorMsg          string    `json:"error_msg"`
	ResultStatus      string    `json:"result_status"`
	ResultRefNum      string    `json:"result_ref_num"`
	ResultService     string    `json:"result_service"`
	ResultDestination string    `json:"result_destination"`
	ResultAmount      int       `json:"result_amount"`
	ResultState       string    `json:"result_state"`
	IsChecked         bool      `json:"is_checked"`
	Client            string    `json:"client"`
}

func GetTransactions(ctx *gin.Context) {
	rows, err := db.DB.Query(context.Background(), GetTransactionsQuery, 0, 10)
	if err != nil {
		logger.Error(err)
	}
	defer rows.Close()

	transactions := make([]Transaction, 0)
	for rows.Next() {
		var trxn Transaction

		if err := rows.Scan(&trxn.UUID, &trxn.CreatedAt, &trxn.UpdatedAt, &trxn.RequestLocalID, &trxn.RequestService, &trxn.RequestPhone, &trxn.RequestAmount, &trxn.Status, &trxn.ErrorCode, &trxn.ErrorMsg, &trxn.ResultStatus, &trxn.ResultRefNum, &trxn.ResultService, &trxn.ResultDestination, &trxn.ResultAmount, &trxn.ResultState, &trxn.IsChecked, &trxn.Client); err != nil {
			logger.Errorf("row scan error %v", err)
		}
		transactions = append(transactions, trxn)
	}

	ctx.JSON(http.StatusOK, transactions)
}
