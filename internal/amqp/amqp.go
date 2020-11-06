package amqp

import (
	"time"

	"github.com/sirupsen/logrus"
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
	logger       *logrus.Entry
}

func NewClient(url, exchangeName, contentType string, logger *logrus.Entry) *Client {
	c := new(Client)
	c.url = url
	c.exchangeName = exchangeName
	c.contentType = contentType
	c.logger = logger

	c.connect()
	go c.reconnector()

	return c
}

func (c *Client) Close() {
	c.logger.Info("closing amqp broker connection")
	c.closed = true
	c.channel.Close()
	c.connection.Close()
}

func (c *Client) connect() {
	for {
		c.logger.Infof("connecting to amqp broker on %s\n", c.url)
		conn, err := amqp.Dial(c.url)
		if err == nil {
			c.connection = conn
			c.errorChannel = make(chan *amqp.Error)
			c.connection.NotifyClose(c.errorChannel)

			c.logger.Info("connection amqp broker established!")

			c.openChannel()
			c.declareExchange()

			return
		}
		c.logError("connection to amqp broker failed, retrying in 1 sec...", err)
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) reconnector() {
	for {
		err := <-c.errorChannel
		if !c.closed {
			c.logError("reconnecting after connection closed", err)
			c.connect()
		}
	}
}

func (c *Client) openChannel() {
	channel, err := c.connection.Channel()
	c.logError("opening channel failed", err)
	c.channel = channel
}

func (c *Client) declareExchange() {
	err := c.channel.ExchangeDeclare(
		c.exchangeName, // name
		"fanout",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)

	c.logError("Failed to declare an exchange", err)
}

func (c *Client) Send(message []byte) {
	err := c.channel.Publish(
		c.exchangeName, // exchange
		"",             // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  c.contentType,
			Body:         message,
		},
	)
	c.logError("failed to publish a message", err)
}

func (c *Client) logError(message string, err error) {
	if err != nil {
		c.logger.Errorf("%s: %s", message, err)
	}
}
