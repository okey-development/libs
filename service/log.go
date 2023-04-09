package service

import (
	"encoding/json"
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

type local struct {
	lang Lang
}

func newLocal(lang Lang) *local {
	return &local{lang: lang}
}

func (local *local) ParseError(err error) string {
	return GetLocal(GetErrorKey(err), local.lang)
}

func (local *local) GetLang() Lang {
	return local.lang
}

func NewError(key, details string, arg ...interface{}) error {
	errorBody, _ := json.Marshal(&map[string]string{
		"key":     key,
		"details": fmt.Sprintf(details, arg...),
	})
	return fmt.Errorf(string(errorBody))
}

func GetErrorKey(err error) string {
	errorBody := make(map[string]string)
	_ = json.Unmarshal([]byte(err.Error()), errorBody)
	if _, ok := errorBody["key"]; ok {
		return errorBody["key"]
	}
	return UnknownError
}

func GetErrorDetails(err error) string {
	errorBody := make(map[string]string)
	_ = json.Unmarshal([]byte(err.Error()), errorBody)
	if _, ok := errorBody["details"]; ok {
		return errorBody["details"]
	}
	return err.Error()
}
