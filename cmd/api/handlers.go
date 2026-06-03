package main

import (
	"encoding/json"
	"go-core/internal/data"
	"net/http"
)

type JsonResponse map[string]interface{}

func writeJsonResponse(w http.ResponseWriter, statusCode int, data JsonResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Application struct holds dependencies for our route handlers
type application struct {
	store *data.Store
}

func (app *application) handleListTeams(w http.ResponseWriter, r *http.Request) {
	entities := data.ListEntities(app.store)
	writeJsonResponse(w, http.StatusOK, JsonResponse{"entities": entities})
}

type tipPayload struct {
	Amount float64 `json:"amount"`
}

// HandleProcessTip simulates an incoming mobile wallet deposit to a entity
func (app *application) handleProcessTip(w http.ResponseWriter, r *http.Request) {
	// Parse the entity ID directly from the URL path placeholder
	entityID := r.PathValue("id")

	var payload tipPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, JsonResponse{"error": "invalid JSON payload"})
		return
	}

	transaction, err := data.ProcessTip(app.store, entityID, payload.Amount)
	if err != nil {
		switch err {
		case data.ErrEntityNotFound:
			writeJsonResponse(w, http.StatusNotFound, JsonResponse{"error": "entity not found"})
		case data.ErrInvalidAmount:
			writeJsonResponse(w, http.StatusBadRequest, JsonResponse{"error": "amount must be greater than zero"})
		default:
			writeJsonResponse(w, http.StatusInternalServerError, JsonResponse{"error": "internal server error"})
		}
		return
	}
	writeJsonResponse(w, http.StatusOK, JsonResponse{"transaction": transaction})
}

type withdrawPayload struct {
	Amount float64 `json:"amount"`
}

// HandleProcessWithdrawal simulates a withdrawal request from a entity to their mobile wallet
func (app *application) handleProcessWithdrawal(w http.ResponseWriter, r *http.Request) {
	// Parse the entity ID directly from the URL path placeholder
	entityID := r.PathValue("id")
	var payload withdrawPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		writeJsonResponse(w, http.StatusBadRequest, JsonResponse{"error": "invalid JSON payload"})
		return
	}
	transaction, err := data.ProcessWithdrawal(app.store, entityID, payload.Amount)
	if err != nil {
		switch err {
		case data.ErrEntityNotFound:
			writeJsonResponse(w, http.StatusNotFound, JsonResponse{"error": "entity not found"})
		case data.ErrBelowThreshold:
			writeJsonResponse(w, http.StatusBadRequest, JsonResponse{"error": "withdrawal amount is below the minimum threshold"})
		case data.ErrInsufficientFunds:
			writeJsonResponse(w, http.StatusBadRequest, JsonResponse{"error": "insufficient funds in entity wallet"})
		default:
			writeJsonResponse(w, http.StatusInternalServerError, JsonResponse{"error": "internal server error"})
		}
		return
	}
	writeJsonResponse(w, http.StatusOK, JsonResponse{"transaction": transaction})
}
