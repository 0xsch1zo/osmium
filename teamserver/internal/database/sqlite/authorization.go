package sqlite

import (
	"database/sql"
	"errors"

	"github.com/sentientbottleofwine/osmium/teamserver/service"
)

func (authr *AuthorizationRepository) Register(username, passwordHash string) error {
	query := "INSERT INTO Users (Username, PasswordHash) VALUES(?, ?)"
	_, err := authr.databaseHandle.Exec(query, username, passwordHash)
	return err
}

func (authr *AuthorizationRepository) GetPasswordHash(username string) (string, error) {
	query := "SELECT PasswordHash FROM Users WHERE Username = ?"
	row := authr.databaseHandle.QueryRow(query, username)
	var passwordHash string
	err := row.Scan(&passwordHash)
	if err == sql.ErrNoRows {
		return "", service.NewRepositoryErrInvalidCredentials()
	}

	return passwordHash, err
}

func (authr *AuthorizationRepository) SetSessionToken(username, sessionToken string) error {
	query := "UPDATE Users SET SessionToken = ? WHERE Username = ?"
	_, err := authr.databaseHandle.Exec(query, sessionToken, username)
	return err
}

func (authr *AuthorizationRepository) GetSessionToken(username string) (string, error) {
	query := "SELECT SessionToken FROM Users WHERE Username = ?"
	row := authr.databaseHandle.QueryRow(query, username)

	var sessionToken string
	err := row.Scan(&sessionToken)
	return sessionToken, err
}

func (authr *AuthorizationRepository) UsernameExists(username string) (bool, error) {
	query := "SELECT Username FROM Users WHERE Username = ?"
	row := authr.databaseHandle.QueryRow(query, username)
	var temp string
	err := row.Scan(&temp)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
