package service

import (
	"encoding/json"
	"fmt"
	"runtime"

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

type ErrorConstructor struct {
	key string
}

func NewError(key string) *ErrorConstructor {
	return &ErrorConstructor{key: key}
}

func (e *ErrorConstructor) Error(err error) error {
	var details = GetErrorDetails(err).Error()
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	if n != 0 {
		// Извлекаем информацию о вызывающей функции
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()

		details = fmt.Sprintf("%s: %s", frame.Function, details)
	}

	errorBody, _ := json.Marshal(&map[string]string{
		"key":     e.key,
		"details": details,
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

func GetErrorDetails(err error) error {
	errorBody := make(map[string]string)
	_ = json.Unmarshal([]byte(err.Error()), errorBody)
	if _, ok := errorBody["details"]; ok {
		return fmt.Errorf(errorBody["details"])
	}
	return err
}

func Errorf(format string, arg ...interface{}) error {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return fmt.Errorf(format, arg...)
	}

	// Извлекаем информацию о вызывающей функции
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	return fmt.Errorf("%s: %s", frame.Function, fmt.Sprintf(format, arg...))
}

func GetCallerName() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	if n == 0 {
		return ""
	}

	// Извлекаем информацию о вызывающей функции
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	// Получаем имя вызывающей функции
	return frame.Function
}
