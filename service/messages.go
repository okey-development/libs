package service

import "encoding/json"

type FirebaseOperation string

const (
	FIREBASE_SUBSCRIBE   FirebaseOperation = "FIREBASE_SUBSCRIBE"
	FIREBASE_UNSUBSCRIBE FirebaseOperation = "FIREBASE_UNSUBSCRIBE"
)

type TypeMessages string

const (
	FIREBASE_PUSH TypeMessages = "FIREBASE_PUSH"
	MAIL          TypeMessages = "MAIL"
)

type SubTypeMessages string

const (
	PUSH_GET_DEVICE     SubTypeMessages = "PUSH_GET_DEVICE"
	PUSH_ORDER          SubTypeMessages = "PUSH_ORDER"
	PUSH_PRO_SCREEN     SubTypeMessages = "PUSH_PRO_SCREEN"
	PUSH_PROGRESS       SubTypeMessages = "PUSH_PROGRESS"
	PUSH_TRANSACTION    SubTypeMessages = "PUSH_TRANSACTION"
	MAIL_RESET_PASSWORD SubTypeMessages = "MAIL_RESET_PASSWORD"
	MAIL_EMAIL_VERIFITY SubTypeMessages = "MAIL_EMAIL_VERIFITY"
	MAIL_ORDER          SubTypeMessages = "MAIL_ORDER"
)

type Message struct {
	TypeMessages    TypeMessages
	SubTypeMessages SubTypeMessages
	Message         string
	Delay           int64
	Lang            Lang
}

func SendMessage(message Message) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return Errorf(err.Error())
	}

	if err := SendRedisMessage("MESSAGES", jsonData); err != nil {
		return Errorf(err.Error())
	}
	return nil
}
