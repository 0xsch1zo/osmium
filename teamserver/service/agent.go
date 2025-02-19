package service

import (
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
		return nil, teamserver.NewServerError(err.Error())
	}

	return agent, nil
}

func (as *AgentService) getAgent(agentId uint64) (*teamserver.Agent, error) {
	return as.agentRepository.GetAgent(agentId)
}

func (as *AgentService) agentExists(agentId uint64) (bool, error) {
	return as.agentRepository.AgentExists(agentId)
}

func (as *AgentService) getAgentTaskProgress(agentId uint64) (uint64, error) {
	return as.agentRepository.GetAgentTaskProgress(agentId)
}

func (as *AgentService) updateAgentTaskProgress(agentId uint64) error {
	return as.agentRepository.UpdateAgentTaskProgress(agentId)
}
