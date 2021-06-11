package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/CanalTP/stomptherabbit/internal/scoreboard"
	"github.com/CanalTP/stomptherabbit/internal/webstomp"
	"github.com/streadway/amqp"
)

type responseStatus struct {
	ApplicationName                 string    `json:"applicationName"`
	ApplicationVersion              string    `json:"applicationVersion"`
	Hostname                        string    `json:"hostname"`
	AmqpHostname                    string    `json:"amqpHostname"`
	WebstompHostname                string    `json:"webstompHostname"`
	GoRuntimeVersion                string    `json:"goRuntimeVersion"`
	IsAmqpConnActive                bool      `json:"isAmqpConnActive"`
	IsWebStompConnActive            bool      `json:"isWebStompConnActive"`
	DateLastMessageReceiveSNCF      time.Time `json:"dateLastMessageReceiveSNCF"`
	DateLastAttempSendRabbitMQ      time.Time `json:"dateLastMessageSendRabbitMQ"`
	DateLastSuccessfulWriteRabbitMQ time.Time `json:"dateLastMessageSuccessfulWriteRabbitMQ"`
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
func (r *responseStatus) setStatusStomp(stompOpts webstomp.Options) {
	r.IsWebStompConnActive = false
	r.WebstompHostname = stompOpts.Target
	if websocketconn, err := webstomp.Dial(stompOpts.Target, stompOpts.Protocol); err == nil {
		r.IsWebStompConnActive = true
		websocketconn.Close()
	}
}

func (r *responseStatus) setStatusScoreBoard(mapScore map[string]time.Time) {
	r.DateLastMessageReceiveSNCF = mapScore["dateLastMessageReceiveSNCF"]
	r.DateLastAttempSendRabbitMQ = mapScore["dateLastMessageAttempSendRabbitMQ"]
	r.DateLastSuccessfulWriteRabbitMQ = mapScore["dateLastMessageSuccessfulWriteRabbitMQ"]
}

// StatusHandler : web handler for stomptherabbit status
type StatusHandler struct {
	amqpURL            string
	webstompOptions    webstomp.Options
	applicationVersion string
	applicationName    string
	scoreboard         *scoreboard.ScoreBoard
}

//NewStatusHandler : constuctor for statusHandler
func NewStatusHandler(amqpURL string, webstompOptions webstomp.Options, version, applicationName string, scoreboard *scoreboard.ScoreBoard) *StatusHandler {
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
