package tools

import (
	"database/sql"
)

type SQLiteDatabase struct {
	databaseHandle *sql.DB
}

func (sqliteDb *SQLiteDatabase) SetupDatabase() error {
	databaseHandle, err := sql.Open("sqlite3", "teamserver.db")
	if err != nil {
		return err
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

CREATE TABLE IF NOT EXISTS Failures(
    AgentId INT NOT NULL,
    TaskId INT NOT NULL,
    DateTime DATETIME NOT NULL,
    Error VARCHAR NOT NULL,
    FOREIGN KEY (AgentId) REFERENCES Agents(AgentId) ON DELETE CASCADE,
    FOREIGN KEY (TaskId)  REFERENCES TaskQueue(TaskId) ON DELETE CASCADE
);`
	_, err = databaseHandle.Exec(query)
	if err != nil {
		return err
	}

	sqliteDb.databaseHandle = databaseHandle
	return nil
}

func (sqliteDb *SQLiteDatabase) AddAgent() (*Agent, error) {
	rsaPriv, err := GenerateKey()
	if err != nil {
		return nil, err
	}

	query := "INSERT INTO Agents (AgentId, TaskProgress, PrivateKey) values(NULL, 0, ?);"
	_, err = sqliteDb.databaseHandle.Exec(query, PrivRsaToPem(rsaPriv))
	if err != nil {
		return nil, err
	}

	// Get last row in db to get the AgentId of the newly created Agent
	query = "SELECT AgentId FROM Agents ORDER BY AgentId DESC LIMIT 1;" // in sqlite integer primary key will autoicrement as long as null is passed in
	AgentIdSqlRow := sqliteDb.databaseHandle.QueryRow(query)

	var AgentId uint64
	err = AgentIdSqlRow.Scan(&AgentId)

	return &Agent{
		AgentId:    AgentId,
		PrivateKey: rsaPriv,
	}, err
}

func (sqliteDb *SQLiteDatabase) GetAgent(agentId uint64) (*Agent, error) {
	query := "SELECT AgentId, TaskProgress, PrivateKey FROM Agents WHERE AgentId = ?"
	AgentSqlRow := sqliteDb.databaseHandle.QueryRow(query, agentId)
	var agent Agent
	var agentPrivateKeyPem string
	err := AgentSqlRow.Scan(&agent.AgentId, &agent.TaskProgress, &agentPrivateKeyPem)
	if err != nil {
		return nil, err
	}

	agent.PrivateKey, err = PemToPrivRsa(agentPrivateKeyPem)
	return &agent, err
}

func (sqliteDb *SQLiteDatabase) GetAgentTaskProgress(agentId uint64) (uint64, error) {
	query := "SELECT TaskProgress FROM Agents WHERE AgentId = ?"
	AgentSqlRow := sqliteDb.databaseHandle.QueryRow(query, agentId)
	var taskProgress uint64
	err := AgentSqlRow.Scan(&taskProgress)
	return taskProgress, err
}

func (sqliteDb *SQLiteDatabase) UpdateAgentTaskProgress(agentId uint64) error {
	query := "UPDATE Agents SET TaskProgress = (SELECT MAX(TaskId) FROM TaskQueue)"
	_, err := sqliteDb.databaseHandle.Exec(query)
	return err
}

func (sqliteDb *SQLiteDatabase) GetTasks(agentId uint64) ([]string, error) {
	taskProgress, err := sqliteDb.GetAgentTaskProgress(agentId)
	if err != nil {
		return nil, err
	}

	query := "SELECT Task FROM TaskQueue WHERE TaskId >= ?"
	tasksSqlRows, err := sqliteDb.databaseHandle.Query(query, taskProgress)
	if err != nil {
		return nil, err
	}

	var tasks []string
	for tasksSqlRows.Next() {
		tasks = append(tasks, "")
		tasksSqlRows.Scan(&tasks[len(tasks)-1])
	}

	return tasks, nil
}

func (sqliteDb *SQLiteDatabase) TaskQueuePush(task string) error {
	query := "INSERT INTO TaskQueue values(NULL, ?)"
	_, err := sqliteDb.databaseHandle.Exec(query, task)
	return err
}

/*
func (sqliteDb *SQLiteDatabase) TaskQueuePop() error {
	query := "DELETE FROM TaskQueue ORDER BY TaskId LIMIT 1;"
	_, err := sqliteDb.databaseHandle.Exec(query)
	return err
}*/
