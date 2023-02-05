package service

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func Debug(str string, args ...any) {
	log.Debug().Msgf(str, args...)
}

func Info(str string, args ...any) {
	log.Info().Msgf(str, args...)
	insertToDb(appConfig.AppName, fmt.Sprintf(str, args...), INFO)
}

func Error(err error) {
	log.Error().Err(err)
	insertToDb(appConfig.AppName, err.Error(), ERROR)
}

func Warn(str string, args ...any) {
	log.Warn().Msgf(str, args...)
	insertToDb(appConfig.AppName, fmt.Sprintf(str, args...), WARN)
}

func insertToDb(appName, message string, typeLog typeLog) {
	if _, err := Exec(`INSERT INTO audit.logs (service, type, message) VALUES ((select id from admin.apps where name = $1 ), $2, $3);`, appName, typeLog, message); err != nil {
		log.Error().Err(err)
	}
}

type typeLog int

const (
	ERROR typeLog = 0
	WARN  typeLog = 1
	INFO  typeLog = 2
	DEBUG typeLog = 3
)
