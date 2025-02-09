package teamserver

import (
	"crypto/rsa"
	"database/sql"
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

type TaskResult struct {
	Task
	Output string
}

type AgentService interface {
	AddAgent() (*Agent, error)
	GetAgent(agentId uint64) (*Agent, error)
	GetAgentTaskProgress(agentId uint64) (uint64, error)
	UpdateAgentTaskProgress(agentId uint64) error
}

type TaskQueueService interface {
	TaskQueuePush(task string) error // TODO: Use api structs as input
	GetTasks(agentId uint64) ([]Task, error)
}

// TOOO: Figure out how to translate api model to domain model when they're not one to one
type TaskResultsService interface {
	SaveTaskResults(agentId uint64) error
	GetTaskResult(agentId uint64, taskId uint64) (*TaskResult, error)
	GetTaskResults(agentId uint64) ([]TaskResult, error)
}

type Serivces interface {
	NewAgentService(*sql.DB) (AgentService, error)
	NewTaskQueueService(*sql.DB) (TaskQueueService, error)
	NewTaskResultService(*sql.DB) (TaskResultsService, error)
}
