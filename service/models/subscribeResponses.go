package models

type Block struct {
	Type           string `json:"type"`
	Account        string `json:"account"`
	Previous       string `json:"previous"`
	Representative string `json:"representative"`
	Balance        string `json:"balance"`
	Link           string `json:"link"`
	LinkAsAccount  string `json:"link_as_account"`
	Signature      string `json:"signature"`
	Work           string `json:"work"`
	Subtype        string `json:"subtype"`
}

type ConfirmationMessage struct {
	Account          string `json:"account"`
	Amount           string `json:"amount"`
	Hash             string `json:"hash"`
	ConfirmationType string `json:"confirmation_type"`
	Block            Block  `json:"block"`
}

type WebsocketConfirmationReceipt struct {
	Topic string `json:"topic"`
	Time string `json:"time"`
	Message ConfirmationMessage `json:"message"`
}
