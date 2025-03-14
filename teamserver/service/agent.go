package service

import (
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (as *AgentService) AddAgent() (*teamserver.Agent, error) {
	rsaPriv, err := tools.GenerateKey()
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	agent, err := as.agentRepository.AddAgent(rsaPriv)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

func (as *AgentService) GetAgent(agentId uint64) (*teamserver.Agent, error) {
	err := as.AgentExists(agentId)
	if err != nil {
		return nil, err
	}

	agent, err := as.agentRepository.GetAgent(agentId)
	return agent, err
}

func (as *AgentService) AgentExists(agentId uint64) error {
	exists, err := as.agentRepository.AgentExists(agentId)
	if err != nil {
		return err
	}

	if !exists {
		return teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId))
	}
	return nil
}

func (as *AgentService) GetAgentTaskProgress(agentId uint64) (uint64, error) {
	err := as.AgentExists(agentId)
	if err != nil {
		return 0, err
	}

	taskProgress, err := as.agentRepository.GetAgentTaskProgress(agentId)
	return taskProgress, err
}

func (as *AgentService) UpdateAgentTaskProgress(agentId uint64) error {
	err := as.AgentExists(agentId)
	if err != nil {
		return err
	}

	return as.agentRepository.UpdateAgentTaskProgress(agentId)
}

func (as *AgentService) ListAgents() ([]teamserver.AgentView, error) {
	agentViews, err := as.agentRepository.ListAgents()
	if err != nil {
		return nil, err
	}
	return agentViews, nil
}
