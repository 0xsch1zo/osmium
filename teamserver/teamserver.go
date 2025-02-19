package teamserver

import (
	"crypto/rsa"
)

type Agent struct {
	AgentId      uint64
	TaskProgress uint64
	PrivateKey   *rsa.PrivateKey
}

type Task struct {
	TaskId uint64
	Task   string
}

type TaskResultIn struct {
	TaskId uint64
	Output string
}

type TaskResultOut struct {
	TaskId uint64
	Task   string
	Output string
}

type ClientError struct {
	Err string
}

func (clientError *ClientError) Error() string {
	return clientError.Err
}

func NewClientError(err string) *ClientError {
	return &ClientError{
		Err: err,
	}
}

type ServerError struct {
	Err string
}

func (serverError *ServerError) Error() string {
	return serverError.Err
}

func NewServerError(err string) *ServerError {
	return &ServerError{
		Err: err,
	}
}

type AgentRepository interface {
	AddAgent() (*Agent, error)
	GetAgent(agentId uint64) (*Agent, error)
	GetAgentTaskProgress(agentId uint64) (uint64, error)
	UpdateAgentTaskProgress(agentId uint64) error
	AgentExists(agentId uint64) (bool, error)
}

type TaskQueueRepository interface {
	TaskQueuePush(task string) error
	GetTasks(agentId uint64) ([]Task, error)
	TaskExists(taskId uint64) (bool, error)
}

// TOOO: Figure out how to translate api model to domain model when they're not one to one
type TaskResultsRepository interface {
	SaveTaskResults(agentId uint64, taskResults []TaskResultIn) error
	GetTaskResult(agentId uint64, taskId uint64) (*TaskResultOut, error)
	GetTaskResults(agentId uint64) ([]TaskResultOut, error)
}
