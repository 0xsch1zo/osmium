package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (auths *AuthorizationService) Register(username, password string) error {
	err := auths.UsernameExists(username)
	if err == nil {
		return teamserver.NewClientError(fmt.Sprintf(errAlreadyExistsFmt, "username"), http.StatusConflict)
	} else if err.Error() != errInvalidCredentials {
		return err
	}

	if len(username) == 0 || len(password) == 0 {
		return teamserver.NewClientError(fmt.Sprintf(errEmptyString, "username or password"), http.StatusBadRequest)
	}

	passwordHash, err := tools.HashPassword(password)
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return err
	}

	err = auths.authorizationRepository.Register(username, passwordHash)
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return err
	}
	return nil
}

func (auths *AuthorizationService) Login(username, password string) (*teamserver.AuthToken, error) {
	passwordHash, err := auths.GetPasswordHash(username)
	if err != nil {
		return nil, err
	}

	match, err := tools.CompareHashToPassword(passwordHash, password)
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return nil, err
	}

	if !match {
		auths.eventLogService.LogEvent(teamserver.Warn, "Unsuccessful login")
		return nil, teamserver.NewClientError(errInvalidCredentials, http.StatusUnauthorized)
	}

	expiryTime := time.Now().Add(jwtExpiryTime)
	token, err := tools.GenerateJWT(username, expiryTime, auths.jwtKey)
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
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
		return teamserver.NewClientError(err.Error(), http.StatusUnauthorized)
	}

	if !authorized {
		auths.eventLogService.LogEvent(teamserver.Warn, "Unauthorized request was made")
		return teamserver.NewClientError(errInvalidCredentials, http.StatusUnauthorized)
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
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return nil, err
	}

	if time.Until(claims.ExpiresAt.Time) > jwtRefreshWindow {
		return nil, teamserver.NewClientError(errTokenNotOld, http.StatusUnauthorized)
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
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return "", err
	}
	return passwordHash, nil
}

func (auths *AuthorizationService) UsernameExists(username string) error {
	exists, err := auths.authorizationRepository.UsernameExists(username)
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return err
	}

	if !exists {
		return teamserver.NewClientError(errInvalidCredentials, http.StatusUnauthorized)
	}

	return nil
}

func (auths *AuthorizationService) GetRefreshTime(token string) (string, error) {
	err := auths.Authorize(token)
	if err != nil {
		return "", err
	}

	claims, err := tools.GetJWTClaims(token, auths.jwtKey)
	if err != nil {
		ServiceServerErrHandler(err, authorizationServiceStr, auths.eventLogService)
		return "", err
	}

	refTime := claims.ExpiresAt.Time.Add(-jwtRefreshWindow)
	return refTime.Format(time.RFC3339), nil
}
