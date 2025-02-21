package http

import (
	"log"
	"net/http"
	"strconv"

	//"github.com/a-h/templ"
	"github.com/sentientbottleofwine/osmium/teamserver/api"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/ui/templates"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Server struct {
	mux                *http.ServeMux
	server             *http.Server
	AgentService       *service.AgentService
	TaskQueueService   *service.TaskQueueService
	TaskResultsService *service.TaskResultsService
}

var counter uint64

func NewServer(port int, db *database.Database) *Server {
	mux := http.NewServeMux()

	agentRepo := (*db).NewAgentRepository()
	taskQueueRepo := (*db).NewTaskQueueRepository()
	taskResultsRepo := (*db).NewTaskResultsRepository()

	server := Server{
		mux: mux,
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
		AgentService:       service.NewAgentService(*agentRepo),
		TaskQueueService:   service.NewTaskQueueService(*taskQueueRepo, *agentRepo),
		TaskResultsService: service.NewTaskResultsService(*taskResultsRepo, *agentRepo, *taskQueueRepo),
	}

	server.registerHandlers()

	return &server
}

func (server *Server) registerHandlers() {
	server.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	server.mux.HandleFunc("POST /register", server.Register)
	server.mux.HandleFunc("POST /taskQueue", server.PushTask)
	server.mux.HandleFunc("GET /agents/{id}/tasks", server.GetTasks)
	server.mux.HandleFunc("POST /agents/{id}/results", server.SaveTaskResults)
	server.mux.HandleFunc("GET /agents/{id}/results", server.GetTaskResults)
	server.mux.HandleFunc("GET /agents", server.ListAgents)

	server.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := templates.Page(counter)
		err := component.Render(r.Context(), w)
		if err != nil {
			api.InternalErrorHandler(w)
		}
	})
	server.mux.HandleFunc("POST /count", func(w http.ResponseWriter, r *http.Request) {
		counter++
		component := templates.Counter(counter)
		err := component.Render(r.Context(), w)
		if err != nil {
			api.InternalErrorHandler(w)
		}
	})
}

func (server *Server) ListenAndServe() {
	log.Print("Starting listening on: " + server.server.Addr)
	log.Fatal(server.server.ListenAndServe())
}

func (server *Server) ListenAndServeTLS(cert string, key string) {
	log.Print("Starting listening on: " + server.server.Addr)
	log.Fatal(server.server.ListenAndServeTLS(cert, key))
}
