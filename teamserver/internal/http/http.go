package http

import (
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/teamserver"
)

type Server struct {
	Server *http.Server
	teamserver.AgentService
	teamserver.TaskQueueService
	teamserver.TaskResultsService
}

func NewServer(port int, serveMux *http.ServeMux) *Server {
	return &Server{
		Server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
	}
}
