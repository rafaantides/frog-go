package dto

type ImportTxnMessage struct {
	JobID    string `json:"job_id"`
	Filename string `json:"filename"`
	Action   string `json:"action"`
	Data     struct {
		Transaction TransactionRequest `json:"transaction"`
	} `json:"data"`
}
