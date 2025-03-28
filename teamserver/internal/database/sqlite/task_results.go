package sqlite

import (
	"database/sql"
	"log"

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
	query := "SELECT Task, Output FROM TaskResults INNER JOIN Tasks ON Tasks.TaskId = TaskResults.TaskId WHERE TaskResults.AgentId = ? AND TaskResults.TaskId = ?"
	taskResultsSqlRow := trr.databaseHandle.QueryRow(query, agentId, taskId)
	taskResult := teamserver.TaskResultOut{}
	err := taskResultsSqlRow.Scan(&taskResult.Task, &taskResult.Output)
	log.Print()
	if err != nil {
		return nil, err
	}

	taskResult.TaskId = taskId

	return &taskResult, nil
}

func (trr *TaskResultsRepository) TaskResultExists(agentId, taskId uint64) (bool, error) {
	query := "SELECT TaskId FROM TaskResults WHERE AgentId = ? AND TaskId = ?"
	row := trr.databaseHandle.QueryRow(query, agentId, taskId)

	var temp uint64
	err := row.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
