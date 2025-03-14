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
	taskProgress, err := ts.agentService.GetAgentTaskProgress(agentId)
	if err != nil {
		return nil, err // GetAgentTaskProgress returns the custom error type already
	}

	tasks, err := ts.tasksRepository.GetTasks(agentId, taskProgress)
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
