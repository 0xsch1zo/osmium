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
	TasksService       *service.TasksService
	TaskResultsService *service.TaskResultsService
}

func NewServer(port int, db *database.Database) *Server {
	mux := http.NewServeMux()

	agentRepo := (*db).NewAgentRepository()
	tasksRepo := (*db).NewTasksRepository()
	taskResultsRepo := (*db).NewTaskResultsRepository()

	server := Server{
		mux: mux,
		server: &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: mux,
		},
		AgentService:       service.NewAgentService(*agentRepo),
		TasksService:       service.NewTasksService(*tasksRepo, *agentRepo),
		TaskResultsService: service.NewTaskResultsService(*taskResultsRepo, *agentRepo, *tasksRepo),
	}

	server.registerAgentApiRouter()
	server.registerFrontendRouter()

	return &server
}

func (server *Server) registerAgentApiRouter() {
	router := http.NewServeMux()
	router.HandleFunc("POST /register", server.Register)

	router.HandleFunc("GET /agents/{agentId}/tasks", server.GetTasks)
	router.HandleFunc("POST /agents/{agentId}/tasks", server.AddTask)

	router.HandleFunc("GET /agents/{agentId}/results/{taskId}", server.GetTaskResult)
	router.HandleFunc("POST /agents/{agentId}/results/{taskId}", server.SaveTaskResult)

	middlewareStack := middleware.CreateStack(middleware.JsonContentType)
	server.mux.Handle("/api/", middlewareStack(http.StripPrefix("/api", router)))
}

func (server *Server) registerFrontendRouter() {
	router := http.NewServeMux()
	views := ui.NewUi(server.AgentService, server.TasksService, server.TaskResultsService)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("./node_modules"))))
	router.HandleFunc("/", views.RootHandler)
	server.mux.Handle("/", router)
}

func (server *Server) ListenAndServe() error {
	log.Print("Starting listening on: " + server.server.Addr)
	return server.server.ListenAndServe()
}

func (server *Server) ListenAndServeTLS(cert string, key string) error {
	log.Print("Starting listening on: " + server.server.Addr)
	return server.server.ListenAndServeTLS(cert, key)
}

func (server *Server) Close() {
	server.server.Close()
}
