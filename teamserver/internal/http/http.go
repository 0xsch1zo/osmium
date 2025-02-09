package http

import (
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/teamserver"
)

type Server struct {
	server *http.Server
	teamserver.AgentService
	teamserver.TaskQueueService
	teamserver.TaskResultsService
}

func NewServer(port int, serveMux *http.ServeMux) *Server {
	return &Server{
		server: &http.Server{
			Addr: ":" + strconv.Itoa(port),
		},
	}
}
