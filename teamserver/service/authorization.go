package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (auths *AuthorizationService) Register(username, password string) error {
	err := auths.UsernameExists(username)
	if err == nil {
		return teamserver.NewClientError(fmt.Sprintf(errAlreadyExistsFmt, "username"))
	} else if err.Error() != errInvalidCredentials {
		return err
	}

	if len(username) == 0 || len(password) == 0 {
		return teamserver.NewClientError(fmt.Sprintf(errEmptyString, "username or password"))
	}

	passwordHash, err := tools.HashPassword(password)
	if err != nil {
		return teamserver.NewServerError(err.Error())
	}

	err = auths.authorizationRepository.Register(username, passwordHash)
	return err
}

func (auths *AuthorizationService) Login(username, password string) (*teamserver.AuthToken, error) {
	passwordHash, err := auths.GetPasswordHash(username)
	if err != nil {
		return nil, err
	}

	match, err := tools.CompareHashToPassword(passwordHash, password)
	if err != nil {
		return nil, err
	}

	if !match {
		auths.eventLogService.LogEvent(teamserver.Warn, "Unsuccessful login")
		return nil, teamserver.NewClientError(errInvalidCredentials)
	}

	expiryTime := time.Now().Add(jwtExpiryTime)
	token, err := tools.GenerateJWT(username, expiryTime, auths.jwtKey)
	if err != nil {
		return nil, err
	}

	auths.eventLogService.LogEvent(teamserver.Info, "Successful login")
	return &teamserver.AuthToken{
		Token:      token,
		ExpiryTime: expiryTime,
	}, nil
}

func (auths *AuthorizationService) Authorize(token string) error {
	authorized, err := tools.VerifyJWT(token, auths.jwtKey)
	if err != nil {
		auths.eventLogService.LogEvent(teamserver.Warn, "Unauthorized request was made")
		return teamserver.NewClientError(err.Error())
	}

	if !authorized {
		auths.eventLogService.LogEvent(teamserver.Warn, "Unauthorized request was made")
		return teamserver.NewClientError(errInvalidCredentials)
	}

	return nil
}

func (auths *AuthorizationService) RefreshToken(token string) (*teamserver.AuthToken, error) {
	err := auths.Authorize(token)
	if err != nil {
		return nil, err
	}

	claims, err := tools.GetJWTClaims(token, auths.jwtKey)
	if err != nil {
		return nil, err
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return nil, teamserver.NewClientError(errTokenNotOld)
	}

	expiryTime := time.Now().Add(jwtExpiryTime)
	claims.ExpiresAt = jwt.NewNumericDate(expiryTime)
	tokenRefreshed, err := tools.NewJWTWithClaims(claims, auths.jwtKey)
	return &teamserver.AuthToken{
		Token:      tokenRefreshed,
		ExpiryTime: expiryTime,
	}, err
}

func (auths *AuthorizationService) GetPasswordHash(username string) (string, error) {
	err := auths.UsernameExists(username)
	if err != nil {
		return "", err
	}

	passwordHash, err := auths.authorizationRepository.GetPasswordHash(username)
	return passwordHash, err
}

func (auths *AuthorizationService) UsernameExists(username string) error {
	exists, err := auths.authorizationRepository.UsernameExists(username)
	if err != nil {
		return err
	}

	if !exists {
		return teamserver.NewClientError(errInvalidCredentials)
	}

	return nil
}
