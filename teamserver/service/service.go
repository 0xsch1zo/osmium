package service

import (
	"crypto/rsa"
	"errors"

	"github.com/sentientbottleofwine/osmium/teamserver"
)

const ErrAgentIdNotFoundFmt = "AgentId not found: %d"
const ErrTaskIdNotFoundFmt = "TaskId not found: %d"

type RepositoryErrNotFound struct {
	Err string
}

func (serr *RepositoryErrNotFound) Error() string {
	return serr.Err
}

func NewRepositoryErrNotFound(err string) *RepositoryErrNotFound {
	return &RepositoryErrNotFound{Err: err}
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

func repositoryErrWrapper(err error) error {
	if err == nil {
		return nil
	}

	target := &RepositoryErrNotFound{}
	if errors.As(err, &target) {
		return teamserver.NewClientError(err.Error())
	}

	return teamserver.NewServerError(err.Error())
}
