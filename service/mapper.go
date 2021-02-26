package service

import (
	"log"
	"strconv"
	"strings"
)

type accountTriggerViewOption func(*AccountTriggerView)

func newAccountTriggerView(receipt *websocketConfirmationReceipt, opts ...accountTriggerViewOption) AccountTriggerView {
	message := receipt.Message

	time, err := strconv.ParseInt(receipt.Time, 10, 64)
	if err != nil {
		log.Println(err)
	}

	amount, err := strconv.ParseFloat(convertRawToMNano(message.Amount), 32)
	balance, err := strconv.ParseFloat(convertRawToMNano(message.Block.Balance), 32)

	v := AccountTriggerView{
		Time: time,
		Type:    message.Block.Subtype,
		Amount:  amount,
		RawAmount: message.Amount,
		Balance: balance,
		RawBalance: message.Block.Balance,
		Hash:    message.Hash,
	}

	//Apply functional options
	for _, opt := range opts {
		opt(&v)
	}

	return v
}

func convertRawToMNano(s string) string {
	//integerPlaces := 9
	decimalPlaces := 30
	length := len(s)

	if length <= decimalPlaces {
		return "0." + strings.Repeat("0", decimalPlaces - length) + s
	} else {
		decimalIndex := length - decimalPlaces
		return s[:decimalIndex] + "." + s[decimalIndex:]
	}
}

func asIncoming() accountTriggerViewOption {
	return func(v *AccountTriggerView) {
		v.Transaction = "incoming"
	}
}

func asOutgoing() accountTriggerViewOption {
	return func(v *AccountTriggerView) {
		v.Transaction = "outgoing"
	}
}

func withAccounts(sendingAccount, receivingAccount string) accountTriggerViewOption {
	return func(v *AccountTriggerView) {
		v.SendingAccount = sendingAccount
		v.ReceivingAccount = receivingAccount
	}
}
