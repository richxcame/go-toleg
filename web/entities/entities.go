package entities

import (
	"time"

	"github.com/google/uuid"
)

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
	ResultReason      string    `json:"result_reason"`
	IsChecked         bool      `json:"is_checked"`
	Client            string    `json:"client"`
}

type User struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Password  string    `json:"password"`
}
