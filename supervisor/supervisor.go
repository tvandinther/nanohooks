package supervisor

import (
	"github.com/golang/protobuf/proto"
	"github.com/streadway/amqp"
	spec "github.com/tvandinther/nanohooks/proto"
	"log"
	"net/url"
)

type RabbitMsg struct {
	QueueName string                     `json:"queueName"`
	Message   spec.WebhookJobMessage `json:"message"`
}

var pchan = make(chan RabbitMsg, 10)
var rabbitURL string = "amqp://admin:admin@localhost:5672"

func Start() {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-pchan:
			data, err := proto.Marshal(&msg.Message)
			if err != nil {
				log.Println(err)
				continue
			}

			err = amqpChannel.Publish(
				"",
				msg.QueueName,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body: data,
				},
			)
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("INFO: published msg: %v\n", msg.Message.Uid)
		}
	}
}

func Test() {
	recipient := url.URL{
		Scheme:      "http",
		Host:        "localhost:8080",
	}

	pchan <-RabbitMsg{
		QueueName: "webhooks",
		Message:   spec.WebhookJobMessage{
			Uid:        "1",
			Action:     spec.WebhookJobAction_CREATE,
			WebhookJob: &spec.WebhookJob{
				Id:        "1",
				Accounts:  []string{
					"nano_3xaz74n68af4oa9jfn8kuan44xz1j5nr69ztt7qo8bu1wgqns9upcfntgkc7",
					"nano_1byxifump4h4999rqbe7ewa1ysc3hm88wuqmprcxi9ko8qaijsds3zzafeef",
				},
				Recipient: recipient.String(),
			},
			ReplyTo:    "",
		},
	}
}
