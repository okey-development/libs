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
	user_id, err := t.validate(token)
	if err != nil {
		return nil, NewError(AuthenticationError).Error(err.Error())
	}
	if user_id == -1 {
		return nil, NewError(TokenHasExpired).Error(TokenHasExpired)
	}

	user, err := GetUserByID(user_id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, NewError(AccountDontExist).Error(AccountDontExist)
	}

	return user, err
}

func (t *token) validate(tokenValidating string) (int64, error) {

	var tokenClaims tokenClaims
	token, err := jwt.ParseWithClaims(tokenValidating, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return []byte(t.signingKey), nil
	})

	if err != nil {
		return -1, err
	}

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return -1, fmt.Errorf("token expired")
			}
		}
	}

	return tokenClaims.Id, nil
}

type TokenType string

type tokenClaims struct {
	Id         int64 `json:"id"`
	Authorized bool  `json:"authorized"`
	Exp        int64 `json:"exp"`
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
