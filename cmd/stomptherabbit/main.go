package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CanalTP/stomptherabbit"
	"github.com/CanalTP/stomptherabbit/internal/amqp"
	"github.com/CanalTP/stomptherabbit/internal/webstomp"
)

func main() {
	c, err := config()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	logger := getLogger(stomptherabbit.Version, c.Logger.JSON)

	amqpClient := amqp.NewClient(c.AMQP.URL, c.AMQP.Exchange.Name, c.AMQP.ContentType, logger)
	defer amqpClient.Close()

	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		close(done)
	}()

	opts := webstomp.Options{
		Protocol:    c.Webstomp.Protocol,
		Login:       c.Webstomp.Login,
		Passcode:    c.Webstomp.Passcode,
		SendTimeout: c.Webstomp.SendTimeout,
		RecvTimeout: c.Webstomp.RecvTimeout,
		Target:      c.Webstomp.Target,
		Destination: c.Webstomp.Destination,
	}

	var wsClient *webstomp.Client
	go func() {
		wsClient = webstomp.NewClient(opts, logger)
		wsClient.Consume(func(msg []byte) {
			amqpClient.Send(msg)
		})
	}()

	fmt.Println(c.ToString())
	logger.Println("Waiting for messages...")
	<-done
	logger.Println("Gracefully exiting...")
	wsClient.Disconnect()
}
