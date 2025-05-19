package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (server *Server) SaveTaskResult(w http.ResponseWriter, r *http.Request) {
	var taskResults api.PostTaskResultRequest
	err := json.NewDecoder(r.Body).Decode(&taskResults)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	taskId, err := strconv.ParseUint(r.PathValue("taskId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	domainTaskResults := PostTaskResultRequestToTaskResultsIn(&taskResults, taskId)
	err = server.TaskResultsService.SaveTaskResult(agentId, domainTaskResults)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}

func (server *Server) GetTaskResult(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	taskId, err := strconv.ParseUint(r.PathValue("taskId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	taskResultsDomain, err := server.TaskResultsService.GetTaskResult(agentId, taskId)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	resp := TaskResultsOutToGetTaskResultsResponse(taskResultsDomain)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}

func (server *Server) GetTaskResults(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	taskResults, err := server.TaskResultsService.GetTaskResults(agentId)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}

	err = templates.TaskResults(agentId, taskResults).Render(context.Background(), w)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}

func (server *Server) TaskResultsListen(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	server.TaskResultsService.AddOnTaskResultSavedCallback(func(resultAgentId uint64, taskResult teamserver.TaskResultIn) {
		if agentId == resultAgentId {
			taskResultOut, err := server.TaskResultsService.GetTaskResult(agentId, taskResult.TaskId)
			if err != nil {
				log.Print(err)
				return
			}

			buf := bytes.Buffer{}
			err = templates.TaskResultOOB(*taskResultOut).Render(context.Background(), &buf)
			if err != nil {
				log.Print(err)
				return
			}

			err = sendSSE(w, "task-result", buf.String())
			if err != nil {
				log.Print(err)
				return
			}
		}
	})

	wg.Wait()
}
