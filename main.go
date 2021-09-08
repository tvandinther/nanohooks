package main

import (
	"github.com/tvandinther/nanohooks/service"
	"github.com/tvandinther/nanohooks/supervisor"
	"os"
	"os/signal"
	"time"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go service.Start()
	go supervisor.Start()

	time.Sleep(5 * time.Second)
	supervisor.Test()

	for {
		select {
		case <-interrupt:
			return
		}
	}
}
