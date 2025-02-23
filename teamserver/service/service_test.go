package service_test

import (
	"testing"

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

func TestGetAgent(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	agentReturned, err := testedServices.agentService.GetAgent(agent.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	if agent.PrivateKey.Equal(*agentReturned.PrivateKey) {
		t.Fatal("Keys don't match")
	}
}

func TestListAgents(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	agentList, err := testedServices.agentService.ListAgents()
	if err != nil {
		t.Fatal(err)
	}

	for _, agentListed := range agentList {
		if agentListed.AgentId == agent.AgentId {
			return
		}
	}
	t.Fatal("Agent not found when listed")
}
