package http

import (
	"bytes"
	"log"
	"net/http"
	"sync"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (s *Server) EventLogListen(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	s.EventLogService.AddEventLoggedListener(func() {
		eventLog, err := s.EventLogService.GetEventLog()
		if err != nil {
			sendSSE(w, "error", "An error occured, Couldn't retrieve the event log")
			log.Print(err)
			return
		}

		buf := bytes.Buffer{}
		templates.EventView(eventLog).Render(r.Context(), &buf)
		sendSSE(w, "eventLogView", buf.String())
	})
	wg.Add(1)

	wg.Wait()
}
