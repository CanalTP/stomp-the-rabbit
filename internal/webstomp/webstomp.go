package webstomp

import (
	"log"
	"time"

	"github.com/go-stomp/stomp"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Protocol    string
	Login       string
	Passcode    string
	SendTimeout int
	RecvTimeout int
	Target      string
	Destination string
}

type Client struct {
	conn   *stomp.Conn
	opts   Options
	sub    *stomp.Subscription
	closed bool
	logger *logrus.Entry
}

type messageConsumer func([]byte)

func NewClient(opts Options, logger *logrus.Entry) *Client {
	c := new(Client)
	c.opts = opts
	c.logger = logger

	c.connect()

	return c
}

func (c *Client) connect() {
	for {
		if !c.closed {
			c.logger.Infof("connecting to stomp on %s\n", c.opts.Target)
			websocketconn, err := Dial(c.opts.Target, c.opts.Protocol)
			if err == nil {
				conn, err := stomp.Connect(websocketconn,
					stomp.ConnOpt.Login(c.opts.Login, c.opts.Passcode),
					stomp.ConnOpt.HeartBeat(time.Duration(c.opts.RecvTimeout)*time.Millisecond, time.Duration(c.opts.SendTimeout)*time.Millisecond),
				)
				if err == nil {
					c.conn = conn
					sub, err := c.conn.Subscribe(c.opts.Destination, stomp.AckClient)
					if err != stomp.ErrClosedUnexpectedly {
						c.sub = sub
						c.logger.Info("connection stomp established!")
						return
					}
				}
			}

			c.logError("connection to stomp failed, retrying in 1 sec...", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *Client) Consume(consumer messageConsumer) {
	for {
		msg := <-c.sub.C
		if msg != nil && msg.Err != nil {
			if c.sub.Active() {
				c.logger.Infof("cannot handle message received, NACKing..., err: %v\n", msg.Err)
				// Unacknowledge the message
				err := c.conn.Nack(msg)
				if err != nil {
					log.Printf("failed to unacknowledge the message, err: %v\n", err)
				}
			} else {
				c.connect()
			}
		} else {
			// Acknowledge the message
			err := c.conn.Ack(msg)
			if err != nil {
				c.logger.Infof("failed to aknowledge message, err: %v\n", err)
			} else {
				consumer(msg.Body)
			}
		}
	}
}

func (c *Client) Disconnect() error {
	if c != nil {
		c.logger.Info("closing stomp connection")
		c.closed = true
		return c.conn.Disconnect()
	}

	return nil
}

func (c *Client) logError(message string, err error) {
	if err != nil {
		c.logger.Errorf("%s: %s", message, err)
	}
}
