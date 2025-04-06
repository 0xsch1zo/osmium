package service

import (
	"strings"
	"sync"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (es *EventLogService) LogEvent(eventType teamserver.EventType, text string) error {
	event := &teamserver.Event{Type: eventType, Time: time.Now(), Contents: text}
	err := es.eventLogRepository.LogEvent(event)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(es.callbacks))
	for _, listener := range es.callbacks {
		if listener != nil {
			go func() {
				listener(event)
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
		eventFormat := es.FormatEvent(&event)
		eventLogFormat = append(eventLogFormat, eventFormat)
	}

	return eventLogFormat, nil
}

// Consider using normal function
func (es *EventLogService) FormatEvent(event *teamserver.Event) string {
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

	return eventFormatBuilder.String()
}

func (es *EventLogService) AddOnEventLoggedCallback(listener func(event *teamserver.Event)) teamserver.CallbackHandle {
	es.callbacks = append(es.callbacks, listener)
	return teamserver.CallbackHandle(len(es.callbacks) - 1)
}

func (es *EventLogService) RemoveOnEventLoggedCallback(listenerHandle teamserver.CallbackHandle) {
	for i := range es.callbacks {
		if teamserver.CallbackHandle(i) == listenerHandle {
			es.callbacks[i] = nil // Deleting the element would cause handles to be invalid
			break
		}
	}
}
