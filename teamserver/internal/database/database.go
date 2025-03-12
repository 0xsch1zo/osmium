package database

import (
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database/sqlite"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Database interface {
	NewAgentRepository() *service.AgentRepository
	NewTasksRepository() *service.TasksRepository
	NewTaskResultsRepository() *service.TaskResultsRepository
	NewAuthorizationRepository() *service.AuthorizationRepository
}

func NewDatabase(sourceString string) (*Database, error) {
	var db Database
	db, err := sqlite.SetupDatabase(sourceString)
	return &db, err
}
