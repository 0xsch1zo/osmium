package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (server *Server) AgentRegister(w http.ResponseWriter, r *http.Request) {
	agent, err := server.AgentService.AddAgent()
	if err != nil {
		ApiErrorHandler(err, w)
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
		api.RequestErrorHandler(w, err, http.StatusBadRequest) // Clients error
		return
	}

	tasks, err := server.TasksService.GetTasks(agentId)
	if err != nil {
		ApiErrorHandler(err, w)
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
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	var addTaskReq api.AddTaskRequest
	err = json.NewDecoder(r.Body).Decode(&addTaskReq)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	_, err = server.TasksService.AddTask(agentId, addTaskReq.Task)
	if err != nil {
		ApiErrorHandler(err, w)
		return
	}
}

func (server *Server) AgentSocket(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("agentId"), 10, 64)
	if err != nil {
		api.RequestErrorHandler(w, err, http.StatusBadRequest)
		return
	}

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		SocketErrorHandler(err, conn)
		return
	}

	for {
		_, messageReader, err := conn.NextReader()
		if err != nil {
			SocketErrorHandler(err, conn)
			return
		}

		var task api.AddTaskRequest
		err = json.NewDecoder(messageReader).Decode(&task)
		if err != nil {
			SocketErrorHandler(err, conn)
			return
		}

		taskId, err := server.TasksService.AddTask(agentId, task.Task)
		if err != nil {
			SocketErrorHandler(err, conn)
			return
		}

		// Waiting for task result to be added
		wg := sync.WaitGroup{}
		wg.Add(1)
		_ = server.TaskResultsService.AddOnTaskResultSavedCallback(func(agentSaved uint64, result teamserver.TaskResultIn) {
			if agentId == agentSaved && taskId == result.TaskId {
				wg.Done()
			}
		})
		wg.Wait()

		taskResult, err := server.TaskResultsService.GetTaskResult(agentId, taskId)
		if err != nil {
			SocketErrorHandler(err, conn)
			return
		}

		messageWriter, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			SocketErrorHandler(err, conn)
			return
		}

		err = json.NewEncoder(messageWriter).Encode(taskResult)
		if err != nil {
			SocketErrorHandler(err, conn)
			return
		}

		messageWriter.Close()
	}
}

func (s *Server) AddAgentListen(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	s.AgentService.AddOnAgentAddedCallback(func(agent teamserver.Agent) {
		buf := bytes.Buffer{}
		err := templates.AgentOOB(&agent).Render(r.Context(), &buf)
		if err != nil {
			log.Print(err)
		}

		err = sendSSE(w, "agent", buf.String())
		if err != nil {
			log.Print(err)
		}
	})

	wg.Wait()
}
