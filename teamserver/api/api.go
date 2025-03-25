package api

import (
	"encoding/json"
	"net/http"
)

type RegisterResponse struct {
	AgentId   uint64
	PublicKey string
}

type GetTasksResponse struct {
	Tasks []struct {
		TaskId uint64
		Task   string
	}
}

type AddTaskRequest struct {
	Task string
}

type PostTaskResultRequest struct {
	Output string
}

type GetTaskResultRequest struct {
	TaskId uint64
}

type GetTaskResultResponse struct {
	TaskId uint64
	Task   string
	Output string
}

type ListAgentsResponse struct {
	AgentViews []struct {
		AgentId uint64
	}
}

type LoginResponse struct {
	ExpiryTime int64
}

type Error struct {
	Code    int
	Message string
}

func writeError(w http.ResponseWriter, error Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.Code)

	_ = json.NewEncoder(w).Encode(error)
}

func RequestErrorHandler(w http.ResponseWriter, err error) {
	writeError(w, Error{Code: http.StatusBadRequest, Message: err.Error()})
}

func InternalErrorHandler(w http.ResponseWriter) {
	writeError(w, Error{Code: http.StatusInternalServerError, Message: "Internal server error occured!"})
}
