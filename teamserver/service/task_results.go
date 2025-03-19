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

	err = trs.agentService.UpdateAgentTaskProgress(agentId)
	if err != nil {
		return err
	}

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
