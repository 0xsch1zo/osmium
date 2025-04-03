package http

import (
	"bytes"
	"net/http"
	"sync"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (s *Server) EventLogListen(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	s.EventLogService.AddOnEventLoggedCallback(func(event *teamserver.Event) {
		buf := bytes.Buffer{}
		templates.Event(s.EventLogService.FormatEvent(event)).Render(r.Context(), &buf)
		sendSSE(w, "event", buf.String())
	})

	wg.Wait()
}
