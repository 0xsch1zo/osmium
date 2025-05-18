package service_test

import (
	"testing"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func TestGetAgent(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	registerInfo := teamserver.AgentRegisterInfo{
		Username: "some",
		Hostname: "host",
	}
	agent, err := testedServices.agentService.AddAgent(registerInfo)
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

	if agent.AgentInfo.Hostname != registerInfo.Hostname {
		t.Fatal("Hostnames don't match")
	}

	if agent.AgentInfo.Username != registerInfo.Username {
		t.Fatal("Usernames don't match")
	}
}

func TestAgentExists(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	agent, err := testedServices.agentService.AddAgent(teamserver.AgentRegisterInfo{})
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

	agent, err := testedServices.agentService.AddAgent(teamserver.AgentRegisterInfo{})
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
