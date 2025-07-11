package http

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/templates"
)

func (s *Server) EventLogListen(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	handle := s.EventLogService.AddOnEventLoggedCallback(func(event teamserver.Event) {
		buf := bytes.Buffer{}
		err := templates.EventOOB(s.EventLogService.FormatEvent(&event)).Render(context.Background(), &buf)
		if err != nil {
			log.Print(err)
			return
		}
		err = sendSSE(w, "event", buf.String())
		if err != nil {
			log.Print(err)
			return
		}
	})

	defer s.EventLogService.RemoveOnEventLoggedCallback(handle)
	wg.Wait()
}
