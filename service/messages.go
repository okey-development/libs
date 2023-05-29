package service

import "encoding/json"

type MessagesOperation string

const (
	SEND_MESSAGE    MessagesOperation = "SEND_MESSAGE"
	DELETE_MESSAGES MessagesOperation = "DELETE_MESSAGES"
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
	UserId          int64
	TypeMessages    TypeMessages
	SubTypeMessages SubTypeMessages
	Message         string
	Delay           int64
	Lang            Lang
	TypeOperation   MessagesOperation
}

func SendMessage(message Message) {
	message.TypeOperation = SEND_MESSAGE
	go func() {
		jsonData, err := json.Marshal(message)
		if err != nil {
			Error(err)
			return
		}

		if err := SendRedisMessage("MESSAGES", jsonData); err != nil {
			Error(err)
			return
		}
	}()

}

func DeleteMessages(userId int64, typeMessages TypeMessages, subTypeMessages SubTypeMessages) {
	message := Message{
		UserId:          userId,
		TypeMessages:    typeMessages,
		SubTypeMessages: subTypeMessages,
		TypeOperation:   DELETE_MESSAGES,
	}
	go func() {
		jsonData, err := json.Marshal(message)
		if err != nil {
			Error(err)
			return
		}

		if err := SendRedisMessage("MESSAGES", jsonData); err != nil {
			Error(err)
			return
		}
	}()

}
