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
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	query := "SELECT TaskId From Tasks WHERE AgentId = ? ORDER BY TaskId DESC LIMIT 1"
	row := tx.QueryRow(query, agentId)

	var lastTaskId uint64
	var taskId uint64
	err = row.Scan(&lastTaskId)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	} else {
		taskId = lastTaskId + 1
	}

	query = "INSERT INTO Tasks (AgentId, TaskId, Task, Status) VALUES(?, ?, ?, ?)"
	_, err = tx.Exec(query, agentId, taskId, task, teamserver.TaskUnfinished)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	return taskId, err
}

func (tr *TasksRepository) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	query := "SELECT TaskId, Task FROM Tasks WHERE AgentId = ? AND Status = ?"
	tasksSqlRows, err := tr.databaseHandle.Query(query, agentId, teamserver.TaskUnfinished)
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

func (tr *TasksRepository) UpdateTaskStatus(agentId uint64, taskId uint64, status teamserver.TaskStatus) error {
	query := "UPDATE Tasks SET Status = ? WHERE AgentId = ? AND TaskId = ?"
	_, err := tr.databaseHandle.Exec(query, status, agentId, taskId)
	return err
}
