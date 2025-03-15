package teamserver

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
)

type Agent struct {
	AgentId      uint64
	TaskProgress uint64
	PrivateKey   *rsa.PrivateKey
}

type AgentView struct {
	AgentId uint64
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

type Claims struct {
	Username string
	jwt.RegisteredClaims
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
