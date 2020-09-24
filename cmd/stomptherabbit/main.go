package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CanalTP/stomptherabbit/internal/rabbitmq"
	"github.com/CanalTP/stomptherabbit/internal/webstomp"
	"github.com/go-stomp/stomp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	c, err := config()
	if err != nil {
		failOnError(err, "failed to load configuration")
	}

	opts := webstomp.WebStompOpts{
		Protocol:    c.Webstomp.Protocol,
		Login:       c.Webstomp.Login,
		Passcode:    c.Webstomp.Passcode,
		SendTimeout: c.Webstomp.SendTimeout,
		RecvTimeout: c.Webstomp.RecvTimeout,
	}
	destination := c.Webstomp.Destination
	target := c.Webstomp.Target

	exchangeName := c.RabbitMQ.Exchange.Name

	m := rabbitmq.NewAmqpManager(c.RabbitMQ.URL, exchangeName, c.RabbitMQ.ContentType)
	defer m.Close()

	done := make(chan struct{})

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		close(done)
	}()

	var conn *stomp.Conn
	go func() {
		for {
			conn, err = webstomp.Dial(target, opts)
			if err != nil {
				log.Printf("failed to connect to webstomp server, err: %v\n", err)
				continue
			}
			sub, err := conn.Subscribe(destination, stomp.AckClient)
			if err != nil {
				log.Printf("failed to subscribe to %s, err: %v\n", destination, err)
				continue
			}
			for {
				msg := <-sub.C
				if msg.Err != nil {
					log.Printf("cannot handle message received, NACKing..., err: %v\n", msg.Err)
					// Unacknowledge the message
					err = conn.Nack(msg)
					if err != nil {
						log.Printf("failed to unacknowledge the message, err: %v\n", err)
					}
				}

				m.Send(msg.Body)

				// Acknowledge the message
				err = conn.Ack(msg)
				if err != nil {
					log.Printf("failed to aknowledge message, err: %v\n", err)
				}
			}
		}
	}()
	fmt.Println(c.ToString())
	fmt.Println("Waiting for messages...")
	<-done
	fmt.Println("Gracefully exiting...")
	conn.Disconnect()
}
