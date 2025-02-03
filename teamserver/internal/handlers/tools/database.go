package tools

type Database interface {
	AddAgent() (uint64, error)
	SetupDatabase() error
}

type Agent struct {
	AgentId      uint64
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
