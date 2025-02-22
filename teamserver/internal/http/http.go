package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/http/middleware"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/ui"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Server struct {
	mux                *http.ServeMux
	server             *http.Server
	AgentService       *service.AgentService
	TaskQueueService   *service.TaskQueueService
	TaskResultsService *service.TaskResultsService
}

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

	server.registerAgentApiRouter()
	server.registerFrontendRouter()

	return &server
}

func (server *Server) registerAgentApiRouter() {
	router := http.NewServeMux()
	router.HandleFunc("POST /register", server.Register)
	router.HandleFunc("POST /taskQueue", server.PushTask)
	router.HandleFunc("GET /agents/{id}/tasks", server.GetTasks)
	router.HandleFunc("POST /agents/{id}/results", server.SaveTaskResults)
	router.HandleFunc("GET /agents/{id}/results", server.GetTaskResults)

	middlewareStack := middleware.CreateStack(middleware.JsonContentType)
	server.mux.Handle("/api/", middlewareStack(http.StripPrefix("/api", router)))
}

func (server *Server) registerFrontendRouter() {
	router := http.NewServeMux()
	views := ui.NewUi(server.AgentService, server.TaskQueueService, server.TaskResultsService)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.HandleFunc("/", views.RootHandler)
	server.mux.Handle("/", router)
}

func (server *Server) ListenAndServe() {
	log.Print("Starting listening on: " + server.server.Addr)
	log.Fatal(server.server.ListenAndServe())
}

func (server *Server) ListenAndServeTLS(cert string, key string) {
	log.Print("Starting listening on: " + server.server.Addr)
	log.Fatal(server.server.ListenAndServeTLS(cert, key))
}
