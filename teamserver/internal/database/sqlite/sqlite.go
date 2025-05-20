package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Sqlite struct {
	databaseHandle *sql.DB
}

type AgentRepository struct {
	databaseHandle *sql.DB
}

type TasksRepository struct {
	databaseHandle *sql.DB
}

type TaskResultsRepository struct {
	databaseHandle *sql.DB
}

type AuthorizationRepository struct {
	databaseHandle *sql.DB
}

type EventLogRepository struct {
	databaseHandle *sql.DB
}

func (s *Sqlite) NewAgentRepository() *service.AgentRepository {
	var a service.AgentRepository = &AgentRepository{
		databaseHandle: s.databaseHandle,
	}
	return &a
}

func (s *Sqlite) NewTasksRepository() *service.TasksRepository {
	var ts service.TasksRepository = &TasksRepository{
		databaseHandle: s.databaseHandle,
	}
	return &ts
}

func (s *Sqlite) NewTaskResultsRepository() *service.TaskResultsRepository {
	var trr service.TaskResultsRepository = &TaskResultsRepository{
		databaseHandle: s.databaseHandle,
	}
	return &trr
}

func (s *Sqlite) NewAuthorizationRepository() *service.AuthorizationRepository {
	var auth service.AuthorizationRepository = &AuthorizationRepository{
		databaseHandle: s.databaseHandle,
	}
	return &auth
}

func (s *Sqlite) NewEventLogRepository() *service.EventLogRepository {
	var er service.EventLogRepository = &EventLogRepository{
		databaseHandle: s.databaseHandle,
	}

	return &er
}

// Shitty code use migrations or something
func SetupDatabase(sourceString string) (*Sqlite, error) {
	databaseHandle, err := sql.Open("sqlite3", sourceString)
	if err != nil {
		return nil, err
	}

	query := `
CREATE TABLE IF NOT EXISTS Tasks(
    AgentId INT NOT NULL,
    TaskId INT NOT NULL,
    Task VARCHAR NOT NULL,
    Status INT NOT NUll,
    FOREIGN KEY (AgentId) REFERENCES Agents(AgentId) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Agents(
    AgentId INTEGER PRIMARY KEY,
    PrivateKey VARCHAR NOT NULL,
	Username VARCHAR NOT NULL,
	Hostname VARCHAR NOT NULL,
	LastCallback INT NOT NULL
); 

CREATE TABLE IF NOT EXISTS TaskResults(
    AgentId INT NOT NULL,
    TaskId INT NOT NULL,
    Output VARCHAR NOT NULL,
    FOREIGN KEY (AgentId) REFERENCES Agents(AgentId) ON DELETE CASCADE,
    FOREIGN KEY (TaskId)  REFERENCES TaskQueue(TaskId) ON DELETE CASCADE,
    UNIQUE (AgentId, TaskId)
);

DROP TABLE IF EXISTS Users;
CREATE TABLE IF NOT EXISTS Users(
    UserId INT PRIMARY KEY,
    Username TEXT,
    PasswordHash TEXT,
    UNIQUE(Username)
);

CREATE TABLE IF NOT EXISTS EventLog(
	Type INTEGER,
	Time INTEGER,
	Contents TEXT
);`
	_, err = databaseHandle.Exec(query)

	return &Sqlite{
		databaseHandle: databaseHandle,
	}, err
}
