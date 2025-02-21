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

type TaskQueueRepository interface {
	TaskQueuePush(task string) error
	GetTasks(agentId uint64, taskProgress uint64) ([]teamserver.Task, error)
	TaskExists(taskId uint64) (bool, error)
}

// TOOO: Figure out how to translate api model to domain model when they're not one to one
type TaskResultsRepository interface {
	SaveTaskResults(agentId uint64, taskResults []teamserver.TaskResultIn) error
	GetTaskResult(agentId uint64, taskId uint64) (*teamserver.TaskResultOut, error)
	GetTaskResults(agentId uint64) ([]teamserver.TaskResultOut, error)
}

type AgentService struct {
	agentRepository AgentRepository
}

type TaskQueueService struct {
	agentService        *AgentService
	taskQueueRepository TaskQueueRepository
}

type TaskResultsService struct {
	agentService          *AgentService
	taskQueueService      *TaskQueueService
	taskResultsRepository TaskResultsRepository
}

func NewAgentService(agentRepository AgentRepository) *AgentService {
	return &AgentService{
		agentRepository: agentRepository,
	}
}

func NewTaskQueueService(taskQueueRepository TaskQueueRepository, agentRepository AgentRepository) *TaskQueueService {
	return &TaskQueueService{
		agentService:        NewAgentService(agentRepository),
		taskQueueRepository: taskQueueRepository,
	}
}

func NewTaskResultsService(
	taskResultsRepository TaskResultsRepository,
	agentRepository AgentRepository,
	taskQueueRepository TaskQueueRepository,
) *TaskResultsService {
	return &TaskResultsService{
		agentService:          NewAgentService(agentRepository),
		taskQueueService:      NewTaskQueueService(taskQueueRepository, agentRepository),
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
