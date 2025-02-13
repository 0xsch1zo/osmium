package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

type Server struct {
	mux    *http.ServeMux
	server *http.Server
	teamserver.AgentService
	teamserver.TaskQueueService
	teamserver.TaskResultsService
}

func NewServer(port int) *Server {
	mux := http.NewServeMux()
	server := Server{
		mux: mux,
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
	}

	server.registerHandlers()

	return &server
}

func (server *Server) registerHandlers() {
	server.mux.HandleFunc("POST /register", server.Register)
	server.mux.HandleFunc("POST /taskQueue", server.PushTask)
	server.mux.HandleFunc("GET /agents/{id}/tasks", server.GetTasks)
	server.mux.HandleFunc("GET /agents/{id}/results", server.SaveTaskResults)
}

func (server *Server) ListenAndServe() {
	log.Print("Starting listening on: " + server.server.Addr)
	log.Fatal(server.server.ListenAndServe())
}

func (server *Server) ListenAndServeTLS(cert string, key string) {
	log.Print("Starting listening on: " + server.server.Addr)
	log.Fatal(server.server.ListenAndServeTLS(cert, key))
}
