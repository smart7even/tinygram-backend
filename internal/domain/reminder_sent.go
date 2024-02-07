package domain

import "time"

type ReminderSent struct {
	Id         int       `json:"id"`
	ReminderId int       `json:"reminderId"`
	DeviceId   int       `json:"deviceId"`
	UserId     string    `json:"userId"`
	SentAt     time.Time `json:"sentAt"`
}
