package service

import (
	"context"
	"fmt"
	db "libs/service/DB"
	httpserver "libs/service/httpServer"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Debug        bool   `default:"true"`
	ServerConfig string `default:"./router/router.yaml"`
	DB           db.Config
}

type Service struct {
	AppName     string
	Controllers httpserver.Controllers
	Befor       func() error
	After       func() error
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: false}
	log.Logger = log.Output(output)
}

func Start(config *Service) error {

	if config.Befor != nil {
		if err := config.Befor(); err != nil {
			return fmt.Errorf("error in Befor function: %s", err.Error())
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := Config{}

	if err := envconfig.Process(config.AppName, &cfg); err != nil {
		log.Fatal().Err(err)
	}

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Debug().Msgf("Config: %#v", cfg)

	if err := db.Init(&cfg.DB); err != nil {
		log.Error().Err(err).Msg("Error connect ti data base")
	}

	errWg, errCtx := errgroup.WithContext(ctx)

	errWg.Go(func() error {
		return httpserver.Run(errCtx, cfg.ServerConfig, config.Controllers)
	})

	if config.After != nil {
		if err := config.After(); err != nil {
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
