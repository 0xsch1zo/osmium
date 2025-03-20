package service

import (
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (ts *TasksService) AddTask(agentId uint64, task string) (uint64, error) {
	err := ts.agentService.AgentExists(agentId)
	if err != nil {
		return 0, err
	}

	taskId, err := ts.tasksRepository.AddTask(agentId, task)
	return taskId, err
}

func (ts *TasksService) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	tasks, err := ts.tasksRepository.GetTasks(agentId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (ts *TasksService) TaskExists(agentId uint64, taskId uint64) error {
	err := ts.agentService.AgentExists(agentId)
	if err != nil {
		return err
	}

	exists, err := ts.tasksRepository.TaskExists(agentId, taskId)
	if err != nil {
		return err
	}

	if !exists {
		return teamserver.NewClientError(fmt.Sprintf(errTaskIdNotFoundFmt, taskId))
	}

	return nil
}

func (ts *TasksService) UpdateTaskStatus(agentId, taskId uint64, taskStatus teamserver.TaskStatus) error {
	return ts.tasksRepository.UpdateTaskStatus(agentId, taskId, taskStatus)
}
