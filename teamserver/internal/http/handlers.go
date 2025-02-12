package http

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"

	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	agent, err := server.AgentService.AddAgent()
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to add id to database: %v", err)
		return
	}

	publicKeyPem, err := tools.PubRsaToPem(&agent.PrivateKey.PublicKey)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to tranform the public key to PEM: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(api.RegisterResponse{
		AgentId:   agent.AgentId,
		PublicKey: publicKeyPem,
	})
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to serialize register response with: %v", err)
		return
	}
}

func (server *Server) GetTasks(w http.ResponseWriter, r *http.Request) {
	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Print(err)
		return
	}

	tasks, err := server.TaskQueueService.GetTasks(agentId)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to get tasks for agent: %d - %v", agentId, err)
		return
	}

	t1 := teamserver.Task{}
	t2 := api.TaskResponse(t1)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(api.GetTasksResponse{
		Tasks: api.TaskResponse(tasks),
	})
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to serialize register response with: %v", err)
		return
	}
}

func PushTask(w http.ResponseWriter, r *http.Request) {
	database, err := database.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to open database with: %v", err)
		return
	}

	var pushTasksReq api.PushTaskRequest
	err = json.NewDecoder(r.Body).Decode(&pushTasksReq)
	if err != nil {
		api.RequestErrorHandler(w, err)
		log.Printf("Bad request for task: %v", err)
		return
	}

	err = (*database).TaskQueuePush(pushTasksReq.Task)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to push to task queue with: %v", err)
		return
	}
}

func SaveTaskResults(w http.ResponseWriter, r *http.Request) {
	database, err := database.NewDatabase()
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to open database with: %v", err)
		return
	}

	var taskResults api.PostTaskResultsRequest
	err = json.NewDecoder(r.Body).Decode(&taskResults)
	if err != nil {
		api.RequestErrorHandler(w, err)
		log.Printf("Bad request for task: %v", err)
		return
	}

	agentId, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("%v", err)
		return
	}

	err = (*database).SaveTaskResults(agentId, taskResults)
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to save task results: %v", err)
		return
	}
}
