package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	ErrorInternalServer = "internal server error"
	ErrorBadRequest     = "bad request"
	ErrorInvalidSchema  = "invalid schema"
)

var validate *validator.Validate

type server struct {
	service Service	
	register *prometheus.Registry
}

func NewServer(service Service, register *prometheus.Registry) *server {
	return &server{
		service: service,
		register : register,
	}
}

func (s *server) AddHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var logReq Log
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		http.Error(w, ErrorBadRequest, http.StatusBadRequest)
		return
	}

	if err := validate.Struct(logReq); err != nil {
		http.Error(w, ErrorInvalidSchema, http.StatusBadRequest)
		return
	}

	if err := s.service.Add(ctx, logReq); err != nil {
		http.Error(w, ErrorInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := s.service.GetAll(ctx)
	if err != nil {
		http.Error(w, ErrorInternalServer, http.StatusInternalServerError)
		return
	}
	logs := Logs{
		Logs: res,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		http.Error(w, ErrorInternalServer, http.StatusInternalServerError)
		return
	}
}

func (s *server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /add", s.AddHandler)
	mux.HandleFunc("GET /", s.GetAllHandler)
	mux.Handle("GET /metrics", promhttp.HandlerFor(s.register, promhttp.HandlerOpts{}))
	fmt.Println("Server is listening on port 3000")

	return http.ListenAndServe(":3000", mux)

}

func init() {
	validate = validator.New()
}
