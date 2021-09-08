package service

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	spec "github.com/tvandinther/nanohooks/proto"
	"log"
)

type listener struct {
	rabbitMQUrl string
	replyChannel chan RabbitMsg
	registrar registrar
	webhookService *webhookService
}

type RabbitMsg struct {
	QueueName string                   `json:"queueName"`
	Reply     spec.WebhookJobReply `json:"reply"`
}

func newListener(r registrar, w *webhookService) *listener {
	return &listener{
		rabbitMQUrl: "amqp://admin:admin@localhost:5672",
		registrar:      r,
		webhookService: w,
		replyChannel: make(chan RabbitMsg, 10),
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

	go l.startReplyService()

	for {
		select {
		case msg := <-msgChannel:
			hookMsg := &spec.WebhookJobMessage{}
			err := proto.Unmarshal(msg.Body, hookMsg)
			if err != nil {
				log.Println(err)
				continue
			}

			err = msg.Ack(true)
			if err != nil {
				log.Println(err)
			}
			
			l.handleMsg(hookMsg)
		}
	}
}

func (l *listener) handleMsg(msg *spec.WebhookJobMessage) {
	log.Println("Received message: ", msg.Uid)

	AddAccountTrigger(l.registrar, l.webhookService, newWebhookJobFromProto(msg.WebhookJob))

	reply := spec.WebhookJobReply{
		Uid:    msg.Uid,
		Status: spec.ActionResult_SUCCESSFUL,
	}

	replyMsg := RabbitMsg{
		QueueName: msg.ReplyTo,
		Reply: reply,
	}

	l.replyChannel <-replyMsg
}

func (l *listener) startReplyService() {
	conn, err := amqp.Dial(l.rabbitMQUrl)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-l.replyChannel:
			data, err := proto.Marshal(&msg.Reply)
			if err != nil {
				log.Printf("ERROR: fail marshal: %s", err.Error())
				continue
			}

			err = amqpChannel.Publish(
				"",
				msg.QueueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        data,
				},
			)
			if err != nil {
				log.Printf("ERROR: fail publish msg: %s", err.Error())
				continue
			}
		}
	}
}
