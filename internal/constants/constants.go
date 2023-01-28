package constants

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var BASE_URL string
var USERNAME string
var SERVER string
var ACCESS_TOKEN string

var EPOCH_URL string
var BALANCE_URL string
var SERVICES_URL string

var ADD_TRANSACTION_URL string
var TRANSACTION_STATUS_URL string
var TRANSACTION_RETRY_URL string
var CORRECT_DECLINED_URL string
var FORCE_ADD_URL string
var CALLBACK_URL string

var CHECK_DESTINATION_URL string
var POLL_CHECK_DESTINATION_URL string
var CHECK_CALLBACK_URL string

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	BASE_URL = os.Getenv("GARYNJA_BASE_URL")
	USERNAME = os.Getenv("GARYNJA_USERNAME")
	SERVER = os.Getenv("GARYNJA_SERVER")
	ACCESS_TOKEN = os.Getenv("GARYNJA_ACCESS_TOKEN")

	EPOCH_URL = fmt.Sprintf("%s/epoch", BASE_URL)
	BALANCE_URL = fmt.Sprintf("%s/%s/%s/dealer/balance", BASE_URL, USERNAME, SERVER)
	SERVICES_URL = fmt.Sprintf("%s/%s/%s/dealer/services", BASE_URL, USERNAME, SERVER)

	ADD_TRANSACTION_URL = fmt.Sprintf("%s/%s/%s/txn/add", BASE_URL, USERNAME, SERVER)
	TRANSACTION_STATUS_URL = fmt.Sprintf("%s/%s/%s/txn/info", BASE_URL, USERNAME, SERVER)
	TRANSACTION_RETRY_URL = fmt.Sprintf("%s/%s/%s/txn/retry", BASE_URL, USERNAME, SERVER)
	CORRECT_DECLINED_URL = fmt.Sprintf("%s/%s/%s/txn/change-service", BASE_URL, USERNAME, SERVER)
	FORCE_ADD_URL = fmt.Sprintf("%s/%s/%s/txn/force/add", BASE_URL, USERNAME, SERVER)
	CALLBACK_URL = fmt.Sprintf("%s/%s/%s/<callback>", BASE_URL, USERNAME, SERVER)

	CHECK_DESTINATION_URL = fmt.Sprintf("%s/%s/%s/cd/add", BASE_URL, USERNAME, SERVER)
	POLL_CHECK_DESTINATION_URL = fmt.Sprintf("%s/%s/%s/cd/poll", BASE_URL, USERNAME, SERVER)
	CHECK_CALLBACK_URL = fmt.Sprintf("%s/%s/%s/<callback>", BASE_URL, USERNAME, SERVER)
}
