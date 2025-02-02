package tools

import "database/sql"

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
AgentId VARCHAR(36),
TaskProgress INT
); 
CREATE TABLE IF NOT EXISTS TaskQueue
TaskQueue VARCHAR`
	_, err = databaseHandle.Exec(querry)
	if err != nil {
		return err
	}

	sqliteDb.databaseHandle = databaseHandle
	return nil
}

func (sqliteDb *SQLiteDatabase) AddAgent(uuid string, taskProgress uint32) error {
	querry := "INSERT INTO Agents values(\"" + uuid + "\", 0);"
	_, err := sqliteDb.databaseHandle.Exec(querry)
	return err
}

/*func (sqliteDb *SQLiteDatabase) SaveTaskQueue() error {
	querry := "INSERT INTO TaskQueue values(\"" + Task + "\");"
	_, err := sqliteDb.databaseHandle.Exec(querry)
	return err
}

func (sqliteDb *SQLiteDatabase) GetTaskQueue(Task string) error {
	querry := "INSERT INTO TaskQueue values(\"" + Task + "\");"
	_, err := sqliteDb.databaseHandle.Exec(querry)
	return err
}*/
