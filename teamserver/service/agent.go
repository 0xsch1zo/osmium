package service

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (as *AgentService) AddAgent() (*teamserver.Agent, error) {
	rsaPriv, err := tools.GenerateKey()
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return nil, err
	}

	agent, err := as.agentRepository.AddAgent(rsaPriv)
	if err != nil {
		ServiceServerErrHandler(err, agentServiceStr, as.eventLogService)
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(as.callbacks))
	for _, listener := range as.callbacks {
		if listener != nil {
			go func() {
				listener(*agent)
				wg.Done()
			}()
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

func (as *AgentService) AddOnAgentAddedCallback(callback func(teamserver.Agent)) teamserver.CallbackHandle {
	as.callbacks = append(as.callbacks, callback)
	return teamserver.CallbackHandle(len(as.callbacks) - 1)
}

func (as *AgentService) RemoveOnAgentAddedCallback(handle teamserver.CallbackHandle) {
	for i := range as.callbacks {
		if teamserver.CallbackHandle(i) == handle {
			as.callbacks[i] = nil
			break
		}
	}
}
