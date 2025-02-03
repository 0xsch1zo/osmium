package api

import (
	"encoding/json"
	"net/http"
)

type RegisterResponse struct {
	AgentId uint64
}

type Error struct {
	Code    int
	Message string
}

func writeError(w http.ResponseWriter, error Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.Code)

	json.NewEncoder(w).Encode(error)
}

func RequestErrorHandler(w http.ResponseWriter, err error) {
	writeError(w, Error{Code: http.StatusBadRequest, Message: err.Error()})
}

func InternalErrorHandler(w http.ResponseWriter) {
	writeError(w, Error{Code: http.StatusInternalServerError, Message: "Internal server error occured!"})
}
