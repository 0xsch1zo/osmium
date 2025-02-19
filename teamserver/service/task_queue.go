package service

import "github.com/sentientbottleofwine/osmium/teamserver"

func (tqs *TaskQueueService) taskExists(taskId uint64) (bool, error) {
	return tqs.taskQueueRepository.TaskExists(taskId)
}

func (tqs *TaskQueueService) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	taskProgress, err := tqs.agentService.getAgentTaskProgress(agentId)
	if err != nil {
		return nil, err // GetAgentTaskProgress returns the custom error type already
	}

	tasks, err := tqs.taskQueueRepository.GetTasks(agentId, taskProgress)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	return tasks, nil
}

func (tqs *TaskQueueService) TaskQueuePush(task string) error {
	return tqs.taskQueueRepository.TaskQueuePush(task)
}
