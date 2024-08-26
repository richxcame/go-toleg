package handlers

import (
	"encoding/json"
	"fmt"
	"gotoleg/internal/config"
	"gotoleg/internal/constants"
	"gotoleg/internal/db"
	"gotoleg/internal/transaction"
	"gotoleg/internal/utility"
	"gotoleg/pkg/arrs"
	"gotoleg/pkg/hmacsha1"
	"gotoleg/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Transaction struct {
	LocalID string `json:"local_id"`
	Service string `json:"service"`
	Phone   string `json:"phone"`
	Amount  string `json:"amount"`
	Note    string `json:"note"`
	ApiKey  string `json:"api_key"`
	Reason  string `json:"reason"`
}

func AddTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	// bind request body
	var trxn Transaction
	err := c.BindJSON(&trxn)
	if err != nil {
		logger.Errorf("couldn't bind request body: %v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Couldn't bind request body",
		})
		return
	}

	// validate api key
	client, hasInList := arrs.HasMapWithKey(config.Clients, trxn.ApiKey)
	if !hasInList {
		logger.Errorf("wrong api_key: %v, %v", trxn, client)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid_api_key",
			"message": "API key is not valid",
		})
		return
	}

	// Convert amount to tenges
	mnt, err := strconv.ParseFloat(trxn.Amount, 64)
	if err != nil {
		logger.Errorf("couldn't parse request amount to float: %v", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   err.Error(),
			"message": "Couldn't convert request amount to float",
		})
		return
	}
	amount := strconv.Itoa(int(mnt))

	// Insert request to database
	_uuid := uuid.New().String()
	sqlStatement := `
		INSERT INTO transactions (uuid, created_at, updated_at, client, request_local_id, request_service, request_phone, request_amount, note, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
	_, err = db.DB.Exec(ctx, sqlStatement, _uuid, time.Now(), time.Now(), client, trxn.LocalID, trxn.Service, trxn.Phone, amount, trxn.Note, trxn.Reason)
	if err != nil {
		logger.Errorf("Couldn't save request to database: %v", err.Error())
	}

	// Get epoch time
	epochTime, err := utility.GetEpoch()
	if err != nil {
		logger.Error("Couldn't get epoch time: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Coudn't get epoch time",
		})
		return
	}

	// Prepare ts, msg and request body
	ts := strconv.FormatInt(epochTime, 10)
	msg := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", trxn.LocalID, trxn.Service, amount, trxn.Phone, ts, ts, constants.USERNAME)
	data := url.Values{
		"local-id":    {trxn.LocalID},
		"service":     {trxn.Service},
		"amount":      {amount},
		"destination": {trxn.Phone},
		"txn-ts":      {ts},
		"ts":          {ts},
		"hmac":        {hmacsha1.Generate(constants.ACCESS_TOKEN, msg)},
	}

	// send request
	httpClient := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := httpClient.PostForm(constants.ADD_TRANSACTION_URL, data)
	if err != nil {
		logger.Errorf("Couldn't send transaction: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't send transaction",
		})
		return
	}
	defer resp.Body.Close()

	// read response body
	respInBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Couldn't read response body: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't read response body",
		})
		return
	}

	// unmarshall response
	var result transaction.GarynjaResponse
	err = json.Unmarshal(respInBytes, &result)
	if err != nil {
		logger.Errorf("Couldn't unmarshall response: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't read response body",
		})
		return
	}

	// update the transaction
	const sqlStr = `UPDATE transactions SET status=$1, error_code=$2, error_msg=$3, result_status=$4, result_ref_num=$5, result_service=$6, result_destination=$7, result_amount=$8, result_state=$9, updated_at=$10 WHERE uuid=$11`
	_, err = db.DB.Exec(ctx, sqlStr, result.Status, result.ErrorCode, result.ErrorMessage, result.Result.Status, result.Result.RefNum, result.Result.Service, result.Result.Destination, result.Result.Amount, result.Result.State, time.Now(), _uuid)
	if err != nil {
		logger.Errorf("Couldn't update the transaction: %v", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Transaction added successfully",
		"data":    result,
	})
}
