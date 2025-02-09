package database

import (
	"database/sql"
	"strings"

	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/handlers/tools"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/osmium"
)

// Make sure that
//var _ AgentService = (*osmium.AgentService)(nil)

type AgentService struct {
	databaseHandle *sql.DB
}

type TaskQueueService struct {
	agentService   osmium.AgentService
	databaseHandle *sql.DB
}

type TaskResultsService struct {
	databaseHandle *sql.DB
}

// Implements functions of Database
//type SqliteDb struct{}

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
		databaseHandle: dbHandle,
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
    FOREIGN KEY (TaskId)  REFERENCES TaskQueue(TaskId) ON DELETE CASCADE
);`
	_, err = databaseHandle.Exec(query)
	return databaseHandle, err
}

func (agentService *AgentService) AddAgent() (*osmium.Agent, error) {
	rsaPriv, err := tools.GenerateKey()
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO Agents (AgentId, TaskProgress, PrivateKey) values(NULL, 0, ?);"
	_, err = agentService.databaseHandle.Exec(query, tools.PrivRsaToPem(rsaPriv))
	if err != nil {
		return nil, err
	}

	// Get last row in db to get the AgentId of the newly created Agent
	query = "SELECT AgentId FROM Agents ORDER BY AgentId DESC LIMIT 1;" // in sqlite integer primary key will autoicrement as long as null is passed in
	AgentIdSqlRow := agentService.databaseHandle.QueryRow(query)

	var AgentId uint64
	err = AgentIdSqlRow.Scan(&AgentId)

	return &osmium.Agent{
		AgentId:    AgentId,
		PrivateKey: rsaPriv,
	}, err
}

func (agentService *AgentService) GetAgent(agentId uint64) (*osmium.Agent, error) {
	query := "SELECT AgentId, TaskProgress, PrivateKey FROM Agents WHERE AgentId = ?"
	AgentSqlRow := agentService.databaseHandle.QueryRow(query, agentId)
	var agent osmium.Agent
	var agentPrivateKeyPem string
	err := AgentSqlRow.Scan(&agent.AgentId, &agent.TaskProgress, &agentPrivateKeyPem)
	if err != nil {
		return nil, err
	}

	agent.PrivateKey, err = tools.PemToPrivRsa(agentPrivateKeyPem)
	return &agent, err
}

func (agentService *AgentService) GetAgentTaskProgress(agentId uint64) (uint64, error) {
	query := "SELECT TaskProgress FROM Agents WHERE AgentId = ?"
	AgentSqlRow := agentService.databaseHandle.QueryRow(query, agentId)
	var taskProgress uint64
	err := AgentSqlRow.Scan(&taskProgress)
	return taskProgress, err
}

func (agentService *AgentService) UpdateAgentTaskProgress(agentId uint64) error {
	query := "UPDATE Agents SET TaskProgress = (SELECT MAX(TaskId) FROM TaskQueue)"
	_, err := agentService.databaseHandle.Exec(query)
	return err
}

func (taskQueueService *TaskQueueService) GetTasks(agentId uint64) ([]osmium.Task, error) {
	taskProgress, err := taskQueueService.agentService.GetAgentTaskProgress(agentId)
	if err != nil {
		return nil, err
	}

	query := "SELECT TaskId, Task FROM TaskQueue WHERE TaskId >= ?"
	tasksSqlRows, err := taskQueueService.databaseHandle.Query(query, taskProgress)
	if err != nil {
		return nil, err
	}

	var tasks []osmium.Task
	for tasksSqlRows.Next() {
		tasks = append(tasks, osmium.Task{})
		err = tasksSqlRows.Scan(&(tasks[len(tasks)-1].TaskId), &(tasks[len(tasks)-1].Task))
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}

func (taskQueueService *TaskQueueService) TaskQueuePush(task string) error {
	query := "INSERT INTO TaskQueue VALUES(NULL, ?)"
	_, err := taskQueueService.databaseHandle.Exec(query, task)
	return err
}

func (taskResultsService *TaskResultsService) SaveTaskResults(agentId uint64, taskResults api.PostTaskResultsRequest) error {
	var queryBuilder strings.Builder
	queryBuilder.WriteString("INSERT INTO TaskResults (AgentId, TaskId, Output) VALUES")
	values := []interface{}{}

	for _, taskResults := range taskResults.TaskResults {
		queryBuilder.WriteString("(?, ?, ?),")
		values = append(values, agentId, taskResults.TaskId, taskResults.Output)
	}

	query := queryBuilder.String()
	query = strings.TrimSuffix(query, ",")

	stmt, err := taskResultsService.databaseHandle.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(values...)
	return err
}

func (taskResultsService *TaskResultsService) GetTaskResult(agentId uint64, taskId uint64) (*osmium.TaskResult, error) {
	query := "SELECT Output, Task FROM TaskResults WHERE agentId = ? AND taskId = ? INNER JOIN Tasks ON TaskResults.TaskId = Tasks.TaskId"
	taskResultsSqlRow := taskResultsService.databaseHandle.QueryRow(query)
	taskResult := osmium.TaskResult{}
	err := taskResultsSqlRow.Scan(&taskResult.TaskId, &taskResult.Task, &taskResult.Output)
	if err != nil {
		return nil, err
	}

	return &taskResult, nil
}

func (taskResultsService *TaskResultsService) GetTaskResults(agentId uint64) ([]osmium.TaskResult, error) {
	query := "SELECT TaskResults.TaskId,Task, Output, FROM TaskResults WHERE agentId = ? INNER JOIN Tasks ON TaskResults.TaskId = Tasks.TaskId"
	taskResultsSqlRows, err := taskResultsService.databaseHandle.Query(query)
	if err != nil {
		return nil, err
	}

	taskResults := []osmium.TaskResult{}
	for taskResultsSqlRows.Next() {
		taskResult := osmium.TaskResult{}
		err := taskResultsSqlRows.Scan(&taskResult.TaskId, &taskResult.Task, &taskResult.Output)
		if err != nil {
			return nil, err
		}
	}

	return taskResults, nil
}

/*
func (sqliteDb *SQLiteDatabase) TaskQueuePop() error {
	query := "DELETE FROM TaskQueue ORDER BY TaskId LIMIT 1;"
	_, err := sqliteDb.databaseHandle.Exec(query)
	return err
}*/
