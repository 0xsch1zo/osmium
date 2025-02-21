package sqlite

import (
	"crypto/rsa"
	"database/sql"
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

func (ar *AgentRepository) AddAgent(rsaPriv *rsa.PrivateKey) (*teamserver.Agent, error) {
	query := "INSERT INTO Agents (AgentId, TaskProgress, PrivateKey) values(NULL, 1, ?);"
	_, err := ar.databaseHandle.Exec(query, tools.PrivRsaToPem(rsaPriv))
	if err != nil {
		return nil, err
	}

	// Get last row in db to get the AgentId of the newly created Agent
	query = "SELECT AgentId FROM Agents ORDER BY AgentId DESC LIMIT 1;" // in sqlite integer primary key will autoicrement as long as null is passed in
	AgentIdSqlRow := ar.databaseHandle.QueryRow(query)

	var AgentId uint64
	err = AgentIdSqlRow.Scan(&AgentId)
	if err != nil {
		return nil, err
	}

	return &teamserver.Agent{
		AgentId:    AgentId,
		PrivateKey: rsaPriv,
	}, nil
}

func (ar *AgentRepository) GetAgent(agentId uint64) (*teamserver.Agent, error) {
	query := "SELECT AgentId, TaskProgress, PrivateKey FROM Agents WHERE AgentId = ?"
	AgentSqlRow := ar.databaseHandle.QueryRow(query, agentId)

	var agent teamserver.Agent
	var agentPrivateKeyPem string
	err := AgentSqlRow.Scan(&agent.AgentId, &agent.TaskProgress, &agentPrivateKeyPem)
	if err == sql.ErrNoRows {
		return nil, service.NewRepositoryErrNotFound(fmt.Sprintf(service.ErrAgentIdNotFoundFmt, agentId))
	} else if err != nil {
		return nil, err
	}

	agent.PrivateKey, err = tools.PemToPrivRsa(agentPrivateKeyPem)
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

func (ar *AgentRepository) AgentExists(agentId uint64) (bool, error) {
	query := "SELECT AgentId FROM Agents WHERE AgentId = ?"
	sqlRow := ar.databaseHandle.QueryRow(query, agentId)

	var temp uint64
	err := sqlRow.Scan(&temp)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (ar *AgentRepository) GetAgentTaskProgress(agentId uint64) (uint64, error) {
	query := "SELECT TaskProgress FROM Agents WHERE AgentId = ?"
	AgentSqlRow := ar.databaseHandle.QueryRow(query, agentId)
	var taskProgress uint64
	err := AgentSqlRow.Scan(&taskProgress)
	if err == sql.ErrNoRows {
		return 0, service.NewRepositoryErrNotFound(fmt.Sprintf(service.ErrAgentIdNotFoundFmt, agentId))
	} else if err != nil {
		return 0, err
	}

	return taskProgress, nil
}

// TODO: Fix this shit
func (ar *AgentRepository) UpdateAgentTaskProgress(agentId uint64) error {
	query := "UPDATE Agents SET TaskProgress = (SELECT MAX(TaskId) FROM TaskQueue)"
	_, err := ar.databaseHandle.Exec(query)
	return err
}

func (ar *AgentRepository) ListAgents() ([]teamserver.AgentView, error) {
	query := "SELECT AgentId, Task FROM Agents INNER JOIN TaskQueue ON Agents.TaskProgress = TaskQueue.TaskId"
	sqlRow, err := ar.databaseHandle.Query(query)
	if err != nil {
		return nil, err
	}

	var AgentViews []teamserver.AgentView
	for sqlRow.Next() {
		var agent teamserver.AgentView
		err = sqlRow.Scan(&agent.AgentId, &agent.Task)
		if err != nil {
			return nil, err
		}
		AgentViews = append(AgentViews, agent)
	}

	return AgentViews, nil
}
