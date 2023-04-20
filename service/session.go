package service

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Object string

type Action string

type SessionConfig struct {
	Accesses map[Object][]Action
	Rights   map[Object][]Action
	Handle   func(ctx *Context) *Response
}

type Context struct {
	Gin   *gin.Context
	User  *User
	Local *local
}

func BaseHandle(conf *SessionConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Language")
		header := Header{}
		if err := c.ShouldBindHeader(&header); err != nil {
			Error(err)
			NewResponce(http.StatusBadRequest, GetLocal(IncorrectParameter, header.Language), nil).Send(c)
			return
		}

		lang := Lang(strings.ToUpper(string(header.Language)))
		if lang == EMTY {
			lang = EN
		}
		local := newLocal(lang)

		var user *User
		if len(conf.Rights) > 0 || len(conf.Accesses) > 0 {
			var err error
			accessToken := strings.Replace(header.Authorization, "Bearer ", "", -1)
			user, err = newToken(signingKeyAccess).Validate(accessToken)
			if err != nil {
				Error(err)
				NewResponce(GetCodeError(err), local.ParseError(err), nil).Send(c)
				return
			}

			if !user.CheckAccesses(conf.Accesses) {
				NewResponce(http.StatusForbidden, local.ParseError(NewError(ForbbidenAccess).Error(ForbbidenAccess)), nil).Send(c)
				return
			}

			if !user.CheckRights(conf.Rights) {
				NewResponce(http.StatusForbidden, local.ParseError(NewError(ForbbidenRights).Error(ForbbidenRights)), nil).Send(c)
				return
			}
		}

		conf.Handle(&Context{
			Gin:   c,
			User:  user,
			Local: local,
		}).Send(c)
	}
}
