package service

import (
	"fmt"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (auths *AuthorizationService) Register(username, password string) error {
	exists, err := auths.UsernameExists(username)
	if err != nil {
		return err
	}

	if exists {
		return teamserver.NewClientError(fmt.Sprintf(ErrAlreadyExistsFmt, "username"))
	}

	if len(username) == 0 || len(password) == 0 {
		return teamserver.NewClientError(fmt.Sprintf(ErrEmptyString, "username or password"))
	}

	passwordHash, err := tools.HashPassword(password)
	if err != nil {
		return teamserver.NewServerError(err.Error())
	}

	err = auths.authorizationRepository.Register(username, passwordHash)
	return repositoryErrWrapper(err)
}

func (auths *AuthorizationService) Login(username, password string) (string, error) {
	passwordHash, err := auths.GetPasswordHash(username)
	if err != nil {
		return "", err
	}

	match, err := tools.CompareHashToPassword(passwordHash, password)
	if err != nil {
		return "", err
	}

	if !match {
		return "", repositoryErrWrapper(NewRepositoryErrInvalidCredentials())
	}

	sessionToken, err := tools.GenerateToken(TokenSize)
	if err != nil {
		return "", err
	}

	err = auths.SetSessionToken(username, sessionToken)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (auths *AuthorizationService) Authorize(username, sessionToken string) error {
	validSessionToken, err := auths.GetSessionToken(username)
	if err != nil {
		return err
	}

	if sessionToken != validSessionToken {
		return NewRepositoryErrInvalidCredentials()
	}

	return nil
}

func (auths *AuthorizationService) GetPasswordHash(username string) (string, error) {
	passwordHash, err := auths.authorizationRepository.GetPasswordHash(username)
	return passwordHash, repositoryErrWrapper(err)
}

func (auths *AuthorizationService) SetSessionToken(username, sessionToken string) error {
	return repositoryErrWrapper(auths.authorizationRepository.SetSessionToken(username, sessionToken))
}

func (auths *AuthorizationService) GetSessionToken(username string) (string, error) {
	sessionToken, err := auths.authorizationRepository.GetSessionToken(username)
	return sessionToken, repositoryErrWrapper(err)
}

func (auths *AuthorizationService) UsernameExists(username string) (bool, error) {
	exists, err := auths.authorizationRepository.UsernameExists(username)
	return exists, repositoryErrWrapper(err)
}
