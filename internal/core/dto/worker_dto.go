package dto

type MessageData struct {
	ID        string `json:"id"`
	Processed string `json:"processed"`
	Attempt   int    `json:"attempts"`
}

type ResponseData struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ProcessResult struct {
	Retry     bool   `json:"retry"`
	Status    string `json:"status"`
	Processed bool   `json:"processed"`
}

type ImportTxnMessage struct {
	JobID    string `json:"job_id"`
	Filename string `json:"filename"`
	Action   string `json:"action"`
	Data     struct {
		Transaction TransactionRequest `json:"transaction"`
	} `json:"data"`
}
