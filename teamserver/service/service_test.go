package service_test

import (
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type testedServices struct {
	agentService       *service.AgentService
	taskQueueService   *service.TaskQueueService
	taskResultsService *service.TaskResultsService
}

func newTestedServices() (*testedServices, error) {
	database, err := database.NewDatabase(":memory:")
	if err != nil {
		return nil, err
	}

	agentRepo := (*database).NewAgentRepository()
	taskQueueRepo := (*database).NewTaskQueueRepository()
	taskResultsRepo := (*database).NewTaskResultsRepository()

	return &testedServices{
		agentService:       service.NewAgentService(*agentRepo),
		taskQueueService:   service.NewTaskQueueService(*taskQueueRepo, *agentRepo),
		taskResultsService: service.NewTaskResultsService(*taskResultsRepo, *agentRepo, *taskQueueRepo),
	}, nil
}
