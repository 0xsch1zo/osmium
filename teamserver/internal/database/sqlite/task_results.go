package sqlite

import (
	"fmt"
	"strings"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (taskResultsService *TaskResultsService) SaveTaskResults(agentId uint64, taskResults []teamserver.TaskResultIn) error {
	exists, err := taskResultsService.agentService.AgentExists(agentId)
	if err != nil {
		return err
	} else if !exists {
		return teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId))
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO TaskResults (AgentId, TaskId, Output) VALUES")
	values := []interface{}{}

	for _, taskResults := range taskResults {
		queryBuilder.WriteString("(?, ?, ?),")

		exists, err := taskResultsService.taskQueueService.TaskExists(taskResults.TaskId)
		if err != nil {
			return err
		} else if !exists {
			return teamserver.NewClientError(fmt.Sprintf(errTaskIdNotFoundFmt, taskResults.TaskId))
		}

		values = append(values, agentId, taskResults.TaskId, taskResults.Output)
	}

	query := queryBuilder.String()
	query = strings.TrimSuffix(query, ",")

	stmt, err := taskResultsService.databaseHandle.Prepare(query)
	if err != nil {
		return teamserver.NewServerError(err.Error())
	}

	_, err = stmt.Exec(values...)
	if err != nil {
		return teamserver.NewServerError(err.Error())
	}
	return err
}

func (taskResultsService *TaskResultsService) GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error) {
	exists, err := taskResultsService.agentService.AgentExists(agentId)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	} else if !exists {
		return nil, teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId))
	}

	exists, err = taskResultsService.taskQueueService.TaskExists(taskId)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	} else if !exists {
		return nil, teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, taskId))
	}

	query := "SELECT Task, Output FROM TaskResults INNER JOIN TaskQueue ON TaskResults.TaskId = TaskQueue.TaskId WHERE agentId = ? AND taskId = ?"
	taskResultsSqlRow := taskResultsService.databaseHandle.QueryRow(query)
	taskResult := teamserver.TaskResultOut{}
	err = taskResultsSqlRow.Scan(&taskResult.TaskId, &taskResult.Task, &taskResult.Output)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	return &taskResult, nil
}

func (taskResultsService *TaskResultsService) GetTaskResults(agentId uint64) ([]teamserver.TaskResultOut, error) {
	exists, err := taskResultsService.agentService.AgentExists(agentId)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	} else if !exists {
		return nil, teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId))
	}

	query := "SELECT TaskResults.TaskId, Task, Output FROM TaskResults INNER JOIN TaskQueue ON TaskResults.TaskId = TaskQueue.TaskId WHERE agentId = ?"
	taskResultsSqlRows, err := taskResultsService.databaseHandle.Query(query, agentId)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	taskResults := []teamserver.TaskResultOut{}
	for taskResultsSqlRows.Next() {
		taskResult := teamserver.TaskResultOut{}
		err := taskResultsSqlRows.Scan(&taskResult.TaskId, &taskResult.Task, &taskResult.Output)
		if err != nil {
			return nil, teamserver.NewServerError(err.Error())
		}

		taskResults = append(taskResults, taskResult)
	}

	return taskResults, nil
}
