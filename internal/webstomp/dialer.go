package webstomp

import (
	"net"
	"net/url"

	"golang.org/x/net/websocket"
)

// Dial creates a client connection to the given target
func Dial(target string, protocol string) (net.Conn, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	origin, err := u.Parse("/")
	if err != nil {
		return nil, err
	}
	origin.Scheme = "https"
	return websocket.Dial(u.String(), protocol, origin.String())
}
