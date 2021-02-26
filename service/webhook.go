package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type webhookService struct {
	cache map[string][]url.URL
}

func newWebhookService() *webhookService {
	return &webhookService{
		map[string][]url.URL{},
	}
}

func sendWebhook(recipient url.URL, content interface{}) {
	body, err := json.Marshal(content)
	if err != nil {
		log.Print(err)
	}

	address := recipient.String()

	res, err := http.Post(address, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Println(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		log.Printf("Webhook recipient failed to confirm receipt: %s\n", address)
	}
}

func (w *webhookService) addToCache(account string, recipient url.URL) {
	recipients, ok := w.cache[account]
	if ok {
		w.cache[account] = append(recipients, recipient)
	} else {
		w.cache[account] = []url.URL{recipient}
	}
}

func (w *webhookService) removeFromCache(account string, recipient url.URL) error {
	 recipients, ok := w.cache[account]
	 if !ok {
	 	return errors.New("account not registered")
	 }
	 fmt.Println("Cache deletion not implemented: ", recipients)
	 //TODO: Add logic to remove from the cache
	 return nil
}

type block struct {
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

type confirmationMessage struct {
	Account          string `json:"account"`
	Amount           string `json:"amount"`
	Hash             string `json:"hash"`
	ConfirmationType string `json:"confirmation_type"`
	Block            block  `json:"block"`
}

func (w *webhookService) ReceiveAccount(confirmationReceipt *websocketConfirmationReceipt) error {
	message := confirmationReceipt.Message
	confirmationType := message.Block.Subtype
	var senderAccount, receiverAccount string
	var accounts []string
	if confirmationType == "send" {
		senderAccount = message.Account
		receiverAccount = message.Block.LinkAsAccount

		accounts = append(accounts, senderAccount, receiverAccount)
	} else if confirmationType == "receive" {
		senderAccount = ""
		receiverAccount = message.Account

		accounts = append(accounts, senderAccount, receiverAccount)
	} else {
		log.Printf("Unknown confirmation type: %s\n", confirmationType)
	}

	recipients, sendOk := w.cache[senderAccount]
	if sendOk {
		for _, recipient := range recipients {
			go sendWebhook(recipient, newAccountTriggerView(confirmationReceipt, asOutgoing(), withAccounts(senderAccount, receiverAccount)))
		}
	}

	recipients, receiveOk := w.cache[receiverAccount]
	if receiveOk {
		for _, recipient := range recipients {
			go sendWebhook(recipient, newAccountTriggerView(confirmationReceipt, asIncoming(), withAccounts(receiverAccount, receiverAccount)))
		}
	}

	if !sendOk && !receiveOk {
		return errors.New("no accounts registered")
	}

	return nil
}
