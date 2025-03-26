package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

func (server *Server) AgentRegister(w http.ResponseWriter, r *http.Request) {
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
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err) // Clients error
		log.Print(err)
		return
	}

	tasks, err := server.TasksService.GetTasks(agentId)
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

func (server *Server) AddTask(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	var addTaskReq api.AddTaskRequest
	err = json.NewDecoder(r.Body).Decode(&addTaskReq)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	_, err = server.TasksService.AddTask(agentId, addTaskReq.Task)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to push to task queue with: %w", err), w)
		return
	}
}
