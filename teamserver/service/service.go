package service

import (
	"crypto/rsa"
	"fmt"
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
	jwtRefreshWindow      = 30 * time.Second
)

type ServiceString string

const (
	agentServiceStr         ServiceString = "Agent service: "
	tasksServiceStr         ServiceString = "Tasks service: "
	taskResultsServiceStr   ServiceString = "Tasks results service: "
	authorizationServiceStr ServiceString = "Authorization service: "
)

type AgentRepository interface {
	AddAgent(rsaPriv *rsa.PrivateKey, agentInfo teamserver.AgentRegisterInfo) (*teamserver.Agent, error)
	GetAgent(agentId uint64) (*teamserver.Agent, error)
	UpdateLastCallbackTime(agentId uint64) error
	ListAgents() ([]teamserver.AgentView, error)
	AgentExists(agentId uint64) (bool, error)
}

type TasksRepository interface {
	AddTask(agentId uint64, task string) (uint64, error)
	GetTasks(agentId uint64) ([]teamserver.Task, error)
	GetTasksWithStatus(agentId uint64, status teamserver.TaskStatus) ([]teamserver.Task, error)
	TaskExists(agentId uint64, taskId uint64) (bool, error)
	UpdateTaskStatus(agentId uint64, taskId uint64, status teamserver.TaskStatus) error
}

type TaskResultsRepository interface {
	SaveTaskResult(agentId uint64, taskResult *teamserver.TaskResultIn) error
	GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error)
	GetTaskResults(agentId uint64) ([]teamserver.TaskResultOut, error)
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
	agentRepository                AgentRepository
	eventLogService                *EventLogService
	onAgentAddedCallbacks          []func(teamserver.Agent)
	onCallbackTimeUpdatedCallbacks []func(teamserver.Agent)
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

func NewTasksService(tasksRepository *TasksRepository, agentService *AgentService, eventLogService *EventLogService) *TasksService {
	return &TasksService{
		agentService:    agentService,
		tasksRepository: *tasksRepository,
		eventLogService: eventLogService,
	}
}

func NewTaskResultsService(
	taskResultsRepository *TaskResultsRepository,
	agentService *AgentService,
	tasksRepository *TasksRepository,
	eventLogService *EventLogService,
) *TaskResultsService {
	return &TaskResultsService{
		agentService:          agentService,
		tasksService:          NewTasksService(tasksRepository, agentService, eventLogService),
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

func ServiceServerErrHandler(err error, service ServiceString, eventLogService *EventLogService) {
	err = fmt.Errorf("Server error: %s%w", service, err)
	eventLogService.LogEvent(teamserver.Error, err.Error())
}
