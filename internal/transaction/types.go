package transaction

import pb "gotoleg/rpc/gotoleg"

// All results are JSON objects encapsulated in below template
// "status" is one of the "SUCCESS", "ERROR", "EAUTH", "EXPIRED", "INVALID"
// "error-code" is optional, error code for error message
// "error-msg" is optional, error message in english
// "result" is optional, response object on success
type GarynjaResponse struct {
	Status       string `json:"status"`
	ErrorCode    int64  `json:"error-code"`
	ErrorMessage string `json:"error-msg"`
	Result       Result `json:"result"`
}
type Result struct {
	Status        string  `json:"status"`
	RefNum        int64   `json:"ref-num"`
	Service       string  `json:"service"`
	Destination   string  `json:"destination"`
	Amount        int64   `json:"amount"`
	State         string  `json:"state"`
	UpdateTS      float64 `json:"update-ts"`
	ReceivedTS    float64 `json:"received-ts"`
	TransactionTS float64 `json:"txn-ts"`
}

type Server struct {
	pb.UnimplementedTransactionServer
}
