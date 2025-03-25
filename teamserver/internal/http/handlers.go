package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"

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

	_, err = server.TasksService.AddTask(agentId, addTaskReq.Task)
	if err != nil {
		ApiErrorHandler(fmt.Errorf("Failed to push to task queue with: %w", err), w)
		return
	}
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
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err)
		return
	}

	// Don't care about tasks before
	tasksExclude, err := server.TasksService.GetTasks(agentId)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	for {
		// Poll untill we get a newly started task
		var tasks []teamserver.Task
		for len(tasks) == 0 {
			tasks, err = server.TasksService.GetTasks(agentId)
			if err != nil {
				ApiErrorHandler(err, w)
				return
			}

			tasks = slices.DeleteFunc(tasks, func(task teamserver.Task) bool {
				return slices.Contains(tasksExclude, task)
			})
		}

		// Unlikely to have many tasks
		for _, task := range tasks {
			var exists bool
			// Wait for taskResult to be devlivered
			for !exists {
				exists, err = server.TaskResultsService.TaskResultExists(agentId, task.TaskId)
				if err != nil {
					ApiErrorHandler(err, w)
					return
				}
			}

			taskResult, err := server.TaskResultsService.GetTaskResult(agentId, task.TaskId)
			if err != nil {
				ApiErrorHandler(err, w)
				return
			}

			err = sendEventMessage(w, taskResult.Output)
			if err != nil {
				api.InternalErrorHandler(w)
				return
			}
			exists = false
		}
	}
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	username := r.Form["username"]
	password := r.Form["password"]
	if len(username) == 0 || len(password) == 0 {
		api.RequestErrorHandler(w, errors.New(errUnauthorized))
		return
	}

	token, err := server.AuthorizationService.Login(username[0], password[0])
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	w.Header().Add("HX-Redirect", "/")
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  token.ExpiryTime,
		Value:    token.Token,
	})
}

func (server *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err == http.ErrNoCookie {
		api.RequestErrorHandler(w, errors.New(errUnauthorized))
		return
	}

	refreshedToken, err := server.AuthorizationService.RefreshToken(token.Value)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshedToken.ExpiryTime,
		Value:    refreshedToken.Token,
	})
}
