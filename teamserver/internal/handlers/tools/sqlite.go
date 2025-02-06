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

	querry :=
		`CREATE TABLE IF NOT EXISTS Agents(
AgentId INTEGER PRIMARY KEY,
TaskProgress INT,
RSAPrivate VARCHAR
); 
CREATE TABLE IF NOT EXISTS TaskQueue(
TaskQueue VARCHAR
);`
	_, err = databaseHandle.Exec(querry)
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

	querry := "INSERT INTO Agents (AgentId, TaskProgress, RSAPrivate) values(NULL, 0, ?);"
	_, err = sqliteDb.databaseHandle.Exec(querry, PrivRsaToPem(rsaPriv))
	if err != nil {
		return nil, err
	}

	// Get last row in db to get the AgentId of the newly created Agent
	querry = "SELECT AgentId FROM Agents ORDER BY AgentId DESC LIMIT 1;" // in sqlite integer primary key will autoicrement as long as null is passed in
	AgentIdSqlRow := sqliteDb.databaseHandle.QueryRow(querry)

	var AgentId uint64
	err = AgentIdSqlRow.Scan(&AgentId)

	return &Agent{
		AgentId:    AgentId,
		PrivateKey: rsaPriv,
	}, err
}

/*
func (sqliteDb *SQLiteDatabase) SaveTaskQueue() error {
	querry := "INSERT INTO TaskQueue values(?)"
	_, err := sqliteDb.databaseHandle.Exec(querry)
	return err
}

func (sqliteDb *SQLiteDatabase) GetTaskQueue(Task string) error {
	querry := "INSERT INTO TaskQueue values(\"" + Task + "\");"
	_, err := sqliteDb.databaseHandle.Exec(querry)
	return err
}*/
