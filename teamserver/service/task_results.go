package service

import (
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (trs *TaskResultsService) SaveTaskResult(agentId uint64, taskResult *teamserver.TaskResultIn) error {
	exists, err := trs.TaskResultExists(agentId, taskResult.TaskId)
	if err != nil {
		return err
	}

	if exists {
		return teamserver.NewClientError(fmt.Sprintf(errAlreadyExistsFmt, "Task Result"))
	}

	err = trs.taskResultsRepository.SaveTaskResult(agentId, taskResult)
	if err != nil {
		return err
	}

	err = trs.tasksService.UpdateTaskStatus(agentId, taskResult.TaskId, teamserver.TaskFinished)
	if err != nil {
		return err
	}

	for _, callback := range trs.callbacks {
		if callback != nil {
			go callback(agentId, *taskResult)
		}
	}

	trs.eventLogService.LogEvent(
		teamserver.Info,
		fmt.Sprintf("A task result was recieved from agent %d", agentId))
	return nil
}

func (trs *TaskResultsService) GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error) {
	err := trs.agentService.AgentExists(agentId)
	if err != nil {
		return nil, err
	}

	err = trs.tasksService.TaskExists(agentId, taskId)
	if err != nil {
		return nil, err
	}

	taskResult, err := trs.taskResultsRepository.GetTaskResult(agentId, taskId)
	if err != nil {
		return nil, err
	}

	return taskResult, nil
}

func (trs *TaskResultsService) TaskResultExists(agentId, taskId uint64) (bool, error) {
	err := trs.agentService.AgentExists(agentId)
	if err != nil {
		return false, err
	}

	err = trs.tasksService.TaskExists(agentId, taskId)
	if err != nil {
		return false, err
	}

	exists, err := trs.taskResultsRepository.TaskResultExists(agentId, taskId)
	return exists, err
}

func (trs *TaskResultsService) AddOnTaskResultSavedCallback(listener func(agentId uint64, result teamserver.TaskResultIn)) teamserver.CallbackHandle {
	trs.callbacks = append(trs.callbacks, listener)
	return teamserver.CallbackHandle(len(trs.callbacks) - 1)
}

func (trs *TaskResultsService) RemoveOnTaskResultSavedCallback(listenerHandle teamserver.CallbackHandle) {
	for i := range trs.callbacks {
		if teamserver.CallbackHandle(i) == listenerHandle {
			trs.callbacks[i] = nil
		}
	}
}
