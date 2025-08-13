package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "gobank/db/sqlc"
	"net/http"
	"strconv"
	"strings"
)

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required,oneof=USD VND EUR"`
}

func (server *Server) createAccount(w http.ResponseWriter, r *http.Request) {
	// Get the JSON data
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate data
	if err := server.validate.Struct(req); err != nil {
		server.WriteError(w, http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Create new account into database
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0, // Default balance when creating account
		Currency: req.Currency,
	}
	account, err := server.store.CreateAccount(r.Context(), arg)
	if err != nil {
		server.logger.Error("POST /account: failed to create new account", "error", err)
		server.WriteError(w, http.StatusInternalServerError, "Failed to create new account")
		return
	}

	// Return the newly created account back to client
	server.WriteJSON(w, http.StatusCreated, account)
}

func (server *Server) getAccount(w http.ResponseWriter, r *http.Request) {
	// Get the id parameter
	idRaw := strings.TrimSpace(r.PathValue("id"))

	// Try parse ID
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil || id <= 0 {
		server.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid ID: %s", idRaw))
		return
	}

	// Get account by id
	account, err := server.store.GetAccount(r.Context(), id)
	if err != nil {
		// If ID not match any record in database
		if err == sql.ErrNoRows {
			server.logger.Warn("GET /account/{id}: account not found", "account_id", id)
			server.WriteError(w, http.StatusNotFound, "account not found")
			return
		}

		// Other database errors
		server.logger.Error("GET /account/{id}: failed to get account", "account_id", id)
		server.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to get account with ID: %d", id))
		return
	}

	// Return the fetched account to client
	server.WriteJSON(w, http.StatusFound, account)
}

func (server *Server) listAccounts(w http.ResponseWriter, r *http.Request) {
	// Get the page and offset
	params := r.URL.Query()
	pageIdRaw, pageSizeRaw := params.Get("page_id"), params.Get("page_size")

	// Try parse request parameter
	pageId, err := strconv.ParseInt(pageIdRaw, 10, 32)
	if err != nil || pageId <= 0 {
		server.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request parameter page_id: %s", pageIdRaw))
		return
	}

	pageSize, err := strconv.ParseInt(pageSizeRaw, 10, 32)
	if err != nil || pageSize <= 0 {
		server.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid request parameter page_size: %s", pageSizeRaw))
		return
	}

	// Get list of accounts
	accounts, err := server.store.ListAccount(r.Context(), db.ListAccountParams{
		Limit:  int32(pageSize),
		Offset: int32((pageId - 1) * pageSize),
	})
	if err != nil {
		server.logger.Error("GET /accounts: failed to get list of accounts", "error", err)
		server.WriteError(w, http.StatusInternalServerError, "failed to get list of accounts")
		return
	}

	server.WriteJSON(w, http.StatusFound, accounts)
}
