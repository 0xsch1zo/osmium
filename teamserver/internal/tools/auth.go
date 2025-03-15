package tools

import (
	"errors"
	"github.com/sentientbottleofwine/osmium/teamserver"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}

func CompareHashToPassword(passwordHash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewJWTWithClaims(claims *teamserver.Claims, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func GenerateJWT(username string, expiryTime time.Time, jwtKey string) (string, error) {
	claims := &teamserver.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}

	return NewJWTWithClaims(claims, jwtKey)
}

func VerifyJWT(tokenStr, jwtKey string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &teamserver.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err == jwt.ErrSignatureInvalid {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if !token.Valid {
		return false, nil
	}

	return true, err
}

func GetJWTClaims(tokenStr, jwtKey string) (*teamserver.Claims, error) {
	var claims teamserver.Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return &claims, nil
}
