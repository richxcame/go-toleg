package handlers

import (
	"context"
	"gotoleg/internal/db"
	"gotoleg/pkg/logger"
	"gotoleg/web/entities"
	"gotoleg/web/helpers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const GetTransactionsQuery = `SELECT uuid, created_at, updated_at, request_local_id, request_service, request_phone, request_amount, status, error_code, error_msg, result_status, result_ref_num, result_service, result_destination, result_amount, result_state, is_checked, client FROM transactions ORDER BY created_at DESC OFFSET $1 LIMIT $2`
const GetTransactionsCountQuery = `SELECT COUNT(*) FROM transactions`
const GetUserQuery = `SELECT username, created_at, update_at, password FROM users WHERE username = $1`

func GetTransactions(ctx *gin.Context) {
	offsetQuery := ctx.DefaultQuery("offset", "0")
	limitQuery := ctx.DefaultQuery("limit", "20")
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "offset value must be convertable to integer",
		})
		return
	}
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "limit value must be convertable to integer",
		})
		return
	}

	rows, err := db.DB.Query(context.Background(), GetTransactionsQuery, offset, limit)
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

	totalCount := 0
	err = db.DB.QueryRow(context.Background(), GetTransactionsCountQuery).Scan(&totalCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't count total number of transactions",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"offset":       offset,
		"limit":        limit,
		"total_count":  totalCount,
	})
}

func Login(ctx *gin.Context) {
	var user entities.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Couln't parse the request to user",
		})
		return
	}

	var dUser entities.User
	err := db.DB.QueryRow(context.Background(), GetUserQuery, user.Username).Scan(&dUser.Username, &dUser.CreatedAt, &dUser.UpdatedAt, &dUser.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Couldn't find user",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "wrong_password",
			"message": "Invalid password",
		})
		return
	}

	tokens, err := helpers.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Coulnd't create token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func RefreshToken(ctx *gin.Context) {
	ctx.JSON(200, "success login")
}
