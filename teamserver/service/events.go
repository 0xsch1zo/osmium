package service

import (
	"strings"
	"sync"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (es *EventLogService) LogEvent(event *teamserver.Event) error {
	err := es.eventLogRepository.LogEvent(event)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(es.onEventLogged))
	for _, listener := range es.onEventLogged {
		if listener != nil {
			go func() {
				listener()
				wg.Done()
			}()
		}
	}

	return nil
}

func (es *EventLogService) GetEventLog() ([]string, error) {
	eventLog, err := es.eventLogRepository.GetEventLog()
	if err != nil {
		return nil, err
	}

	eventLogFormat := make([]string, 0, len(eventLog))
	for _, event := range eventLog {
		eventFormatBuilder := strings.Builder{}
		eventFormatBuilder.WriteString(event.Time.Format(time.RFC3339) + " ")

		switch event.Type {
		case teamserver.Info:
			eventFormatBuilder.WriteString("[+]")
		case teamserver.Warn:
			eventFormatBuilder.WriteString("[#]")
		case teamserver.Error:
			eventFormatBuilder.WriteString("[!]")
		}

		eventFormatBuilder.WriteString(" ")
		eventFormatBuilder.WriteString(event.Contents)

		eventLogFormat = append(eventLogFormat, eventFormatBuilder.String())
	}

	return eventLogFormat, nil
}

func (es *EventLogService) AddEventLoggedListener(listener func()) teamserver.EventListenerHandle {
	es.onEventLogged = append(es.onEventLogged, listener)
	return teamserver.EventListenerHandle(len(es.onEventLogged) - 1)
}

func (es *EventLogService) RemoveEventLoggedListener(listenerHandle teamserver.EventListenerHandle) {
	for i := range es.onEventLogged {
		if teamserver.EventListenerHandle(i) == listenerHandle {
			es.onEventLogged[i] = nil // Deleting the element would cause handles to be invalid
		}
	}
}
