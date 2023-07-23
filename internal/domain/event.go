package domain

import "time"

type Event struct {
	Id          string
	Name        string
	Description string
	CreatedAt   time.Time
	Payload     map[string]interface{}
}
