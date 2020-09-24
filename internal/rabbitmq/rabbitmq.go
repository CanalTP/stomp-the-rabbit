package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type AmqpManager struct {
	url          string
	exchangeName string
	contentType  string

	errorChannel chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	closed       bool
}

func NewAmqpManager(url, exchangeName, contentType string) *AmqpManager {
	m := new(AmqpManager)
	m.url = url
	m.exchangeName = exchangeName
	m.contentType = contentType

	m.connect()
	go m.reconnector()

	return m
}

func (m *AmqpManager) Close() {
	log.Println("closing rabbitmq connection")
	m.channel.Close()
	m.connection.Close()
}

func (m *AmqpManager) connect() {
	for {
		log.Printf("connecting to rabbitmq on %s\n", m.url)
		conn, err := amqp.Dial(m.url)
		if err == nil {
			m.connection = conn
			m.errorChannel = make(chan *amqp.Error)
			m.connection.NotifyClose(m.errorChannel)

			log.Println("connection rabbitmq established!")

			m.openChannel()
			m.declareExchange()

			return
		}
		logError("connection to rabbitmq failed, retrying in 1 sec...", err)
		time.Sleep(1 * time.Second)
	}
}

func (m *AmqpManager) reconnector() {
	for {
		err := <-m.errorChannel
		if !m.closed {
			logError("reconnnecting after connection closed", err)
			m.connect()
		}
	}
}

func (m *AmqpManager) openChannel() {
	channel, err := m.connection.Channel()
	logError("opening channel failed", err)
	m.channel = channel
}

func (m *AmqpManager) declareExchange() {
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

func (m *AmqpManager) Send(message []byte) {
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
