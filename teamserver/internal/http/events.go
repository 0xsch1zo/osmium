package http

import (
	"bytes"
	"log"
	"net/http"
	"sync"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (s *Server) EventLogListen(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	s.EventLogService.AddOnEventLoggedCallback(func(event teamserver.Event) {
		buf := bytes.Buffer{}
		err := templates.EventOOB(s.EventLogService.FormatEvent(&event)).Render(r.Context(), &buf)
		if err != nil {
			log.Print(err)
		}
		err = sendSSE(w, "event", buf.String())
		if err != nil {
			log.Print(err)
		}
	})

	wg.Wait()
}
