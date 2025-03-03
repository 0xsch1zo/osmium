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

func PostTaskResultRequestToTaskResultsIn(taskResult *api.PostTaskResultRequest, taskId uint64) *teamserver.TaskResultIn {
	return &teamserver.TaskResultIn{
		TaskId: taskId,
		Output: taskResult.Output,
	}
}

func TaskResultsOutToGetTaskResultsResponse(taskResultOut *teamserver.TaskResultOut) *api.GetTaskResultResponse {
	return &api.GetTaskResultResponse{
		TaskId: taskResultOut.TaskId,
		Task:   taskResultOut.Task,
		Output: taskResultOut.Output,
	}
}

func AgentViewsToListAgentsResponse(agentViews []teamserver.AgentView) *api.ListAgentsResponse {
	var listAgentsResponse api.ListAgentsResponse

	for _, agentView := range agentViews {
		listAgentsResponse.AgentViews = append(listAgentsResponse.AgentViews, agentView)
	}

	return &listAgentsResponse
}
