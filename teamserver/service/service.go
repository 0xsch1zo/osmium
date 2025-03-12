package service

import (
	"crypto/rsa"
	"errors"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

const ErrAgentIdNotFoundFmt = "AgentId not found: %d"
const ErrTaskIdNotFoundFmt = "TaskId not found: %d"
const ErrAlreadyExistsFmt = "%s already exists"
const ErrEmptyString = "%s empty"
const TokenSize = 32

type RepositoryErrNotFound struct {
	Err string
}

type RepositoryErrAlreadyExists struct {
	Err string
}

type RepositoryErrInvalidCredentials struct{}

func (err *RepositoryErrNotFound) Error() string {
	return err.Err
}

func (err *RepositoryErrAlreadyExists) Error() string {
	return err.Err
}

func (err *RepositoryErrInvalidCredentials) Error() string {
	return "Invalid credentials"
}

func NewRepositoryErrNotFound(err string) *RepositoryErrNotFound {
	return &RepositoryErrNotFound{Err: err}
}

func NewRepositoryErrAlreadyExists(err string) *RepositoryErrAlreadyExists {
	return &RepositoryErrAlreadyExists{Err: err}
}

func NewRepositoryErrInvalidCredentials() *RepositoryErrInvalidCredentials {
	return &RepositoryErrInvalidCredentials{}
}

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
}

type AuthorizationRepository interface {
	Register(username, passwordHash string) error
	GetPasswordHash(username string) (string, error)
	SetSessionToken(username, sessionToken string) error
	GetSessionToken(username string) (string, error)
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

func NewAuthorizationService(authorizationRepository AuthorizationRepository) *AuthorizationService {
	return &AuthorizationService{
		authorizationRepository: authorizationRepository,
	}
}

func repositoryErrWrapper(err error) error {
	if err == nil {
		return nil
	}

	notFoundTarget := &RepositoryErrNotFound{}
	alreadyExistsTarget := &RepositoryErrAlreadyExists{}
	invalidCredentialsTarget := &RepositoryErrInvalidCredentials{}
	if errors.As(err, &notFoundTarget) ||
		errors.As(err, &alreadyExistsTarget) ||
		errors.As(err, &invalidCredentialsTarget) {
		return teamserver.NewClientError(err.Error())
	}

	return teamserver.NewServerError(err.Error())
}
