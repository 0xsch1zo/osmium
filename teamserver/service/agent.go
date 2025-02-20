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

func (as *AgentService) getAgent(agentId uint64) (*teamserver.Agent, error) {
	agent, err := as.agentRepository.GetAgent(agentId)
	return agent, repositoryErrWrapper(err)
}

func (as *AgentService) agentExists(agentId uint64) (bool, error) {
	exists, err := as.agentRepository.AgentExists(agentId)
	return exists, err
}

func (as *AgentService) getAgentTaskProgress(agentId uint64) (uint64, error) {
	taskProgress, err := as.agentRepository.GetAgentTaskProgress(agentId)
	return taskProgress, repositoryErrWrapper(err)
}

func (as *AgentService) updateAgentTaskProgress(agentId uint64) error {
	return repositoryErrWrapper(as.agentRepository.UpdateAgentTaskProgress(agentId))
}
