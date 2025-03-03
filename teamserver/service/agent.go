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
		return nil, repositoryErrWrapper(err)
	}

	return agent, nil
}

func (as *AgentService) GetAgent(agentId uint64) (*teamserver.Agent, error) {
	agent, err := as.agentRepository.GetAgent(agentId)
	return agent, repositoryErrWrapper(err)
}

func (as *AgentService) AgentExists(agentId uint64) (bool, error) {
	exists, err := as.agentRepository.AgentExists(agentId)
	return exists, repositoryErrWrapper(err)
}

func (as *AgentService) GetAgentTaskProgress(agentId uint64) (uint64, error) {
	taskProgress, err := as.agentRepository.GetAgentTaskProgress(agentId)
	return taskProgress, repositoryErrWrapper(err)
}

func (as *AgentService) UpdateAgentTaskProgress(agentId uint64) error {
	return repositoryErrWrapper(as.agentRepository.UpdateAgentTaskProgress(agentId))
}

func (as *AgentService) ListAgents() ([]teamserver.AgentView, error) {
	agentViews, err := as.agentRepository.ListAgents()
	if err != nil {
		return nil, repositoryErrWrapper(err)
	}
	return agentViews, nil
}
