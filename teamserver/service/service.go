package service

import (
	"crypto/rsa"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

const (
	errAgentIdNotFoundFmt = "AgentId not found: %d"
	errTaskIdNotFoundFmt  = "TaskId not found: %d"
	errAlreadyExistsFmt   = "%s already exists"
	errEmptyString        = "%s empty"
	errTokenNotOld        = "Token is not old enugh"
	errInvalidCredentials = "Invalid credentials"
	jwtExpiryTime         = 15 * time.Minute
)

type AgentRepository interface {
	AddAgent(rsaPriv *rsa.PrivateKey) (*teamserver.Agent, error)
	GetAgent(agentId uint64) (*teamserver.Agent, error)
	ListAgents() ([]teamserver.AgentView, error)
	AgentExists(agentId uint64) (bool, error)
}

type TasksRepository interface {
	AddTask(agentId uint64, task string) (uint64, error)
	GetTasks(agentId uint64) ([]teamserver.Task, error)
	TaskExists(agentId uint64, taskId uint64) (bool, error)
	UpdateTaskStatus(agentId uint64, taskId uint64, status teamserver.TaskStatus) error
}

type TaskResultsRepository interface {
	SaveTaskResult(agentId uint64, taskResult *teamserver.TaskResultIn) error
	GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error)
	TaskResultExists(agentId, taskId uint64) (bool, error)
}

type AuthorizationRepository interface {
	Register(username, passwordHash string) error
	GetPasswordHash(username string) (string, error)
	UsernameExists(username string) (bool, error)
}

type EventLogRepository interface {
	LogEvent(event *teamserver.Event) error
	GetEventLog() ([]teamserver.Event, error)
}

type AgentService struct {
	agentRepository AgentRepository
	eventLogService *EventLogService
	callbacks       []func(teamserver.Agent)
}

type TasksService struct {
	agentService    *AgentService
	eventLogService *EventLogService
	tasksRepository TasksRepository
}

type TaskResultsService struct {
	agentService          *AgentService
	tasksService          *TasksService
	eventLogService       *EventLogService
	taskResultsRepository TaskResultsRepository
	callbacks             []func(agentId uint64, result teamserver.TaskResultIn)
}

type AuthorizationService struct {
	jwtKey                  string
	authorizationRepository AuthorizationRepository
	eventLogService         *EventLogService
}

type EventLogService struct {
	callbacks          []func(teamserver.Event)
	eventLogRepository EventLogRepository
}

func NewAgentService(agentRepository *AgentRepository, eventLogService *EventLogService) *AgentService {
	return &AgentService{
		agentRepository: *agentRepository,
		eventLogService: eventLogService,
	}
}

func NewTasksService(tasksRepository *TasksRepository, agentRepository *AgentRepository, eventLogService *EventLogService) *TasksService {
	return &TasksService{
		agentService:    NewAgentService(agentRepository, eventLogService),
		tasksRepository: *tasksRepository,
		eventLogService: eventLogService,
	}
}

func NewTaskResultsService(
	taskResultsRepository *TaskResultsRepository,
	agentRepository *AgentRepository,
	tasksRepository *TasksRepository,
	eventLogService *EventLogService,
) *TaskResultsService {
	return &TaskResultsService{
		agentService:          NewAgentService(agentRepository, eventLogService),
		tasksService:          NewTasksService(tasksRepository, agentRepository, eventLogService),
		taskResultsRepository: *taskResultsRepository,
		eventLogService:       eventLogService,
	}
}

func NewAuthorizationService(
	authorizationRepository *AuthorizationRepository,
	jwtKey string,
	eventLogService *EventLogService,
) *AuthorizationService {
	return &AuthorizationService{
		jwtKey:                  jwtKey,
		authorizationRepository: *authorizationRepository,
		eventLogService:         eventLogService,
	}
}

func NewEventLogService(eventLogRepository *EventLogRepository) *EventLogService {
	return &EventLogService{
		eventLogRepository: *eventLogRepository,
	}
}
