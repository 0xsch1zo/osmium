package service

import (
	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (ts *TasksService) AddTask(agentId uint64, task string) (uint64, error) {
	taskId, err := ts.tasksRepository.AddTask(agentId, task)
	return taskId, repositoryErrWrapper(err)
}

func (ts *TasksService) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	taskProgress, err := ts.agentService.GetAgentTaskProgress(agentId)
	if err != nil {
		return nil, err // GetAgentTaskProgress returns the custom error type already
	}

	tasks, err := ts.tasksRepository.GetTasks(agentId, taskProgress)
	if err != nil {
		return nil, repositoryErrWrapper(err)
	}

	return tasks, nil
}

func (ts *TasksService) TaskExists(agentId uint64, taskId uint64) (bool, error) {
	exists, err := ts.tasksRepository.TaskExists(agentId, taskId)
	return exists, repositoryErrWrapper(err)
}
