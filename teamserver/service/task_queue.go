package service

import "github.com/sentientbottleofwine/osmium/teamserver"

func (tqs *TaskQueueService) taskExists(taskId uint64) (bool, error) {
	exists, err := tqs.taskQueueRepository.TaskExists(taskId)
	return exists, repositoryErrWrapper(err)
}

func (tqs *TaskQueueService) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	taskProgress, err := tqs.agentService.getAgentTaskProgress(agentId)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error()) // GetAgentTaskProgress returns the custom error type already
	}

	tasks, err := tqs.taskQueueRepository.GetTasks(agentId, taskProgress)
	if err != nil {
		return nil, repositoryErrWrapper(err)
	}

	return tasks, nil
}

func (tqs *TaskQueueService) TaskQueuePush(task string) error {
	return repositoryErrWrapper(tqs.taskQueueRepository.TaskQueuePush(task))
}
