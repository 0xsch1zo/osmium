package service

import (
	"fmt"
	"net/http"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (ts *TasksService) AddTask(agentId uint64, task string) (uint64, error) {
	err := ts.agentService.AgentExists(agentId)
	if err != nil {
		return 0, err
	}

	taskId, err := ts.tasksRepository.AddTask(agentId, task)
	if err != nil {
		ServiceServerErrHandler(err, tasksServiceStr, ts.eventLogService)
		return 0, err
	}

	ts.eventLogService.LogEvent(
		teamserver.Info,
		fmt.Sprintf("Task was assigned for agent %d", agentId),
	)

	return taskId, nil
}

func (ts *TasksService) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	err := ts.agentService.AgentExists(agentId)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.tasksRepository.GetTasks(agentId)
	if err != nil {
		ServiceServerErrHandler(err, tasksServiceStr, ts.eventLogService)
		return nil, err
	}

	return tasks, nil
}

func (ts *TasksService) GetNewTasks(agentId uint64) ([]teamserver.Task, error) {
	err := ts.agentService.AgentExists(agentId)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.tasksRepository.GetTasksWithStatus(agentId, teamserver.TaskUnfinished)
	if err != nil {
		ServiceServerErrHandler(err, tasksServiceStr, ts.eventLogService)
		return nil, err
	}

	err = ts.agentService.UpdateLastCallbackTime(agentId)
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
		ServiceServerErrHandler(err, tasksServiceStr, ts.eventLogService)
		return err
	}

	if !exists {
		return teamserver.NewClientError(fmt.Sprintf(errTaskIdNotFoundFmt, taskId), http.StatusNotFound)
	}

	return nil
}

func (ts *TasksService) UpdateTaskStatus(agentId, taskId uint64, taskStatus teamserver.TaskStatus) error {
	err := ts.tasksRepository.UpdateTaskStatus(agentId, taskId, taskStatus)
	if err != nil {
		ServiceServerErrHandler(err, tasksServiceStr, ts.eventLogService)
		return err
	}
	return nil
}
