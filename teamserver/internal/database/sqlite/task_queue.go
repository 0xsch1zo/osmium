package sqlite

import (
	"database/sql"
	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (taskQueueService *TaskQueueService) TaskExists(taskId uint64) (bool, error) {
	query := "SELECT TaskId FROM TaskQueue WHERE TaskId = ?"
	sqlRow := taskQueueService.databaseHandle.QueryRow(query, taskId)

	var temp uint64
	err := sqlRow.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, teamserver.NewServerError(err.Error())
	}

	return true, nil
}

func (taskQueueService *TaskQueueService) GetTasks(agentId uint64) ([]teamserver.Task, error) {
	taskProgress, err := taskQueueService.agentService.GetAgentTaskProgress(agentId)
	if err != nil {
		return nil, err // GetAgentTaskProgress returns the custom error type already
	}

	query := "SELECT TaskId, Task FROM TaskQueue WHERE TaskId >= ?"
	tasksSqlRows, err := taskQueueService.databaseHandle.Query(query, taskProgress)
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

func (taskQueueService *TaskQueueService) TaskQueuePush(task string) error {
	query := "INSERT INTO TaskQueue VALUES(NULL, ?)"
	_, err := taskQueueService.databaseHandle.Exec(query, task)
	return teamserver.NewServerError(err.Error())
}
