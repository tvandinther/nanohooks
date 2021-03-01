package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	spec "github.com/tvandinther/nanohooks/proto"
	"log"
)

type listener struct {
	rabbitMQUrl string
	registrar *registrar
	webhookService *webhookService
}

func newListener(r *registrar, w *webhookService) *listener {
	return &listener{
		registrar:      r,
		webhookService: w,
	}
}

func (l *listener) listen() {
	conn, err := amqp.Dial(l.rabbitMQUrl)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	queue, err := amqpChannel.QueueDeclare(
		"webhooks",
		true,
		false,
		true,
		false,
		nil,
		)
	if err != nil {
		log.Fatal(err)
	}

	msgChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		true,
		false,
		false,
		nil,
		)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-msgChannel:
			hookMsg := &spec.CreateWebhookJobMessage{}
			err := proto.Unmarshal(msg.Body, hookMsg)
			if err != nil {
				log.Println(err)
				continue
			}

			err = msg.Ack(true)
			if err != nil {
				log.Println(err)
			}
			
			handleMsg(hookMsg)
		}
	}
}

func handleMsg(msg *spec.CreateWebhookJobMessage) {
	
}
