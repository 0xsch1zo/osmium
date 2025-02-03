package handlers

import (
	_ "github.com/mattn/go-sqlite3"

	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/handlers/tools"

	"encoding/json"
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	database, err := tools.NewDatabase()
	if err != nil {
		log.Printf("Failed to open database with: %v", err)
		return
	}

	AgentId, err := (*database).AddAgent()
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to add id to database: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(api.RegisterResponse{AgentId: AgentId})
	if err != nil {
		api.InternalErrorHandler(w)
		log.Printf("Failed to serialize uuid with: %v", err)
		return
	}
}

/*
func (apiHandlers *Api) ServeTaskQueue(w http.ResponseWriter, r *http.Request) {
}
*/
