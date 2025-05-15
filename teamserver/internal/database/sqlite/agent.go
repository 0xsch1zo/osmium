package sqlite

import (
	"crypto/rsa"
	"database/sql"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (ar *AgentRepository) AddAgent(rsaPriv *rsa.PrivateKey, agentInfo teamserver.AgentRegisterInfo) (*teamserver.Agent, error) {
	// race condition
	query := "INSERT INTO Agents (AgentId, PrivateKey, Hostname, Username, LastCallback) values(NULL, ?, ?, ?, 0);"
	_, err := ar.databaseHandle.Exec(query, tools.PrivRsaToPem(rsaPriv), agentInfo.Hostname, agentInfo.Username)
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
		AgentInfo: teamserver.AgentInfo{
			Hostname:     agentInfo.Hostname,
			Username:     agentInfo.Username,
			LastCallback: time.Unix(0, 0),
		},
	}, nil
}

func (ar *AgentRepository) GetAgent(agentId uint64) (*teamserver.Agent, error) {
	query := "SELECT AgentId, PrivateKey, Hostname, Username, LastCallback FROM Agents WHERE AgentId = ?"
	AgentSqlRow := ar.databaseHandle.QueryRow(query, agentId)

	var agent teamserver.Agent
	var agentPrivateKeyPem string
	var lastCallbackUnix int64
	err := AgentSqlRow.Scan(
		&agent.AgentId,
		&agentPrivateKeyPem,
		&agent.AgentInfo.Hostname,
		&agent.AgentInfo.Username,
		&lastCallbackUnix,
	)
	if err != nil {
		return nil, err
	}

	agent.AgentInfo.LastCallback = time.Unix(lastCallbackUnix, 0)

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

func (ar *AgentRepository) ListAgents() ([]teamserver.AgentView, error) {
	query := "SELECT AgentId, Hostname, Username, LastCallback FROM Agents"
	sqlRow, err := ar.databaseHandle.Query(query)
	if err != nil {
		return nil, err
	}

	var AgentViews []teamserver.AgentView
	for sqlRow.Next() {
		var lastCallbackUnix int64
		var agent teamserver.AgentView
		err = sqlRow.Scan(
			&agent.AgentId,
			&agent.AgentInfo.Hostname,
			&agent.AgentInfo.Username,
			&lastCallbackUnix,
		)
		if err != nil {
			return nil, err
		}

		agent.AgentInfo.LastCallback = time.Unix(lastCallbackUnix, 0)

		AgentViews = append(AgentViews, agent)
	}

	return AgentViews, nil
}

func (ar *AgentRepository) UpdateLastCallbackTime(agentId uint64) error {
	query := "UPDATE Agents SET LastCallback = ? WHERE AgentId = ?"
	_, err := ar.databaseHandle.Exec(query, time.Now().Unix(), agentId)
	return err
}
