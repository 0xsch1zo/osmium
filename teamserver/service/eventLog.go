package service

import (
	"log"
	"strings"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func (es *EventLogService) LogEvent(eventType teamserver.EventType, text string) {
	event := &teamserver.Event{Type: eventType, Time: time.Now(), Contents: text}
	err := es.eventLogRepository.LogEvent(event)
	if err != nil {
		log.Printf("Failed to log event: %s", err.Error())
	}

	for _, listener := range es.callbacks {
		if listener != nil {
			// Why the hell in go you can't pass by const-reference. WTF
			go listener(*event)
		}
	}
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
	eventFormatBuilder.WriteString(event.Time.Format(time.DateTime) + " ")

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

func (es *EventLogService) AddOnEventLoggedCallback(listener func(event teamserver.Event)) teamserver.CallbackHandle {
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
