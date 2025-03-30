package sqlite

import (
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (er *EventLogRepository) LogEvent(event *teamserver.Event) error {
	query := "INSERT INTO EventLog (Type, Time, Contents) VALUES(?, ?, ?)"
	_, err := er.databaseHandle.Exec(query, event.Type, event.Time.Unix(), event.Contents)
	return err
}

func (er *EventLogRepository) GetEventLog() ([]teamserver.Event, error) {
	query := "SELECT Type, Time, Contents FROM EventLog"
	results, err := er.databaseHandle.Query(query)
	if err != nil {
		return nil, err
	}

	var eventLog []teamserver.Event
	for results.Next() {
		var event teamserver.Event

		var unixTime int64
		err = results.Scan(&event.Type, &unixTime, &event.Contents)
		if err != nil {
			return nil, err
		}

		event.Time = time.Unix(unixTime, 0)
		eventLog = append(eventLog, event)
	}

	return eventLog, err
}
