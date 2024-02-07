package domain

import "time"

type Reminder struct {
	Id          int       `json:"id"`
	UserId      string    `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RemindAt    time.Time `json:"remindAt"`
	CreatedAt   time.Time `json:"createdAt"`
}
