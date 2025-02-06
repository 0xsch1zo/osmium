package tools

import "crypto/rsa"

type Database interface {
	AddAgent() (*Agent, error)
	GetAgent(agentId uint64) (*Agent, error)
	GetTasks(agentId uint64) ([]string, error)
	//UpdateAgentTaskProgress() error
	TaskQueuePush(task string) error
	//TaskQueuePop() error
	SetupDatabase() error
}

type Agent struct {
	AgentId      uint64
	TaskProgress uint32
	PrivateKey   *rsa.PrivateKey
}

func NewDatabase() (*Database, error) {
	var db Database = &SQLiteDatabase{}

	err := db.SetupDatabase()
	if err != nil {
		return nil, err
	}

	return &db, nil
}
