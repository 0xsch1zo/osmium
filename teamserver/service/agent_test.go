package service_test

import (
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

	err = testedServices.agentService.AgentExists(agent.AgentId)
	if err != nil {
		t.Fatal(err)
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
