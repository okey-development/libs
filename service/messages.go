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
	PUSH_AFTER_SIGNUP             SubTypeMessages = "PUSH_AFTER_SIGNUP"
	PUSH_CHANGE_ORDER_STATUS      SubTypeMessages = "PUSH_CHANGE_ORDER_STATUS"
	PUSH_AFTER_CHANGE_GROUPS      SubTypeMessages = "PUSH_AFTER_CHANGE_GROUPS"
	PUSH_AFTER_VISIT_PP           SubTypeMessages = "PUSH_AFTER_VISIT_PP"
	PUSH_AFTER_SIGNUP_PARTHERS    SubTypeMessages = "PUSH_AFTER_SIGNUP_PARTHERS"
	PUSH_AFTER_RECEIPT            SubTypeMessages = "PUSH_AFTER_RECEIPT"
	PUSH_CHANGE_WITHDRAWAL_STATUS SubTypeMessages = "PUSH_CHANGE_WITHDRAWAL_STATUS"
	MAIL_RESET_PASSWORD           SubTypeMessages = "MAIL_RESET_PASSWORD"
	MAIL_EMAIL_VERIFITY           SubTypeMessages = "MAIL_EMAIL_VERIFITY"
	MAIL_CHANGE_ORDER_STATUS      SubTypeMessages = "MAIL_CHANGE_ORDER_STATUS"
)

type Message struct {
	UserId          int64
	TypeMessages    TypeMessages
	SubTypeMessages SubTypeMessages
	Message         string
	Delay           int64
	Details         map[string]string
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

func DeleteMessages(userId int64, tm TypeMessages, stm SubTypeMessages) {
	message := Message{
		UserId:          userId,
		TypeMessages:    tm,
		SubTypeMessages: stm,
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
