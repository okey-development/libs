package service

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Object string

type Action string

type SessionConfig struct {
	Accesses map[Object]Action
	Rights   map[Object]Action
	Handle   func(ctx *Context) *Response
}

type Context struct {
	Gin   *gin.Context
	User  *User
	Local *local
}

func BaseHandle(conf *SessionConfig) func(c *gin.Context) {
	return func(c *gin.Context) {

		header := Header{}
		if err := c.ShouldBindHeader(&header); err != nil {
			Error(err)
			NewResponce(http.StatusBadRequest, GetLocal(IncorrectParameter, header.Language), nil).Send(c)
			return
		}

		local := newLocal(header.Language)

		var user *User
		accessToken := strings.Replace(header.Authorization, "Bearer ", "", -1)
		user, err := newToken(signingKeyAccess).Validate(accessToken)
		if err != nil {
			Error(err)
			NewResponce(http.StatusUnauthorized, local.ParseError(err), nil)
			return
		}

		if !user.CheckAccesses(conf.Accesses) {
			NewResponce(http.StatusForbidden, ForbbidenAccess, nil)
			return
		}

		if !user.CheckRights(conf.Rights) {
			NewResponce(http.StatusForbidden, ForbbidenRights, nil)
			return
		}

		conf.Handle(&Context{
			Gin:   c,
			User:  user,
			Local: local,
		}).Send(c)
	}
}
