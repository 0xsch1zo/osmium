package service_test

import (
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type testedServices struct {
	agentService       *service.AgentService
	tasksService       *service.TasksService
	taskResultsService *service.TaskResultsService
}

func newTestedServices() (*testedServices, error) {
	database, err := database.NewDatabase(":memory:")
	if err != nil {
		return nil, err
	}

	agentRepo := (*database).NewAgentRepository()
	taskQueueRepo := (*database).NewTasksRepository()
	taskResultsRepo := (*database).NewTaskResultsRepository()

	return &testedServices{
		agentService:       service.NewAgentService(*agentRepo),
		tasksService:       service.NewTasksService(*taskQueueRepo, *agentRepo),
		taskResultsService: service.NewTaskResultsService(*taskResultsRepo, *agentRepo, *taskQueueRepo),
	}, nil
}
