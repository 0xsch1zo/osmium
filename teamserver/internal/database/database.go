package database

import (
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database/sqlite"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Database interface {
	NewAgentRepository() *service.AgentRepository
	NewTaskQueueRepository() *service.TaskQueueRepository
	NewTaskResultsRepository() *service.TaskResultsRepository
}

func NewDatabase() (*Database, error) {
	var db Database
	db, err := sqlite.SetupDatabase()
	return &db, err
}
