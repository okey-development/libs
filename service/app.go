package service

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type config struct {
	Debug        bool   `default:"true"`
	ServerConfig string `default:"./router/router.yaml"`
	DB           dbConfig
}

type Service struct {
	AppName     string
	Controllers Controllers
	Befor       func() error
	After       func() error
}

var appConfig *Service

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
	log.Logger = log.Output(output)
}

func Start(serverConfig *Service) error {

	if serverConfig.Befor != nil {
		if err := serverConfig.Befor(); err != nil {
			return fmt.Errorf("error in Befor function: %s", err.Error())
		}
	}
	appConfig = serverConfig
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config{}

	if err := envconfig.Process(serverConfig.AppName, &cfg); err != nil {
		log.Fatal().Err(err)
	}

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msgf("Config: %#v", cfg)

	if err := initDB(&cfg.DB); err != nil {
		log.Error().Err(err).Msg("Error connect ti data base")
	}

	errWg, errCtx := errgroup.WithContext(ctx)

	errWg.Go(func() error {
		return runHttpServer(errCtx, cfg.ServerConfig, serverConfig.Controllers)
	})

	if serverConfig.After != nil {
		if err := serverConfig.After(); err != nil {
			return fmt.Errorf("error in After function: %s", err.Error())
		}
	}

	if err := errWg.Wait(); err == context.Canceled || err == nil {
		log.Info().Msg("gracefully quit server")
	} else if err != nil {
		log.Error().Err(err).Msg("Error stop service")
	}

	return nil
}
