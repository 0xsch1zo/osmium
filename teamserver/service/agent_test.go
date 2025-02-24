package service_test

import (
	"math/rand/v2"
	"testing"
)

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

func TestAgentExists(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	exists, err := testedServices.agentService.AgentExists(agent.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Fatal("Agent reported as not existing")
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

// Covers both get and update
func TestGetAndUpdateAgentTaskProgress(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent()
	if err != nil {
		t.Fatal(err)
	}

	tasksAssignedCount := rand.Uint32() % 100
	for i := 0; i < int(tasksAssignedCount); i++ {
		err = testedServices.taskQueueService.TaskQueuePush("some task")
		if err != nil {
			t.Fatal(err)
		}
	}

	err = testedServices.agentService.UpdateAgentTaskProgress(agent.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	taskProgress, err := testedServices.agentService.GetAgentTaskProgress(agent.AgentId)
	if err != nil {
		t.Fatal(err)
	}

	if taskProgress != uint64(tasksAssignedCount) {
		t.Fatal("Wrong task progress")
	}
}
