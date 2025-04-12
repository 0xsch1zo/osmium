package service_test

import (
	"testing"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

const testingJwtKey string = "TESTING"

type testedServices struct {
	agentService         *service.AgentService
	tasksService         *service.TasksService
	taskResultsService   *service.TaskResultsService
	authorizationService *service.AuthorizationService
	eventLogService      *service.EventLogService
}

func newTestedServices() (*testedServices, error) {
	database, err := database.NewDatabase(":memory:")
	if err != nil {
		return nil, err
	}

	agentRepo := (*database).NewAgentRepository()
	taskQueueRepo := (*database).NewTasksRepository()
	taskResultsRepo := (*database).NewTaskResultsRepository()
	authorizationRepo := (*database).NewAuthorizationRepository()
	eventLogService := service.NewEventLogService((*database).NewEventLogRepository())

	return &testedServices{
		agentService:         service.NewAgentService(agentRepo, eventLogService),
		tasksService:         service.NewTasksService(taskQueueRepo, agentRepo, eventLogService),
		taskResultsService:   service.NewTaskResultsService(taskResultsRepo, agentRepo, taskQueueRepo, eventLogService),
		authorizationService: service.NewAuthorizationService(authorizationRepo, testingJwtKey, eventLogService),
		eventLogService:      eventLogService,
	}, nil
}

func fatalErrUnexpectedData(t *testing.T, err string, expected, recieved any) {
	t.Error(err)
	t.Error("Expected:")
	t.Error(expected)
	t.Error("Got:")
	t.Fatal(recieved)
}
