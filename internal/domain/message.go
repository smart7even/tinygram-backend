package domain

import (
	"time"
)

type Message struct {
	Id     string    `json:"id"`
	UserId string    `json:"userId"`
	ChatId string    `json:"chatId"`
	SentAt time.Time `json:"sentAt"`
	Text   string    `json:"text"`
}
