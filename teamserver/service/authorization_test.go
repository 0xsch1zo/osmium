package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/sentientbottleofwine/osmium/teamserver"
	"github.com/sentientbottleofwine/osmium/teamserver/internal/tools"
)

const (
	username = "user"
	password = "testing"
)

func getAuthToken(testedServices *testedServices) (string, error) {
	err := testedServices.authorizationService.Register(username, password)
	if err != nil {
		return "", err
	}

	token, err := testedServices.authorizationService.Login(username, password)
	return token.Token, err
}

func TestRegister(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	err = testedServices.authorizationService.Register(username, password)
	if err != nil {
		t.Fatal(err)
	}

	err = testedServices.authorizationService.UsernameExists(username)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(testedServices)
	}

	token, err := getAuthToken(testedServices)
	if err != nil {
		t.Fatal(err)
	}

	err = testedServices.authorizationService.Authorize(token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRefreshToken(t *testing.T) {
	testedServices, err := newTestedServices()
	if err != nil {
		t.Fatal(err)
	}

	expiryTime := time.Now().Add(35 * time.Second)
	token, err := tools.GenerateJWT(username, expiryTime, testingJwtKey)
	if err != nil {
		t.Fatal(err)
	}

	target := &teamserver.ClientError{}
	_, err = testedServices.authorizationService.RefreshToken(token)
	if err == nil {
		t.Fatal("An error wasn't thrown for a newly refreshed token")
	} else if !errors.As(err, &target) {
		t.Fatal(err)
	}

	nearExpiryDuration := time.Until(expiryTime.Add(-30 * time.Second))
	time.Sleep(nearExpiryDuration)

	_, err = testedServices.authorizationService.RefreshToken(token)
	if err != nil {
		t.Fatal(err)
	}
}
