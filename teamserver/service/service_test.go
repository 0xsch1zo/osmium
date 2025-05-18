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
	db, err := database.NewDatabase(":memory:")
	if err != nil {
		return nil, err
	}

	agentRepo := (*db).NewAgentRepository()
	tasksRepo := (*db).NewTasksRepository()
	taskResultsRepo := (*db).NewTaskResultsRepository()
	authRepo := (*db).NewAuthorizationRepository()
	eventLogRepo := (*db).NewEventLogRepository()

	eventLogService := service.NewEventLogService(eventLogRepo)
	agentService := service.NewAgentService(agentRepo, eventLogService)
	tasksService := service.NewTasksService(tasksRepo, agentService, eventLogService)
	taskResultsService := service.NewTaskResultsService(taskResultsRepo, agentService, tasksService, eventLogService)
	authorizationService := service.NewAuthorizationService(authRepo, testingJwtKey, eventLogService)

	return &testedServices{
		agentService:         agentService,
		tasksService:         tasksService,
		taskResultsService:   taskResultsService,
		authorizationService: authorizationService,
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
