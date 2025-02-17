package database

import (
	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database/sqlite"
)

type Database interface {
	teamserver.AgentService
	teamserver.TaskQueueService
	teamserver.TaskResultsService
}

func NewDatabase() (*Database, error) {
	var db Database
	db, err := sqlite.SetupDatabase()
	return &db, err
}
