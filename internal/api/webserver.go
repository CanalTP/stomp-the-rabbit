package api

import (
	"net/http"
)

//Router : return ServerMux for circumstantial
func Router(handler *StatusHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", handler.status)
	return mux
}
