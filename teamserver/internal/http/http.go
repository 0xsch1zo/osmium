package http

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/config"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/database"
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
	EventLogService      *service.EventLogService
}

func NewServer(config *config.Config, db *database.Database) (*Server, error) {
	mux := http.NewServeMux()

	agentRepo := (*db).NewAgentRepository()
	tasksRepo := (*db).NewTasksRepository()
	taskResultsRepo := (*db).NewTaskResultsRepository()
	authRepo := (*db).NewAuthorizationRepository()
	eventLogRepo := (*db).NewEventLogRepository()

	eventLogService := service.NewEventLogService(eventLogRepo)
	agentService := service.NewAgentService(agentRepo, eventLogService)
	tasksService := service.NewTasksService(tasksRepo, agentService, eventLogService)
	taskResultsService := service.NewTaskResultsService(taskResultsRepo, agentService, tasksService, eventLogService)
	authorizationService := service.NewAuthorizationService(authRepo, os.Getenv("JWT_SECRET"), eventLogService)

	server := Server{
		mux: mux,
		server: &http.Server{
			Addr:    ":" + strconv.FormatUint(uint64(config.Port), 10),
			Handler: mux,
		},
		config:               config,
		AgentService:         agentService,
		TasksService:         tasksService,
		TaskResultsService:   taskResultsService,
		AuthorizationService: authorizationService,
		EventLogService:      eventLogService,
	}

	server.registerAgentApiRouter()
	server.registerFrontendRouter()

	for _, user := range config.AuthorizedUsers {
		target := &teamserver.ClientError{}
		err := server.AuthorizationService.UsernameExists(user.Username)
		if errors.As(err, &target) {
			err := server.AuthorizationService.Register(user.Username, user.Password)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

	}

	return &server, nil
}

func (server *Server) registerAgentApiRouter() {
	router := http.NewServeMux()

	// Auth
	router.HandleFunc("POST /auth/login", server.Login)
	router.HandleFunc("POST /auth/refresh", server.RefreshToken)
	router.HandleFunc("GET /auth/refreshTime", server.GetRefreshTime)

	// agent
	router.HandleFunc("POST /agents/register", server.AgentRegister)
	router.HandleFunc("GET /agents/{agentId}/tasks", server.GetNewTasks)
	router.HandleFunc("POST /agents/{agentId}/results/{taskId}", server.SaveTaskResult)

	// user
	router.Handle("POST /agents/{agentId}/tasks", server.Authenticate(
		http.HandlerFunc(server.AddTask),
	))

	router.Handle("GET /agents/{agentId}/tasks/all", server.Authenticate(
		http.HandlerFunc(server.GetTasks),
	))

	router.Handle("GET /agents/{agentId}/results", server.Authenticate(
		http.HandlerFunc(server.GetTaskResults),
	))

	router.Handle("GET /agents/{agentId}/results/listen", server.Authenticate(
		ServerSentEvents(
			http.HandlerFunc(server.TaskResultsListen),
		)),
	)

	router.Handle("GET /agents/{agentId}/results/{taskId}", server.Authenticate(
		http.HandlerFunc(server.GetTaskResult),
	))

	router.Handle("GET /agents/{agentId}/socket", server.Authenticate(
		http.HandlerFunc(server.AgentSocket),
	))

	router.Handle("GET /eventLog", server.Authenticate(
		ServerSentEvents(
			http.HandlerFunc(server.EventLogListen),
		)),
	)

	router.Handle("GET /agents/register/listen", server.Authenticate(
		ServerSentEvents(
			http.HandlerFunc(server.AddAgentListen),
		)),
	)
	// reuse this for other agent updates if possible
	router.Handle("GET /agents/callbackTime/listen", server.Authenticate(
		ServerSentEvents(
			http.HandlerFunc(server.CallbackTimeUpdatedListen),
		)),
	)

	commonMiddleware := CreateStack(JsonContentType)
	server.mux.Handle("/api/", commonMiddleware(http.StripPrefix("/api", router)))
}

func (server *Server) registerFrontendRouter() {
	router := http.NewServeMux()
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Handle("/node_modules/", http.StripPrefix("/node_modules/", http.FileServer(http.Dir("./node_modules"))))
	router.Handle("/dist/", http.StripPrefix("/dist/", http.FileServer(http.Dir("./dist"))))
	router.HandleFunc("/", server.RootHandler)
	server.mux.Handle("/", router)
}

func (server *Server) ListenAndServe() error {
	log.Print("Starting listening on: " + server.server.Addr)
	if server.config.Https {
		return server.server.ListenAndServeTLS(server.config.CertificatePath, server.config.KeyPath)
	}
	return server.server.ListenAndServe()
}

func (server *Server) Close() {
	server.server.Close()
}
