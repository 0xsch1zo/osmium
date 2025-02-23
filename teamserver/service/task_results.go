package service

import (
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (trs *TaskResultsService) validTaskResults(taskResults []teamserver.TaskResultIn) (bool, error) {
	for _, taskResult := range taskResults {
		exists, err := trs.taskQueueService.TaskExists(taskResult.TaskId)
		if err != nil {
			return false, err
		} else if !exists {
			return false, nil
		}
	}

	return true, nil
}

func (trs *TaskResultsService) SaveTaskResults(agentId uint64, taskResults []teamserver.TaskResultIn) error {
	valid, err := trs.agentService.AgentExists(agentId)
	if err != nil {
		return err
	} else if !valid {
		return teamserver.NewClientError(fmt.Sprintf(ErrAgentIdNotFoundFmt, agentId))
	}

	valid, err = trs.validTaskResults(taskResults)
	if err != nil {
		return err
	} else if !valid {
		return teamserver.NewClientError(ErrTaskIdNotFoundFmt)
	}

	err = trs.taskResultsRepository.SaveTaskResults(agentId, taskResults)
	if err != nil {
		return repositoryErrWrapper(err)
	}

	return nil
}

func (trs *TaskResultsService) GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error) {
	exists, err := trs.agentService.AgentExists(agentId)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, teamserver.NewClientError(fmt.Sprintf(ErrAgentIdNotFoundFmt, agentId))
	}

	exists, err = trs.taskQueueService.TaskExists(taskId)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, teamserver.NewClientError(fmt.Sprintf(ErrTaskIdNotFoundFmt, taskId))
	}

	taskResult, err := trs.taskResultsRepository.GetTaskResult(agentId, taskId)
	if err != nil {
		return nil, repositoryErrWrapper(err)
	}

	return taskResult, nil
}

func (trs *TaskResultsService) GetTaskResults(agentId uint64) ([]teamserver.TaskResultOut, error) {
	exists, err := trs.agentService.AgentExists(agentId)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, teamserver.NewClientError(fmt.Sprintf(ErrAgentIdNotFoundFmt, agentId))
	}

	taskResltuts, err := trs.taskResultsRepository.GetTaskResults(agentId)
	if err != nil {
		return nil, repositoryErrWrapper(err)
	}

	return taskResltuts, nil
}
