package check

type CheckDestinationResult struct {
	Status      string `json:"status"`
	Reason      string `json:"reason,omitempty"`
	State       string `json:"state,omitempty"`
	Service     string `json:"service,omitempty"`
	Destination string `json:"destination,omitempty"`
}
type CheckDestinationResp struct {
	Status       string                 `json:"status"`
	ErrorCode    int                    `json:"error-code"`
	ErrorMessage string                 `json:"error-msg"`
	Result       CheckDestinationResult `json:"result,omitempty"`
}
