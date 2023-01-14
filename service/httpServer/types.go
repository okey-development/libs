package httpserver

import (
	"github.com/gin-gonic/gin"
)

type Path string

type Method string

type Controller string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

type Controllers map[Controller]func(c *gin.Context)

type Router map[Path]map[Method]Controller

type ServerConfig struct {
	Port           string `yaml:"Port"`
	MaxHeaderBytes int    `yaml:"MaxHeaderBytes"`
	ReadTimeout    int    `yaml:"ReadTimeout"`
	WriteTimeout   int    `yaml:"WriteTimeout"`
	Router         Router `yaml:"Router"`
}
