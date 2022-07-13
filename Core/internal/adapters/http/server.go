package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	Okafka "github.com/segmentio/kafka-go"

	"transaction/config"
	"transaction/internal/adapters/http/validator"
	"transaction/internal/adapters/kafka"
	"transaction/internal/domain/dErrors"
	innerUser "transaction/internal/users"
	"transaction/internal/users/db"
)

type Server struct {
	server     *http.Server
	l          net.Listener
	address    string
	repository *db.Repository
	context    context.Context
	kafka      *kafka.ClientW
}

type CurrentRequest struct {
	Id          string
	Transaction float64
	Sign        string
}

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	err error
	s   Server
)

func New(ctx context.Context, conf *config.Config, repository *db.Repository, kafkaClient *kafka.ClientW) (*Server, error) {

	httpPort := conf.GetHttpPort()

	s.l, err = net.Listen("tcp", ":"+httpPort)

	if err != nil {
		log.Fatal("Failed listen port", err)
	}
	log.Printf("API server listening at: %s", httpPort)

	s.address = httpPort
	s.repository = repository
	s.context = ctx
	s.kafka = kafkaClient
	s.server = &http.Server{
		Handler: s.routes(),
	}

	return &s, nil
}

func (s *Server) Start() error {
	if err := s.server.Serve(s.l); !errors.Is(err, http.ErrServerClosed) {
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

	r.Post("/create?_id={_id}&name={name}&balance={balance}&email={email}", s.createOne)
	r.Get("/info?_id={_id}", s.getInfo)
	r.Post("/transaction?_id={_id}&sign={sign}&value={value}", s.makeTransaction)

	r.NotFound(methodNotFound)

	return r
}

func get(w http.ResponseWriter, r *http.Request) {
	okResponse("Transaction server is running.").send(w)
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

func bindBody(payload interface{}, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		return dErrors.NewError(err.Error())
	}

	return nil
}

func (s *Server) createOne(w http.ResponseWriter, r *http.Request) {

	request := map[string]string{}
	json.NewDecoder(r.Body).Decode(&request)

	code, message := validator.ValidateUser(request)
	if code != 200 {
		errorResponse(code, message).send(w)
	}

	balance, _ := strconv.ParseFloat(request["balance"], 64)

	currentUser := innerUser.User{
		Id:      request["_id"],
		Name:    request["name"],
		Email:   request["email"],
		Balance: balance,
	}

	err := s.repository.Create(s.context, &currentUser)
	if err != nil {
		log.Fatalf(
			"Error was occurred while creating a new user. ID: %s, name: %s, balance: %f. The error is : %v",
			currentUser.Id,
			currentUser.Name,
			currentUser.Balance,
			err,
		)
	}

	users, err := s.repository.FindAll(s.context)
	if err != nil {
		log.Fatalln(err)
	}

	okResponse(users).send(w)
}

func (s *Server) getInfo(w http.ResponseWriter, r *http.Request) {

	request := map[string]string{}
	json.NewDecoder(r.Body).Decode(&request)

	code, message := validator.ValidateId(request)
	if code != 200 {
		errorResponse(code, message).send(w)
	}

	currentUser, err := s.repository.FindOne(s.context, request["_id"])
	if err != nil {
		log.Fatalf("Could not find a user with this Id. Error: %v", err)
	}

	okResponse(currentUser).send(w)
}

func (s *Server) makeTransaction(w http.ResponseWriter, r *http.Request) {

	request := map[string]string{}
	json.NewDecoder(r.Body).Decode(&request)

	code, message := validator.ValidateValues(request)

	if code != 200 {
		errorResponse(code, message).send(w)
	}

	transaction, _ := strconv.ParseFloat(request["value"], 64)
	sign := request["sign"]
	currentUser, err := s.repository.FindOne(s.context, request["_id"])
	if err != nil {
		log.Fatalf("Could not find a user with this Id. Error: %v", err)
	}

	err = s.repository.Update(s.context, currentUser, transaction, sign)
	if err != nil {

		kafkaMessages := []Okafka.Message{
			{
				Topic: "transaction",
				Key:   []byte("error"),
				Value: []byte(currentUser.Email + " error"),
			},
		}

		s.kafka.SendMessages(kafkaMessages)
		log.Fatalf("Could not update data for the user. Error: %v", err)
	}

	dataBytes, _ := json.Marshal("The transaction was successfully updated.")
	w.Write(dataBytes)

	kafkaMessages := []Okafka.Message{
		{
			Topic: "transaction",
			Key:   []byte("successes"),
			Value: []byte(currentUser.Email + " ok"),
		},
	}

	s.kafka.SendMessages(kafkaMessages)
}
