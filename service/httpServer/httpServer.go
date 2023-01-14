package httpserver

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func Run(ctx context.Context, path string, controllers Controllers) error {
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
	log.Debug().Msgf("Config: %#s", string(configFile))
	config := ServerConfig{Router: make(Router)}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, fmt.Errorf("getRouter - error unmarshal: %v ", err)
	}

	log.Debug().Msgf("Config: %#v", config)

	return &config, nil
}

func handlerInit(router Router, controllers Controllers) *gin.Engine {
	handler := gin.New()

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
