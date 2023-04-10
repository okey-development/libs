package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func runHttpServer(ctx context.Context, path string, controllers Controllers) error {
	config, err := getConfig(path)
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handlerInit(config.Router, controllers),
		MaxHeaderBytes: config.MaxHeaderBytes,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
	}

	go func() {
		<-ctx.Done()
		httpServer.Shutdown(context.Background())
	}()

	return httpServer.ListenAndServe()
}

func getConfig(path string) (*ServerConfig, error) {

	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("getRouter - error read router file: %v ", err)
	}
	config := ServerConfig{Router: make(Router)}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("getRouter - error unmarshal: %v ", err)
	}

	return &config, nil
}

func handlerInit(router Router, controllers Controllers) *gin.Engine {
	handler := gin.Default()
	handler.Static("/source", "./source")
	handler.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
	})
	for route := range router {
		for method := range router[route] {
			if controllers[router[route][method]] == nil {
				log.Error().Msgf("no handler for %s:%s ", method, route)
				continue
			}

			switch method {
			case GET:
				handler.GET(string(route), controllers[router[route][method]])
			case POST:
				handler.POST(string(route), controllers[router[route][method]])
			case PUT:
				handler.PUT(string(route), controllers[router[route][method]])
			case DELETE:
				handler.DELETE(string(route), controllers[router[route][method]])
			default:
				log.Error().Msgf("unknown method %s for route %s and controller %s ", method, route, router[route][method])
			}
		}
	}

	return handler
}

// структура ответа
type Response struct {
	code    int
	message string
	data    interface{}
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewResponce(code int, message string, data interface{}) *Response {
	return &Response{code, message, data}
}

func (resp *Response) Send(c *gin.Context) {
	if resp == nil {
		return
	}
	if resp.code == http.StatusOK && resp.data != nil {
		c.JSON(http.StatusOK, resp.data)
		return
	}
	// блокирует выполнение последующих обработчиков и выводит в ответе сообщение в Json и статус
	c.AbortWithStatusJSON(resp.code, response{resp.code, resp.message})
}

var muhttpCodes sync.RWMutex

var httpCodes = map[string]int{
	IncorrectParameter:  http.StatusBadRequest,
	AuthenticationError: http.StatusUnauthorized,
	ForbbidenAccess:     http.StatusForbidden,
	ForbbidenRights:     http.StatusForbidden,
	UnknownError:        http.StatusInternalServerError,
	TokenHasExpired:     http.StatusUnauthorized,
	AccountDontExist:    http.StatusUnauthorized,
}

func GetCode(text string) int {
	muhttpCodes.RLock()
	defer muhttpCodes.RUnlock()
	code, ok := httpCodes[text]
	if ok {
		return code
	}
	return httpCodes[UnknownError]
}
