package http

import (
	"errors"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"

	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func ApiErrorHandler(err error, w http.ResponseWriter) {
	target := &teamserver.ClientError{}

	if errors.As(err, &target) {
		api.RequestErrorHandler(w, err)
	} else { // Default to internal
		api.InternalErrorHandler(w)
	}
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	agent, err := server.AgentService.AddAgent()
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to add id to database: %v", err)
		return
	}

	resp, err := AgentToRegisterResponse(agent)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("%v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to serialize register response with: %v", err)
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

	tasks, err := server.TaskQueueService.GetTasks(agentId)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to get tasks for agent: %d - %v", agentId, err)
		return
	}

	resp := TasksToGetTasksResponse(tasks)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to serialize register response with: %v", err)
		return
	}
}

func (server *Server) PushTask(w http.ResponseWriter, r *http.Request) {
	var pushTasksReq api.PushTaskRequest
	err := json.NewDecoder(r.Body).Decode(&pushTasksReq)
	if err != nil {
		api.RequestErrorHandler(w, err)
		log.Printf("Bad request for task: %v", err)
		return
	}

	err = server.TaskQueueService.TaskQueuePush(pushTasksReq.Task)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to push to task queue with: %v", err)
		return
	}
}

func (server *Server) SaveTaskResults(w http.ResponseWriter, r *http.Request) {
	var taskResults api.PostTaskResultsRequest
	err := json.NewDecoder(r.Body).Decode(&taskResults)
	if err != nil {
		api.RequestErrorHandler(w, err)
		log.Printf("Bad request: %v", err)
		return
	}

	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		log.Printf("%v", err)
		return
	}

	err = server.TaskResultsService.SaveTaskResults(agentId, PostTaskResultsRequestToTaskResultsIn(&taskResults))
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to save task results: %v", err)
		return
	}
}

func (server *Server) GetTaskResults(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		log.Printf("%v", err)
		return
	}

	taskResultsDomain, err := server.TaskResultsService.GetTaskResults(agentId)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to get task results: %v", err)
		return
	}

	resp := TaskResultsOutToGetTaskResultsResponse(taskResultsDomain)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(err, w)
		log.Printf("Failed to serialize response with: %v", err)
		return
	}
}
