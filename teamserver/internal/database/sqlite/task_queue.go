package sqlite

import (
	"database/sql"
)

func (tqr *TaskQueueRepository) TaskExists(taskId uint64) (bool, error) {
	query := "SELECT TaskId FROM TaskQueue WHERE TaskId = ?"
	sqlRow := tqr.databaseHandle.QueryRow(query, taskId)

	var temp uint64
	err := sqlRow.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (tqr *TaskQueueRepository) GetTaskQueue() ([]string, error) {
	query := "SELECT Task FROM TaskQueue ORDER BY TaskId ASC"
	sqlRow, err := tqr.databaseHandle.Query(query)
	if err != nil {
		return nil, err
	}

	var taskQueue []string
	for sqlRow.Next() {
		var task string
		err = sqlRow.Scan(&task)
		if err != nil {
			return nil, err
		}
		taskQueue = append(taskQueue, task)
	}

	return taskQueue, nil
}
func (tqr *TaskQueueRepository) TaskQueuePush(task string) error {
	query := "INSERT INTO TaskQueue VALUES(NULL, ?)"
	_, err := tqr.databaseHandle.Exec(query, task)
	return err
}
