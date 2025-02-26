package http

import (
	"errors"
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"

	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const errSerializationFmt = "Failed to serialize register response with: %w"

func ApiErrorHandler(err error, w http.ResponseWriter) {
	target := &teamserver.ClientError{}

	if errors.As(err, &target) {
		api.RequestErrorHandler(w, err)
	} else { // Default to internal
		api.InternalErrorHandler(w)
	}

	log.Print(err)
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	agent, err := server.AgentService.AddAgent()
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to add agent: %w", err), w)
		return
	}

	resp, err := AgentToRegisterResponse(agent)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(fmt.Errorf(errSerializationFmt, err), w)
		return
	}
}

func (server *Server) GetTasks(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err) // Clients error
		log.Print(err)
		return
	}

	tasks, err := server.AgentService.GetTasks(agentId)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to get tasks for agent: %d - %w", agentId, err), w)
		return
	}

	resp := TasksToGetTasksResponse(tasks)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(fmt.Errorf(errSerializationFmt, err), w)
		return
	}
}

func (server *Server) PushTask(w http.ResponseWriter, r *http.Request) {
	var pushTasksReq api.PushTaskRequest
	err := json.NewDecoder(r.Body).Decode(&pushTasksReq)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	_, err = server.TaskQueueService.TaskQueuePush(pushTasksReq.Task)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to push to task queue with: %w", err), w)
		return
	}
}

func (server *Server) SaveTaskResults(w http.ResponseWriter, r *http.Request) {
	var taskResults api.PostTaskResultsRequest
	err := json.NewDecoder(r.Body).Decode(&taskResults)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	err = server.TaskResultsService.SaveTaskResults(agentId, PostTaskResultsRequestToTaskResultsIn(&taskResults))
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to save task results: %w", err), w)
		return
	}
}

func (server *Server) GetTaskResults(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	taskResultsDomain, err := server.TaskResultsService.GetTaskResults(agentId)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to get task results: %w", err), w)
		return
	}

	resp := TaskResultsOutToGetTaskResultsResponse(taskResultsDomain)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(fmt.Errorf(errSerializationFmt, err), w)
		return
	}
}
