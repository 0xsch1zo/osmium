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
	GetAgentTaskProgress(agentId uint64) (uint64, error)
	UpdateAgentTaskProgress(agentId uint64) error
	AgentExists(agentId uint64) (bool, error)
}

type TasksRepository interface {
	AddTask(agentId uint64, task string) (uint64, error)
	GetTasks(agentId uint64, taskId uint64) ([]teamserver.Task, error)
	TaskExists(agentId uint64, taskId uint64) (bool, error)
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

type AgentService struct {
	agentRepository AgentRepository
}

type TasksService struct {
	agentService    *AgentService
	tasksRepository TasksRepository
}

type TaskResultsService struct {
	agentService          *AgentService
	tasksService          *TasksService
	taskResultsRepository TaskResultsRepository
}

type AuthorizationService struct {
	jwtKey                  string
	authorizationRepository AuthorizationRepository
}

func NewAgentService(agentRepository AgentRepository) *AgentService {
	return &AgentService{
		agentRepository: agentRepository,
	}
}

func NewTasksService(tasksRepository TasksRepository, agentRepository AgentRepository) *TasksService {
	return &TasksService{
		agentService:    NewAgentService(agentRepository),
		tasksRepository: tasksRepository,
	}
}

func NewTaskResultsService(
	taskResultsRepository TaskResultsRepository,
	agentRepository AgentRepository,
	tasksRepository TasksRepository,
) *TaskResultsService {
	return &TaskResultsService{
		agentService:          NewAgentService(agentRepository),
		tasksService:          NewTasksService(tasksRepository, agentRepository),
		taskResultsRepository: taskResultsRepository,
	}
}

func NewAuthorizationService(authorizationRepository AuthorizationRepository, jwtKey string) *AuthorizationService {
	return &AuthorizationService{
		jwtKey:                  jwtKey,
		authorizationRepository: authorizationRepository,
	}
}
