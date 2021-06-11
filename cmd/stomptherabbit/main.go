package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/CanalTP/stomptherabbit"
	"github.com/CanalTP/stomptherabbit/internal/amqp"
	"github.com/CanalTP/stomptherabbit/internal/api"
	"github.com/CanalTP/stomptherabbit/internal/scoreboard"
	"github.com/CanalTP/stomptherabbit/internal/webstomp"
)

func main() {
	c, err := config()
	if err != nil {
		log.Fatalf("failed to load configuration: %s", err)
	}

	logger := getLogger(stomptherabbit.Version, c.Logger.JSON)
	scoreboard := scoreboard.NewScoreBoard()

	opts := webstomp.Options{
		Protocol:    c.Webstomp.Protocol,
		Login:       c.Webstomp.Login,
		Passcode:    c.Webstomp.Passcode,
		SendTimeout: c.Webstomp.SendTimeout,
		RecvTimeout: c.Webstomp.RecvTimeout,
		Target:      c.Webstomp.Target,
		Destination: c.Webstomp.Destination,
	}

	// Lanch web service for supervision : Listen port 8080
	router := api.Router(
		api.NewStatusHandler(
			c.AMQP.URL,
			opts,
			stomptherabbit.Version,
			"Stomptherabbit",
			&scoreboard))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()

	amqpClient := amqp.NewClient(c.AMQP.URL, c.AMQP.Exchange.Name, c.AMQP.ContentType, logger)
	defer amqpClient.Close()
	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		close(done)
	}()
	var wsClient *webstomp.Client
	go func() {
		wsClient = webstomp.NewClient(opts, logger)
		wsClient.Consume(func(msg []byte) {
			scoreboard.Set("dateLastMessageAttempSendRabbitMQ")
			amqpClient.Send(msg, &scoreboard)
		}, &scoreboard)
	}()

	fmt.Println(c.ToString())
	logger.Println("Waiting for messages...")
	<-done
	logger.Println("Gracefully exiting...")
	wsClient.Disconnect()
}
