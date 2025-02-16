package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

// Make sure that db implements domain serevices
var _ teamserver.AgentService = (*AgentService)(nil)
var _ teamserver.TaskQueueService = (*TaskQueueService)(nil)
var _ teamserver.TaskResultsService = (*TaskResultsService)(nil)

var errAgentIdNotFoundFmt string = "AgentId not found: %d"
var errTaskIdNotFoundFmt string = "TaskId not found: %d"

type AgentService struct {
	databaseHandle *sql.DB
}

type TaskQueueService struct {
	agentService   teamserver.AgentService
	databaseHandle *sql.DB
}

type TaskResultsService struct {
	databaseHandle   *sql.DB
	agentService     teamserver.AgentService
	taskQueueService teamserver.TaskQueueService
}

func NewAgentService(dbHandle *sql.DB) *AgentService {
	return &AgentService{
		databaseHandle: dbHandle,
	}
}

func NewTaskQueueService(dbHandle *sql.DB) *TaskQueueService {
	return &TaskQueueService{
		agentService:   NewAgentService(dbHandle),
		databaseHandle: dbHandle,
	}
}

func NewTaskResultsService(dbHandle *sql.DB) *TaskResultsService {
	return &TaskResultsService{
		databaseHandle:   dbHandle,
		agentService:     NewAgentService(dbHandle),
		taskQueueService: NewTaskQueueService(dbHandle),
	}
}

// Shitty code use migrations or something
func SetupDatabase() (*sql.DB, error) {
	databaseHandle, err := sql.Open("sqlite3", "teamserver.db")
	if err != nil {
		return nil, err
	}

	query := `
CREATE TABLE IF NOT EXISTS TaskQueue(
    TaskId INTEGER PRIMARY KEY,
    Task VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS Agents(
    AgentId INTEGER PRIMARY KEY,
    TaskProgress INT NOT NULL,
    PrivateKey VARCHAR NOT NULL,
    FOREIGN KEY (TaskProgress) REFERENCES TaskQueue(TaskId) ON DELETE CASCADE
); 

CREATE TABLE IF NOT EXISTS TaskResults(
    AgentId INT NOT NULL,
    TaskId INT NOT NULL,
    Output VARCHAR NOT NULL,
    FOREIGN KEY (AgentId) REFERENCES Agents(AgentId) ON DELETE CASCADE,
    FOREIGN KEY (TaskId)  REFERENCES TaskQueue(TaskId) ON DELETE CASCADE,
    UNIQUE (AgentId, TaskId)
);`
	_, err = databaseHandle.Exec(query)
	return databaseHandle, err
}

func (agentService *AgentService) AddAgent() (*teamserver.Agent, error) {
	rsaPriv, err := tools.GenerateKey()
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	query := "INSERT INTO Agents (AgentId, TaskProgress, PrivateKey) values(NULL, 0, ?);"
	_, err = agentService.databaseHandle.Exec(query, tools.PrivRsaToPem(rsaPriv))
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	// Get last row in db to get the AgentId of the newly created Agent
	query = "SELECT AgentId FROM Agents ORDER BY AgentId DESC LIMIT 1;" // in sqlite integer primary key will autoicrement as long as null is passed in
	AgentIdSqlRow := agentService.databaseHandle.QueryRow(query)

	var AgentId uint64
	err = AgentIdSqlRow.Scan(&AgentId)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	return &teamserver.Agent{
		AgentId:    AgentId,
		PrivateKey: rsaPriv,
	}, nil
}

func (agentService *AgentService) GetAgent(agentId uint64) (*teamserver.Agent, error) {
	query := "SELECT AgentId, TaskProgress, PrivateKey FROM Agents WHERE AgentId = ?"
	AgentSqlRow := agentService.databaseHandle.QueryRow(query, agentId)

	var agent teamserver.Agent
	var agentPrivateKeyPem string
	err := AgentSqlRow.Scan(&agent.AgentId, &agent.TaskProgress, &agentPrivateKeyPem)
	if err == sql.ErrNoRows {
		return nil, teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId))
	} else if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}

	agent.PrivateKey, err = tools.PemToPrivRsa(agentPrivateKeyPem)
	if err != nil {
		return nil, teamserver.NewServerError(err.Error())
	}
	return &agent, nil
}

func (agentService *AgentService) AgentExists(agentId uint64) (bool, error) {
	query := "SELECT AgentId FROM Agents WHERE AgentId = ?"
	sqlRow := agentService.databaseHandle.QueryRow(query, agentId)

	var temp uint64
	err := sqlRow.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, teamserver.NewServerError(err.Error())
	}

	return true, nil
}

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

func (agentService *AgentService) GetAgentTaskProgress(agentId uint64) (uint64, error) {
	query := "SELECT TaskProgress FROM Agents WHERE AgentId = ?"
	AgentSqlRow := agentService.databaseHandle.QueryRow(query, agentId)
	var taskProgress uint64
	err := AgentSqlRow.Scan(&taskProgress)
	if err == sql.ErrNoRows {
		return 0, teamserver.NewClientError(fmt.Sprintf(errAgentIdNotFoundFmt, agentId))
	} else if err != nil {
		return 0, teamserver.NewServerError(err.Error())
	}

	return taskProgress, nil
}

// TODO: Fix this shit
func (agentService *AgentService) UpdateAgentTaskProgress(agentId uint64) error {
	query := "UPDATE Agents SET TaskProgress = (SELECT MAX(TaskId) FROM TaskQueue)"
	_, err := agentService.databaseHandle.Exec(query)
	return err
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

/*
func (sqliteDb *SQLiteDatabase) TaskQueuePop() error {
	query := "DELETE FROM TaskQueue ORDER BY TaskId LIMIT 1;"
	_, err := sqliteDb.databaseHandle.Exec(query)
	return err
}*/
