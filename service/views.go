package service

type AccountTriggerView struct {
	Time             int64   `json:"time"`
	Transaction      string	 `json:"transaction"`
	Type             string  `json:"type"`
	SendingAccount   string  `json:"sending_account"`
	ReceivingAccount string  `json:"receiving_account"`
	Amount           float64 `json:"amount"`
	RawAmount        string  `json:"raw_amount"`
	Balance          float64 `json:"balance"`
	RawBalance       string  `json:"raw_balance"`
	Hash             string  `json:"hash"`
}
