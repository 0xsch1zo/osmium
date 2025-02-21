package http

import (
	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func AgentToRegisterResponse(agent *teamserver.Agent) (*api.RegisterResponse, error) {
	pubPem, err := tools.PubRsaToPem(&agent.PrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &api.RegisterResponse{
		AgentId:   agent.AgentId,
		PublicKey: pubPem,
	}, nil
}

func TasksToGetTasksResponse(tasks []teamserver.Task) *api.GetTasksResponse {
	tasksResponse := api.GetTasksResponse{}
	for _, task := range tasks {
		tasksResponse.Tasks = append(tasksResponse.Tasks, task)
	}
	return &tasksResponse
}

func PostTaskResultsRequestToTaskResultsIn(taskResults *api.PostTaskResultsRequest) []teamserver.TaskResultIn {
	domainTaskResults := []teamserver.TaskResultIn{}
	for _, taskResult := range taskResults.TaskResults {
		domainTaskResults = append(domainTaskResults, taskResult)
	}

	return domainTaskResults
}

func TaskResultsOutToGetTaskResultsResponse(taskResultsOut []teamserver.TaskResultOut) *api.GetTaskResultsResponse {
	taskResultsResponse := api.GetTaskResultsResponse{}
	for _, taskResultOut := range taskResultsOut {
		taskResultsResponse.TaskResults = append(taskResultsResponse.TaskResults, struct {
			TaskId uint64
			Task   string
			Output string
		}{
			TaskId: taskResultOut.TaskId,
			Task:   taskResultOut.Task,
			Output: taskResultOut.Output,
		})
	}

	return &taskResultsResponse
}

func AgentViewsToListAgentsResponse(agentViews []teamserver.AgentView) *api.ListAgentsResponse {
	var listAgentsResponse api.ListAgentsResponse

	for _, agentView := range agentViews {
		listAgentsResponse.AgentViews = append(listAgentsResponse.AgentViews, agentView)
	}

	return &listAgentsResponse
}
