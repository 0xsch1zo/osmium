package sqlite

import (
	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (trr *TaskResultsRepository) SaveTaskResult(agentId uint64, taskResult *teamserver.TaskResultIn) error {
	query := "INSERT INTO TaskResults (AgentId, TaskId, Output) VALUES(?, ?, ?)"

	_, err := trr.databaseHandle.Exec(query, agentId, taskResult.TaskId, taskResult.Output)
	if err != nil {
		return err
	}

	return err
}

func (trr *TaskResultsRepository) GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error) {
	query := "SELECT Task, Output FROM TaskResults INNER JOIN Tasks ON Tasks.TaskId = ? WHERE TaskResults.AgentId = ?"
	taskResultsSqlRow := trr.databaseHandle.QueryRow(query, taskId, agentId)
	taskResult := teamserver.TaskResultOut{}
	err := taskResultsSqlRow.Scan(&taskResult.Task, &taskResult.Output)
	if err != nil {
		return nil, err
	}

	taskResult.TaskId = taskId

	return &taskResult, nil
}
