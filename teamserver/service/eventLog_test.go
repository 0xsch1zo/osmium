package service_test

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

func TestEventLog(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	eventGiven := struct {
		teamserver.EventType
		text string
	}{teamserver.Info, "testing"}
	err = testedServices.eventLogService.LogEvent(eventGiven.EventType, eventGiven.text)
	if err != nil {
		t.Fatal(err)
	}

	eventLogGot, err := testedServices.eventLogService.GetEventLog()
	if err != nil {
		t.Fatal(err)
	}

	if len(eventLogGot) != 1 {
		t.Fatal("Event log has wrong size")
	}

	if !strings.Contains(eventLogGot[0], eventGiven.text) {
		t.Fatal("Event log returned does not contain the original logged value")
	}
}

func TestAddEventLogCallback(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	eventGiven := struct {
		teamserver.EventType
		text string
	}{teamserver.Info, "testing"}

	wg := sync.WaitGroup{}
	wg.Add(1)
	_ = testedServices.eventLogService.AddOnEventLoggedCallback(func(event teamserver.Event) {
		if event.Contents != eventGiven.text &&
			event.Type != eventGiven.EventType {
			t.Fatal("Event recieved on callback doesn't match with the event originally logged")
		}
		wg.Done()
	})

	err = testedServices.eventLogService.LogEvent(eventGiven.EventType, eventGiven.text)
	if err != nil {
		t.Fatal(err)
	}
	wg.Wait()
}

func TestRemoveEventLogCallback(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan struct{})
	cHandle := testedServices.eventLogService.AddOnEventLoggedCallback(func(event teamserver.Event) {
		ch <- struct{}{}
	})

	testedServices.eventLogService.RemoveOnEventLoggedCallback(cHandle)

	err = testedServices.eventLogService.LogEvent(teamserver.Warn, "garbage")
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-ch:
		t.Fatal("removed callback function has been run")
	case <-time.After(1 * time.Second):
		break
	}

}
