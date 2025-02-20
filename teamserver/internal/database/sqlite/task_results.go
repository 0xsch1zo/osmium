package sqlite

import (
	"strings"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (trr *TaskResultsRepository) SaveTaskResults(agentId uint64, taskResults []teamserver.TaskResultIn) error {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO TaskResults (AgentId, TaskId, Output) VALUES")
	values := []interface{}{}

	for _, taskResults := range taskResults {
		queryBuilder.WriteString("(?, ?, ?),")
		values = append(values, agentId, taskResults.TaskId, taskResults.Output)
	}

	query := queryBuilder.String()
	query = strings.TrimSuffix(query, ",")

	stmt, err := trr.databaseHandle.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)
	return err
}

func (trr *TaskResultsRepository) GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error) {
	query := "SELECT Task, Output FROM TaskResults INNER JOIN TaskQueue ON TaskResults.TaskId = TaskQueue.TaskId WHERE agentId = ? AND taskId = ?"
	taskResultsSqlRow := trr.databaseHandle.QueryRow(query)
	taskResult := teamserver.TaskResultOut{}
	err := taskResultsSqlRow.Scan(&taskResult.TaskId, &taskResult.Task, &taskResult.Output)
	if err != nil {
		return nil, err
	}

	return &taskResult, nil
}

func (trr *TaskResultsRepository) GetTaskResults(agentId uint64) ([]teamserver.TaskResultOut, error) {
	query := "SELECT TaskResults.TaskId, Task, Output FROM TaskResults INNER JOIN TaskQueue ON TaskResults.TaskId = TaskQueue.TaskId WHERE agentId = ?"
	taskResultsSqlRows, err := trr.databaseHandle.Query(query, agentId)
	if err != nil {
		return nil, err
	}

	taskResults := []teamserver.TaskResultOut{}
	for taskResultsSqlRows.Next() {
		taskResult := teamserver.TaskResultOut{}
		err := taskResultsSqlRows.Scan(&taskResult.TaskId, &taskResult.Task, &taskResult.Output)
		if err != nil {
			return nil, err
		}

		taskResults = append(taskResults, taskResult)
	}

	return taskResults, nil
}
