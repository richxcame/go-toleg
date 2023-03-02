package handlers

import (
	"context"
	"fmt"
	"gotoleg/internal/db"
	"gotoleg/pkg/arrs"
	"gotoleg/pkg/logger"
	"gotoleg/web/entities"
	"gotoleg/web/helpers"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const GetTransactionsQuery = `SELECT uuid, created_at, updated_at, request_local_id, request_service, request_phone, request_amount, status, error_code, error_msg, result_status, result_ref_num, result_service, result_destination, result_amount, result_state, is_checked, client FROM transactions ORDER BY created_at DESC OFFSET $1 LIMIT $2`
const GetTransactionsCountQuery = `SELECT COUNT(*) FROM transactions`
const GetUserQuery = `SELECT username, password FROM users WHERE username = $1`

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

	urlQueries := ctx.Request.URL.Query()
	index := 0
	var values []interface{}
	var queries []string
	for k, v := range urlQueries {
		if arrs.Contains([]string{"uuid", "request_local_id", "request_service", "request_phone"}, k) {
			str := ""
			for _, v := range v {
				str += v + "|"
			}
			str = strings.TrimSuffix(str, "|")
			str += ""
			values = append(values, str)
			index++

			queries = append(queries, fmt.Sprintf("%s ~* $", k)+strconv.Itoa(index))
		}
	}
	valuesWithPagination := append(values, offset, limit)

	sqlStatement := "SELECT uuid, created_at, updated_at, request_local_id, request_service, request_phone, request_amount, status, error_code, error_msg, result_status, result_ref_num, result_service, result_destination, result_amount, result_state, is_checked, client FROM transactions"
	sqlFilters := ""
	if len(queries) > 0 {
		sqlFilters += " WHERE "
		sqlFilters += strings.Join(queries, " AND ")
	}
	sqlStatement += sqlFilters
	sqlStatement += " ORDER BY created_at DESC "
	sqlStatement += fmt.Sprintf(" offset $%v limit $%v", index+1, index+2)
	rows, err := db.DB.Query(context.Background(), sqlStatement, valuesWithPagination...)
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
	err = db.DB.QueryRow(context.Background(), "SELECT COUNT(*) FROM transactions"+sqlFilters, values...).Scan(&totalCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Couldn't count total number of transactions",
		})
		return
	}

	// Get all distinct status strings including filters
	statusRows, err := db.DB.Query(context.Background(), "SELECT DISTINCT(status) from transactions"+sqlFilters, values...)
	if err != nil {
		logger.Error(err)
	}
	defer statusRows.Close()

	statuses := make([]any, 0)
	for statusRows.Next() {
		val, err := statusRows.Values()
		if err != nil || len(val) < 1 {
			break
		}
		statuses = append(statuses, val[0])
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"statuses":     statuses,
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
	err := db.DB.QueryRow(context.Background(), GetUserQuery, user.Username).Scan(&dUser.Username, &dUser.Password)
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
