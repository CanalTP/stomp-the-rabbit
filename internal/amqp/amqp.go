package amqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Client struct {
	url          string
	exchangeName string
	contentType  string

	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	closed       bool
}

func NewClient(url, exchangeName, contentType string) *Client {
	m := new(Client)
	m.url = url
	m.exchangeName = exchangeName
	m.contentType = contentType

	m.connect()
	go m.reconnector()

	return m
}

func (m *Client) Close() {
	log.Println("closing amqp broker connection")
	m.closed = true
	m.channel.Close()
	m.connection.Close()
}

func (m *Client) connect() {
	for {
		log.Printf("connecting to amqp broker on %s\n", m.url)
		conn, err := amqp.Dial(m.url)
		if err == nil {
			m.connection = conn
			m.errorChannel = make(chan *amqp.Error)
			m.connection.NotifyClose(m.errorChannel)

			log.Println("connection amqp broker established!")

			m.openChannel()
			m.declareExchange()

			return
		}
		logError("connection to amqp broker failed, retrying in 1 sec...", err)
		time.Sleep(1 * time.Second)
	}
}

func (m *Client) reconnector() {
	for {
		err := <-m.errorChannel
		if !m.closed {
			logError("reconnecting after connection closed", err)
			m.connect()
		}
	}
}

func (m *Client) openChannel() {
	channel, err := m.connection.Channel()
	logError("opening channel failed", err)
	m.channel = channel
}

func (m *Client) declareExchange() {
	err := m.channel.ExchangeDeclare(
		m.exchangeName, // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	logError("Failed to declare an exchange", err)
}

func (m *Client) Send(message []byte) {
	err := m.channel.Publish(
		m.exchangeName, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  m.contentType,
			Body:         message,
		},
	)
	logError("failed to publish a message", err)
}

func logError(message string, err error) {
	if err != nil {
		log.Printf("%s: %s", message, err)
	}
}
