package service

import (
	"fmt"
	"net/http"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (as *AgentService) AddAgent(agentInfo teamserver.AgentRegisterInfo) (*teamserver.Agent, error) {
	rsaPriv, err := tools.GenerateKey()
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return nil, err
	}

	agent, err := as.agentRepository.AddAgent(rsaPriv, agentInfo)
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return nil, err
	}

	for _, listener := range as.onAgentAddedCallbacks {
		if listener != nil {
			go listener(*agent)
		}
	}

	as.eventLogService.LogEvent(
		teamserver.Info,
		fmt.Sprintf("Agent registered\n>> AgentId: %d", agent.AgentId),
	)

	return agent, nil
}

func (as *AgentService) GetAgent(agentId uint64) (*teamserver.Agent, error) {
	err := as.AgentExists(agentId)
	if err != nil {
		return nil, err
	}

	agent, err := as.agentRepository.GetAgent(agentId)
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return nil, err
	}
	return agent, nil
}

func (as *AgentService) AgentExists(agentId uint64) error {
	exists, err := as.agentRepository.AgentExists(agentId)
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return err
	}

	if !exists {
		return teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId), http.StatusNotFound)
	}
	return nil
}

func (as *AgentService) ListAgents() ([]teamserver.AgentView, error) {
	agentViews, err := as.agentRepository.ListAgents()
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return nil, err
	}
	return agentViews, nil
}

func (as *AgentService) UpdateLastCallbackTime(agentId uint64) error {
	err := as.agentRepository.UpdateLastCallbackTime(agentId)
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return err
	}

	agent, err := as.GetAgent(agentId)
	if err != nil {
		return err
	}

	for _, callback := range as.onCallbackTimeUpdatedCallbacks {
		if callback != nil {
			go callback(*agent)
		}
	}

	return nil
}

func (as *AgentService) AddOnAgentAddedCallback(callback func(teamserver.Agent)) teamserver.CallbackHandle {
	as.onAgentAddedCallbacks = append(as.onAgentAddedCallbacks, callback)
	return teamserver.CallbackHandle(len(as.onAgentAddedCallbacks) - 1)
}

func (as *AgentService) RemoveOnAgentAddedCallback(handle teamserver.CallbackHandle) {
	for i := range as.onAgentAddedCallbacks {
		if teamserver.CallbackHandle(i) == handle {
			as.onAgentAddedCallbacks[i] = nil
			break
		}
	}
}

func (as *AgentService) AddOnCallbackTimeUpdatedCallback(callback func(teamserver.Agent)) teamserver.CallbackHandle {
	as.onCallbackTimeUpdatedCallbacks = append(as.onCallbackTimeUpdatedCallbacks, callback)
	return teamserver.CallbackHandle(len(as.onCallbackTimeUpdatedCallbacks) - 1)
}

func (as *AgentService) RemoveOnCallbackTimeUpdatedCallback(handle teamserver.CallbackHandle) {
	for i := range as.onCallbackTimeUpdatedCallbacks {
		if teamserver.CallbackHandle(i) == handle {
			as.onCallbackTimeUpdatedCallbacks[i] = nil
			break
		}
	}
}
