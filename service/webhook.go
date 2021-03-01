package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/tvandinther/nanohooks/service/models"
	"io"
	"log"
	"net/http"
)

type webhookService struct {
	cache cache
}

func newWebhookService() *webhookService {
	return &webhookService{
		cache: newCache(),
	}
}

func sendWebhook(job webhookJob, content interface{}) {
	body, err := json.Marshal(content)
	if err != nil {
		log.Print(err)
	}

	address := job.recipient.String()

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

func (w *webhookService) ReceiveAccount(confirmationReceipt *models.WebsocketConfirmationReceipt) error {
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

	accountSet, sendOk := w.cache.get(senderAccount)
	if sendOk {
		for _, job := range accountSet {
			go sendWebhook(job, newAccountTriggerView(confirmationReceipt, asOutgoing(), withAccounts(senderAccount, receiverAccount)))
		}
	}

	accountSet, receiveOk := w.cache.get(receiverAccount)
	if receiveOk {
		for _, job := range accountSet {
			go sendWebhook(job, newAccountTriggerView(confirmationReceipt, asIncoming(), withAccounts(receiverAccount, receiverAccount)))
		}
	}

	if !sendOk && !receiveOk {
		return errors.New("no accounts registered")
	}

	return nil
}
