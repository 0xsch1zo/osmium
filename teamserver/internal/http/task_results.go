package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

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
