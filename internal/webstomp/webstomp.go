package webstomp

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/go-stomp/stomp"
	"golang.org/x/net/websocket"
)

type WebStompOpts struct {
	Protocol    string
	Login       string
	Passcode    string
	SendTimeout int
	RecvTimeout int
}

func Dial(target string, opts WebStompOpts) (*stomp.Conn, error) {
	u, err := url.Parse(target)
	if err != nil {
		log.Fatal(err)
	}
	origin, err := u.Parse("/")
	if err != nil {
		return nil, err
	}
	origin.Scheme = "https"
	c, err := websocket.Dial(u.String(), opts.Protocol, origin.String())

	if err != nil {
		return nil, fmt.Errorf("failed to open a new client connection to a the WebSocket %s, err: %v", target, err)
	}

	conn, err := stomp.Connect(c,
		stomp.ConnOpt.Login(opts.Login, opts.Passcode),
		stomp.ConnOpt.HeartBeat(time.Duration(opts.RecvTimeout)*time.Millisecond, time.Duration(opts.SendTimeout)*time.Millisecond),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect, err: %v", err)
	}

	return conn, nil
}
