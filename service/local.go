package service

import "sync"

const (
	UnknownError        = "UnknownError"
	AuthenticationError = "AuthenticationError"
	ForbbidenAccess     = "ForbbidenAccess"
	ForbbidenRights     = "ForbbidenRights"
	IncorrectParameter  = "IncorrectParameter"
	TokenHasExpired     = "TokenHasExpired"
	AccountDontExist    = "AccountDontExist"
)

type Lang string

const (
	EMTY Lang = ""
	EN   Lang = "EN"
	RU   Lang = "RU"
)

var mu sync.RWMutex

var localities = map[string]map[Lang]string{
	IncorrectParameter: {
		EMTY: "Incorrect parameter",
		EN:   "Incorrect parameter",
		RU:   "Некорректный параметр",
	},
	AuthenticationError: {
		EMTY: "Authentication error. Please, login",
		EN:   "Authentication error.Please, login",
		RU:   "Ошибка аутентификации. Пожалуйста, ввойдите в свой аккаунт",
	},
	ForbbidenAccess: {
		EMTY: "You do not have access to this resource. Try logging in with a different user",
		EN:   "You do not have access to this resource. Try logging in with a different user",
		RU:   "У вас нет доступа к этому ресурсу. Попробуйте войти под другим пользователем",
	},
	ForbbidenRights: {
		EMTY: "You do not have permission for this action. Try logging in with a different user",
		EN:   "You do not have permission for this action. Try logging in with a different user",
		RU:   "У вас нет прав для данного действия. Попробуйте войти под другим пользователем",
	},
	UnknownError: {
		EMTY: "An unknown error has occurred.Try again later",
		EN:   "An unknown error has occurred.Try again later",
		RU:   "Произошла неизвестная ошибка. Повторите попытку позже.",
	},
	TokenHasExpired: {
		EMTY: "token has expired",
		EN:   "token has expired",
		RU:   "срок действия токена истек",
	},
	AccountDontExist: {
		EMTY: "We couldn’t find your account. Please try again",
		EN:   "We couldn’t find your account. Please try again",
		RU:   "Мы не смогли найти Ваш аккаунт. Пожалуйста, попробуйте снова",
	},
}

func GetLocal(text string, lang Lang) string {
	mu.RLock()
	defer mu.RUnlock()
	local, ok := localities[text][lang]
	if ok {
		return local
	}
	return localities[text][EN]
}

func SetLocal(newLocalities map[string]map[Lang]string) {
	mu.RLock()
	defer mu.RUnlock()

	for locality := range newLocalities {
		localities[locality] = newLocalities[locality]
	}
}
