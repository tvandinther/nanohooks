package service

import (
	"github.com/gorilla/websocket"
	"github.com/tvandinther/nanohooks/service/models"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type options = struct {
	Accounts []string `json:"accounts,omitempty"`
}

type message = struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Ack     bool   `json:"ack,omitempty"`
	Options interface{} `json:"options,omitempty"`
}

func Start() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	address := url.URL{
		Scheme:      "ws",
		Host:        "localhost:7078",
	}

	conn, _, err := websocket.DefaultDialer.Dial(address.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer func(conn *websocket.Conn) {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}(conn)

	done := make(chan struct{})

	var registrar registrar
	registrar = newWebsocketRegistrar(conn)

	var webhookService *webhookService
	webhookService = newWebhookService()

	go func() {
		defer close(done)
		for {
			payload := models.WebsocketConfirmationReceipt{}
			err := conn.ReadJSON(&payload)
			if err != nil {
				log.Println("read:", err)
				return
			}

			switch payload.Topic {
			case "confirmation":
				if err != nil {
					log.Println("invalid message")
				} else {
					err = webhookService.ReceiveAccount(&payload)
					if err != nil {
						registrar.register("confirmation", models.UpdateAccountOptions{
							AccountsDel: []string{payload.Message.Account},
						})
					}
				}
			default:
				log.Printf("No route for topic: %s\n", payload.Topic)
			}
		}
	}()

	recipient := url.URL{
		Scheme:      "http",
		Host:        "localhost:8080",
	}
	account := "nano_3xaz74n68af4oa9jfn8kuan44xz1j5nr69ztt7qo8bu1wgqns9upcfntgkc7"

	job := webhookJob{
		id: "1",
		accounts: []string{account},
		recipient: recipient,
	}

	AddAccountTrigger(registrar, webhookService, job)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupted")

			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}
