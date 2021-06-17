package webstomp

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/streadway/amqp"
)

type responseStatus struct {
	ApplicationName             string    `json:"applicationName"`
	ApplicationVersion          string    `json:"applicationVersion"`
	Hostname                    string    `json:"hostname"`
	AmqpHostname                string    `json:"amqpHostname"`
	WebstompHostname            string    `json:"webstompHostname"`
	GoRuntimeVersion            string    `json:"goRuntimeVersion"`
	IsAmqpConnActive            bool      `json:"isAmqpConnActive"`
	IsWebStompConnActive        bool      `json:"isWebStompConnActive"`
	LastMessageReceiveSNCF      time.Time `json:"lastMessageReceiveSNCF"`
	LastAttempSendRabbitMQ      time.Time `json:"lastAttempSendRabbitMQ"`
	LastSuccessfulWriteRabbitMQ time.Time `json:"lastSuccessfulWriteRabbitMQ"`
}

type datesScoreBoard struct {
	name      string
	timestamp time.Time
}

func newResponseStatus(applicationName, applicationVersion, hostname string) *responseStatus {
	return &responseStatus{
		ApplicationName:    applicationName,
		ApplicationVersion: applicationVersion,
		GoRuntimeVersion:   runtime.Version(),
		Hostname:           hostname,
	}
}

func (r *responseStatus) setStatusAmqp(amqpURL string) {
	u, _ := url.Parse(amqpURL)
	r.AmqpHostname = u.Hostname()
	r.IsAmqpConnActive = false
	if conn, err := amqp.Dial(amqpURL); err == nil {
		r.IsAmqpConnActive = true
		conn.Close()
	}
}
func (r *responseStatus) setStatusStomp(stompOpts Options) {
	r.IsWebStompConnActive = false
	r.WebstompHostname = stompOpts.Target
	if websocketconn, err := Dial(stompOpts.Target, stompOpts.Protocol); err == nil {
		r.IsWebStompConnActive = true
		websocketconn.Close()
	}
}

func (r *responseStatus) setStatusScoreBoard(mapScore map[string]time.Time) {
	r.LastMessageReceiveSNCF = mapScore["lastMessageReceiveSNCF"]
	r.LastAttempSendRabbitMQ = mapScore["lastAttempSendRabbitMQ"]
	r.LastSuccessfulWriteRabbitMQ = mapScore["lastSuccessfulWriteRabbitMQ"]
}

// StatusHandler : web handler for stomptherabbit status
type StatusHandler struct {
	amqpURL            string
	webstompOptions    Options
	applicationVersion string
	applicationName    string
	scoreboard         *ScoreBoard
}

//NewStatusHandler : constuctor for statusHandler
func NewStatusHandler(amqpURL string, webstompOptions Options, version, applicationName string, scoreboard *ScoreBoard) *StatusHandler {
	return &StatusHandler{
		amqpURL:            amqpURL,
		webstompOptions:    webstompOptions,
		applicationVersion: version,
		applicationName:    applicationName,
		scoreboard:         scoreboard,
	}
}

func (h *StatusHandler) getStatus() *responseStatus {
	hostname, _ := os.Hostname()
	response := newResponseStatus(h.applicationName, h.applicationVersion, hostname)
	response.setStatusAmqp(h.amqpURL)
	response.setStatusStomp(h.webstompOptions)
	response.setStatusScoreBoard(h.scoreboard.All())
	return response
}

func (h *StatusHandler) status(w http.ResponseWriter, r *http.Request) {
	// set the request header Content-Type for json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.getStatus())
}
