package sqlite

import (
	"database/sql"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (tr *TasksRepository) TaskExists(agentId uint64, taskId uint64) (bool, error) {
	query := "SELECT TaskId FROM Tasks WHERE TaskId = ? AND AgentId = ?"
	sqlRow := tr.databaseHandle.QueryRow(query, taskId, agentId)

	var temp uint64
	err := sqlRow.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (tr *TasksRepository) AddTask(agentId uint64, task string) (uint64, error) {
	tx, err := tr.databaseHandle.Begin()
	defer tx.Rollback()

	query := "SELECT TaskId From Tasks WHERE AgentId = ? ORDER BY TaskId DESC LIMIT 1"
	row := tx.QueryRow(query, agentId)

	var lastTaskId uint64
	var taskId uint64
	err = row.Scan(&lastTaskId)
	if err == sql.ErrNoRows {
		taskId = 1
	} else if err != nil {
		return 0, err
	} else {
		taskId = lastTaskId + 1
	}

	query = "INSERT INTO Tasks (AgentId, TaskId, Task) VALUES(?, ?, ?)"
	_, err = tx.Exec(query, agentId, taskId, task)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	return taskId, err
}

func (tr *TasksRepository) GetTasks(agentId uint64, taskProgress uint64) ([]teamserver.Task, error) {
	query := "SELECT TaskId, Task FROM Tasks WHERE AgentId = ? AND TaskId > ?"
	tasksSqlRows, err := tr.databaseHandle.Query(query, agentId, taskProgress)
	if err != nil {
		return nil, err
	}

	var tasks []teamserver.Task
	for tasksSqlRows.Next() {
		tasks = append(tasks, teamserver.Task{})
		err = tasksSqlRows.Scan(&(tasks[len(tasks)-1].TaskId), &(tasks[len(tasks)-1].Task))
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

/*
func (tqr *TaskQueueRepository) GetTaskQueue(agentId uint64) ([]string, error) {
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
*/
/*
func (as *AgentRepository) GetTasks(agentId uint64, taskProgress uint64) ([]teamserver.Task, error) {
	query := "SELECT TaskId, Task FROM TaskQueue WHERE TaskId >= ?"
	tasksSqlRows, err := as.databaseHandle.Query(query, taskProgress)
	if err != nil {
		return nil, err
	}

	var tasks []teamserver.Task
	for tasksSqlRows.Next() {
		tasks = append(tasks, teamserver.Task{})
		err = tasksSqlRows.Scan(&(tasks[len(tasks)-1].TaskId), &(tasks[len(tasks)-1].Task))
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

func (tqr *TaskQueueRepository) AddTasks(agentId uint64, task string) (uint64, error) {
	query := "INSERT INTO TaskQueue VALUES(NULL, ?)"
	_, err := tqr.databaseHandle.Exec(query, task)
	if err != nil {
		return 0, nil
	}

	query = "SELECT TaskId From TaskQueue ORDER BY TaskId DESC LIMIT 1"
	row := tqr.databaseHandle.QueryRow(query)

	var taskId uint64
	err = row.Scan(&taskId)
	if err != nil {
		return 0, err
	}

	return taskId, err
}*/
