package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var signingKeyAccess string

func SetSigningKeyAccess(newSigningKeyAccess string) {
	signingKeyAccess = newSigningKeyAccess
}

type token struct {
	signingKey string
}

func newToken(signingKey string) *token {
	return &token{signingKey: signingKey}
}

func (t *token) Validate(token string) (*User, error) {
	user_uuid, err := t.validate(token)
	if err != nil {
		return nil, NewError(AuthenticationError).Error(err.Error())
	}
	if user_uuid == "" {
		return nil, NewError(TokenHasExpired).Error(TokenHasExpired)
	}

	user, err := getUserByUUID(user_uuid)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, NewError(AccountDontExist).Error(AccountDontExist)
	}

	return user, err
}

func (t *token) validate(tokenValidating string) (string, error) {

	var tokenClaims tokenClaims
	token, err := jwt.ParseWithClaims(tokenValidating, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return []byte(t.signingKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return "", fmt.Errorf("token expired")
			}
		}
	}

	return tokenClaims.UUID, nil
}

type TokenType string

type tokenClaims struct {
	UUID       string `json:"uuid"`
	Authorized bool   `json:"authorized"`
	Exp        int64  `json:"exp"`
}

func (token *tokenClaims) Valid() error {
	if !token.Authorized {
		return fmt.Errorf("Unauthorized")
	}
	if time.Now().After(time.Unix(token.Exp, 0)) {
		return fmt.Errorf("Unauthorized")
	}
	return nil
}

const (
	ttlUUID = 24 * 7 * time.Hour
)

func getUserByUUID(uuid string) (*User, error) {
	userId, _ := GetStringInt64(uuid)
	if userId == 0 {
		user, err := GetUserByUUID(uuid)
		if err != nil {
			return nil, err
		}

		if user != nil {
			err = SetInt64Value(uuid, user.Id, ttlUUID)
			if err != nil {
				Error(Errorf(err.Error()))
			}
		}

		return user, nil
	}
	return GetUserByID(userId)
}
