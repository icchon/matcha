package helper

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/icchon/matcha/api/internal/apperrors"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	log.Printf("Responding with error: status=%d, message=%s", code, message)
	RespondWithJSON(w, code, ErrorResponse{Message: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	if payload == nil {
		payload = map[string]interface{}{}
	}
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Could not encode response payload: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "failed to encode response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func HandleError(w http.ResponseWriter, err error) {
	log.Printf("API Error: %v", err)

	if errors.Is(err, apperrors.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
		RespondWithError(w, http.StatusNotFound, "The requested resource was not found.")
		return
	}
	if errors.Is(err, apperrors.ErrInvalidInput) {
		RespondWithError(w, http.StatusBadRequest, "Invalid input provided.")
		return
	}
	if errors.Is(err, apperrors.ErrUnauthorized) {
		RespondWithError(w, http.StatusUnauthorized, "Authentication failed.")
		return
	}
	if errors.Is(err, apperrors.ErrUnhandled) {
		RespondWithError(w, http.StatusInternalServerError, "An unhandled error occurred.")
		return
	}
	if errors.Is(err, apperrors.ErrInternalServer) {
		RespondWithError(w, http.StatusInternalServerError, "Internal server error.")
		return
	}

	RespondWithError(w, http.StatusInternalServerError, "An unexpected internal error occurred.")
}
