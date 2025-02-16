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

type PushTaskRequest struct {
	Task string
}

type PostTaskResultsRequest struct {
	TaskResults []struct {
		TaskId uint64
		Output string
	}
}

type GetTaskResultRequest struct {
	TaskIds []uint64
}

type GetTaskResultsResponse struct {
	TaskResults []struct {
		TaskId uint64
		Task   string
		Output string
	}
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
