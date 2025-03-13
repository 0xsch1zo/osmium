package http

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/config"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/ui"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

type Server struct {
	mux                  *http.ServeMux
	server               *http.Server
	config               *config.Config
	AgentService         *service.AgentService
	TasksService         *service.TasksService
	TaskResultsService   *service.TaskResultsService
	AuthorizationService *service.AuthorizationService

	awaitedTaskIdChannel chan uint64
	taskResultsChannel   chan *teamserver.TaskResultIn
}

func NewServer(config *config.Config, db *database.Database) (*Server, error) {
	mux := http.NewServeMux()

	agentRepo := (*db).NewAgentRepository()
	tasksRepo := (*db).NewTasksRepository()
	taskResultsRepo := (*db).NewTaskResultsRepository()
	authRepo := (*db).NewAuthorizationRepository()

	server := Server{
		mux: mux,
		server: &http.Server{
			Addr:    ":" + strconv.FormatUint(uint64(config.Port), 10),
			Handler: mux,
		},
		config:               config,
		AgentService:         service.NewAgentService(*agentRepo),
		TasksService:         service.NewTasksService(*tasksRepo, *agentRepo),
		TaskResultsService:   service.NewTaskResultsService(*taskResultsRepo, *agentRepo, *tasksRepo),
		AuthorizationService: service.NewAuthorizationService(*authRepo, os.Getenv("JWT_SECRET")),
		awaitedTaskIdChannel: make(chan uint64),
		taskResultsChannel:   make(chan *teamserver.TaskResultIn),
	}

	server.registerAgentApiRouter()
	server.registerFrontendRouter()

	exists, err := server.AuthorizationService.UsernameExists(config.Username)
	if err != nil {
		return nil, err
	}

	if !exists {
		err := server.AuthorizationService.Register(config.Username, config.Password)
		if err != nil {
			return nil, err
		}
	}

	return &server, nil
}

func (server *Server) registerAgentApiRouter() {
	router := http.NewServeMux()

	// agent
	router.HandleFunc("POST /agents/register", server.AgentRegister)
	router.HandleFunc("GET /agents/{agentId}/tasks", server.GetTasks)
	router.HandleFunc("POST /agents/{agentId}/results/{taskId}", server.SaveTaskResult)

	// user
	router.HandleFunc("POST /auth/login", server.Login)
	router.Handle("POST /agents/{agentId}/tasks", server.Authenticate(http.HandlerFunc(server.AddTask)))
	router.Handle("GET /agents/{agentId}/results/{taskId}", server.Authenticate(http.HandlerFunc(server.GetTaskResult)))
	router.Handle("GET /agents/{agentId}/tasks/listen", ServerSentEvents(http.HandlerFunc(server.ListenAndServeTaskResults)))

	commonMiddleware := CreateStack(JsonContentType)
	server.mux.Handle("/api/", commonMiddleware(http.StripPrefix("/api", router)))
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
	if server.config.Https == true {
		return server.server.ListenAndServeTLS(server.config.CertificatePath, server.config.KeyPath)
	}
	return server.server.ListenAndServe()
}

func (server *Server) Close() {
	server.server.Close()
}
