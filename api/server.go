package api

import (
	"encoding/json"
	"fmt"
	db "gobank/db/sqlc"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	store    db.Store
	mux      *http.ServeMux
	logger   *slog.Logger
	validate *validator.Validate
}

func NewServer(store db.Store, logger *slog.Logger) *Server {
	server := &Server{
		store:    store,
		mux:      http.NewServeMux(),
		logger:   logger,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	server.RegisterHandler()

	return server
}

func (server *Server) RegisterHandler() {
	// Account route
	server.mux.HandleFunc("POST /account", server.createAccount)
	server.mux.HandleFunc("GET /account/{id}", server.getAccount)
	server.mux.HandleFunc("GET /accounts", server.listAccounts)
}

func (server *Server) Start(domain, port string) error {
	server.logger.Info(fmt.Sprintf("Server start at %s:%s", domain, port))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), server.mux)
}

func (server *Server) WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (server *Server) WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"data": data,
	})
}
