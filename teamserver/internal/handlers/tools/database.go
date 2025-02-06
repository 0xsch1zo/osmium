package tools

import "crypto/rsa"

type Database interface {
	AddAgent() (*Agent, error)
	SetupDatabase() error
}

type Agent struct {
	AgentId      uint64
	PrivateKey   *rsa.PrivateKey
	TaskProgress uint32
}

func NewDatabase() (*Database, error) {
	var db Database = &SQLiteDatabase{}

	err := db.SetupDatabase()
	if err != nil {
		return nil, err
	}

	return &db, nil
}
