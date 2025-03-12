package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

const (
	errSerializationFmt = "Failed to serialize register response with: %w"
	errUnauthorized     = "Unauthorized"
)

func sendEventMessage(w http.ResponseWriter, message string) error {
	_, err := w.Write([]byte(
		`event: message
data: ` + message + "\n"))
	return err
}

func ApiErrorHandler(err error, w http.ResponseWriter) {
	target := &teamserver.ClientError{}

	if errors.As(err, &target) {
		api.RequestErrorHandler(w, err)
	} else { // Default to internal
		api.InternalErrorHandler(w)
	}

	log.Print(err)
}

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

	taskId, err := server.TasksService.AddTask(agentId, addTaskReq.Task)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to push to task queue with: %w", err), w)
		return
	}

	server.awaitedTaskIdChannel <- taskId
}

func (server *Server) SaveTaskResult(w http.ResponseWriter, r *http.Request) {
	var taskResults api.PostTaskResultRequest
	err := json.NewDecoder(r.Body).Decode(&taskResults)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	taskId, err := strconv.ParseUint(r.PathValue("taskId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	domainTaskResults := PostTaskResultRequestToTaskResultsIn(&taskResults, taskId)
	err = server.TaskResultsService.SaveTaskResult(agentId, domainTaskResults)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to save task results: %w", err), w)
		return
	}

	server.taskResultsChannel <- domainTaskResults
}

func (server *Server) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	taskId, err := strconv.ParseUint(r.PathValue("taskId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	taskResultsDomain, err := server.TaskResultsService.GetTaskResult(agentId, taskId)
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

func (server *Server) ListenAndServeTaskResults(w http.ResponseWriter, r *http.Request) {
	for {
		taskId, ok := <-server.awaitedTaskIdChannel
		if !ok {
			break
		}

		taskResults, ok := <-server.taskResultsChannel
		if !ok {
			break
		}

		if taskResults.TaskId == taskId {
			err := sendEventMessage(w, taskResults.Output)
			if err != nil {
				api.InternalErrorHandler(w)
			}
			break
		}
	}
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	var creds api.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		api.InternalErrorHandler(w)
		return
	}

	if len(creds.Username) == 0 || len(creds.Password) == 0 {
		api.RequestErrorHandler(w, errors.New(errUnauthorized))
	}

	sessionToken, err := server.AuthorizationService.Login(creds.Username, creds.Password)
	if err != nil {
		ApiErrorHandler(err, w)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(tokenExpiryTime),
		HttpOnly: true,
	})
}
