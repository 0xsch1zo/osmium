package database

import (
	"database/sql"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/teamserver"
)

/*type Database interface {
	SetupDatabase() error
	AddAgent() (*osmium.Agent, error)
	GetAgent(agentId uint64) (*osmium.Agent, error)
	GetTasks(agentId uint64) ([]osmium.Task, error)
	UpdateAgentTaskProgress(agentId uint64) error
	TaskQueuePush(task string) error // TODO: Use api structs as input
	SaveTaskResults(agentId uint64, taskResults api.PostTaskResultsRequest) error
	GetTaskResult(agentId uint64, taskId uint64) (*osmium.TaskResult, error)
	GetTaskResults(agentId uint64) ([]osmium.TaskResult, error)
	//TaskQueuePop() error
}*/

type Database interface {
	NewAgentService(*sql.DB) (*teamserver.AgentService, error)
	NewTaskQueueService(*sql.DB) (*teamserver.TaskQueueService, error)
	NewTaskResultService(*sql.DB) (*teamserver.TaskResultsService, error)
}

/*func NewDatabase() (*Database, error) {
	var db Database = &SqliteDb{}

	err := db.NewAgentService()
	if err != nil {
		return nil, err
	}

	return &db, nil
}*/
