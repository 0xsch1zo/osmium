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
    FOREIGN KEY (AgentId) REFERENCES Agents(AgentId) ON DELETE CASCADE
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

	return &Sqlite{
		databaseHandle: databaseHandle,
	}, err
}
