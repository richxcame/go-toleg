package transaction

type TransactionResult struct {
	RefNum        int64   `json:"ref-num,omitempty"`
	Service       string  `json:"service,omitempty"`
	Destination   string  `json:"destination,omitempty"`
	Amount        int     `json:"amount,omitempty"`
	State         string  `json:"state,omitempty"`
	UpdateTS      float64 `json:"update-ts,omitempty"`
	ReceivedTS    float64 `json:"received-ts,omitempty"`
	TransactionTS float64 `json:"txn-ts,omitempty"`
}
type TransactionResp struct {
	Status       string            `json:"status,omitempty"`
	ErrorCode    int               `json:"error-code,omitempty"`
	ErrorMessage string            `json:"error-msg,omitempty"`
	Result       TransactionResult `json:"result,omitempty"`
}
