package tools

type Database interface {
	AddAgent(string, uint32) error
	SetupDatabase() error
}

type Agent struct {
	Uuid         string
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
