package sqlite

import (
	"database/sql"
	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (tqr *TaskQueueRepository) TaskExists(taskId uint64) (bool, error) {
	query := "SELECT TaskId FROM TaskQueue WHERE TaskId = ?"
	sqlRow := tqr.databaseHandle.QueryRow(query, taskId)

	var temp uint64
	err := sqlRow.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, teamserver.NewServerError(err.Error())
	}

	return true, nil
}

func (tqr *TaskQueueRepository) GetTasks(agentId uint64, taskProgress uint64) ([]teamserver.Task, error) {
	query := "SELECT TaskId, Task FROM TaskQueue WHERE TaskId >= ?"
	tasksSqlRows, err := tqr.databaseHandle.Query(query, taskProgress)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	var tasks []teamserver.Task
	for tasksSqlRows.Next() {
		tasks = append(tasks, teamserver.Task{})
		err = tasksSqlRows.Scan(&(tasks[len(tasks)-1].TaskId), &(tasks[len(tasks)-1].Task))
		if err != nil {
			return nil, teamserver.NewServerError(err.Error())
		}
	}

	return tasks, nil
}

func (tqr *TaskQueueRepository) TaskQueuePush(task string) error {
	query := "INSERT INTO TaskQueue VALUES(NULL, ?)"
	_, err := tqr.databaseHandle.Exec(query, task)
	return teamserver.NewServerError(err.Error())
}
