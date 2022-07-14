package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mail-sender/config"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func New(conf *config.Config, logger *log.Logger) (server *Server, err error) {

	port := conf.GetHttpPort()
	server.listener, err = net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatalf("Failed listen to the port. Error: %v", err)
	}

	logger.Printf("API server listening at: %s", port)

	server.adress = port
	server.server = &http.Server{
		Handler: server.routes(),
	}

	return server, nil

}

func (s *Server) Start() error {
	if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s Server) routes() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", get)

	return r
}

func okResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Status:  http.StatusOK,
		Message: "ok",
		Data:    data,
	}
}

func errorResponse(status int, err string) *APIResponse {
	return &APIResponse{
		Status:  status,
		Message: err,
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	okResponse("Get request for testing is Server running").send(w)
}

func methodNotFound(w http.ResponseWriter, r *http.Request) {
	errorResponse(http.StatusNotFound, "Not found").send(w)
}

func (r *APIResponse) send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	dataBytes, err := json.Marshal(r.Data)
	if err != nil {
		fmt.Println("some error:", err)
		return
	}
	w.Write(dataBytes)
}
