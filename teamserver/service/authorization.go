package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

func (auths *AuthorizationService) Register(username, password string) error {
	err := auths.UsernameExists(username)
	if err == errors.New(errInvalidCredentials) {
		return teamserver.NewClientError(fmt.Sprintf(errAlreadyExistsFmt, "username"))
	} else if err != nil {
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
		return "", teamserver.NewClientError(errInvalidCredentials)
	}

	expiryTime := time.Now().Add(jwtExpiryTime)
	token, err := tools.GenerateJWT(username, expiryTime, auths.jwtKey)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (auths *AuthorizationService) Authorize(token string) error {
	authorized, err := tools.VerifyJWT(token, auths.jwtKey)
	if err != nil {
		return teamserver.NewClientError(err.Error())
	}

	if !authorized {
		return teamserver.NewClientError(errInvalidCredentials)
	}

	return nil
}

func (auths *AuthorizationService) RefreshToken(token string) (string, error) {
	err := auths.Authorize(token)
	if err != nil {
		return "", err
	}

	claims, err := tools.GetJWTClaims(token, auths.jwtKey)
	if err != nil {
		return "", err
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return "", teamserver.NewClientError(errTokenNotOld)
	}

	expiryTime := time.Now().Add(jwtExpiryTime)
	claims.ExpiresAt = jwt.NewNumericDate(expiryTime)
	tokenRefreshed, err := tools.NewJWTWithClaims(claims, auths.jwtKey)
	return tokenRefreshed, nil
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
