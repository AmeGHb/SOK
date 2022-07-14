package httpServer

import (
	"net"
	"net/http"
)

type Server struct {
	server   *http.Server
	listener net.Listener
	adress   string
}

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
