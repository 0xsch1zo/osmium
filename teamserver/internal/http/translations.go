package http

import (
	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
)

func PostTaskResultToDomain(taskResult *api.PostTaskResult) *teamserver.TaskResultIn {
	return &teamserver.TaskResultIn{
		TaskId: taskResult.TaskId,
		Output: taskResult.Output,
	}
}
