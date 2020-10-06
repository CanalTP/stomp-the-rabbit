package webstomp

import (
	"log"
	"time"

	"github.com/go-stomp/stomp"
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
	conn *stomp.Conn
	opts Options
	sub  *stomp.Subscription
}

type messageConsumer func([]byte)

func NewClient(opts Options) *Client {
	c := new(Client)
	c.opts = opts

	c.connect()

	return c
}

func (c *Client) connect() {
	for {
		log.Printf("connecting to stomp on %s\n", c.opts.Target)
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
					log.Println("connection stomp established!")
					return
				}
			}
		}
		logError("connection to stomp failed, retrying in 1 sec...", err)
		time.Sleep(1 * time.Second)
	}
}

func (c *Client) Consume(consumer messageConsumer) {
	for {
		msg := <-c.sub.C
		if msg != nil && msg.Err != nil {
			if c.sub.Active() {
				log.Printf("cannot handle message received, NACKing..., err: %v\n", msg.Err)
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
				log.Printf("failed to aknowledge message, err: %v\n", err)
			} else {
				consumer(msg.Body)
			}
		}
	}
}

func (c *Client) Disconnect() error {
	if c != nil {
		log.Println("closing stomp connection")
		return c.conn.Disconnect()
	}

	return nil
}

func logError(message string, err error) {
	if err != nil {
		log.Printf("%s: %s", message, err)
	}
}
