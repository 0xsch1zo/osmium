package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

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
		return nil, repositoryErrWrapper(NewRepositoryErrInvalidCredentials())
	}

	expiryTime := time.Now().Add(jwtExpiryTime)
	token, err := tools.GenerateJWT(username, expiryTime, auths.jwtKey)
	if err != nil {
		return nil, err
	}

	return &teamserver.AuthToken{
		Token:      token,
		ExpiryTime: expiryTime,
	}, nil
}

func (auths *AuthorizationService) Authorize(token string) error {
	authorized, err := tools.VerifyJWT(token, auths.jwtKey)
	if err != nil {
		return err
	}

	if !authorized {
		return NewRepositoryErrInvalidCredentials()
	}

	return nil
}

func (auths *AuthorizationService) RefreshToken(token string) (*teamserver.AuthToken, error) {
	authorized, err := tools.VerifyJWT(token, auths.jwtKey)
	if err != nil {
		return nil, err
	}

	if !authorized {
		return nil, NewRepositoryErrInvalidCredentials()
	}

	claims, err := tools.GetJWTClaims(token, auths.jwtKey)
	if err != nil {
		return nil, err
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return nil, NewRepositoryErrTokenNotOld()
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
	passwordHash, err := auths.authorizationRepository.GetPasswordHash(username)
	return passwordHash, repositoryErrWrapper(err)
}

func (auths *AuthorizationService) UsernameExists(username string) (bool, error) {
	exists, err := auths.authorizationRepository.UsernameExists(username)
	return exists, repositoryErrWrapper(err)
}
