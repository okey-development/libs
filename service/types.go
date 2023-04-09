package service

import (
	"github.com/gin-gonic/gin"
)

type path string

type method string

type Controller string

const (
	GET    method = "GET"
	POST   method = "POST"
	PUT    method = "PUT"
	DELETE method = "DELETE"
)

type Controllers map[Controller]func(c *gin.Context)

type Router map[path]map[method]Controller

type ServerConfig struct {
	Port           string `yaml:"Port"`
	MaxHeaderBytes int    `yaml:"MaxHeaderBytes"`
	ReadTimeout    int    `yaml:"ReadTimeout"`
	WriteTimeout   int    `yaml:"WriteTimeout"`
	Router         Router `yaml:"Router"`
}

type Header struct {
	Authorization string
	Language      Lang
}
