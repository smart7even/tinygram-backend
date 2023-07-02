package service

import "github.com/smart7even/golang-do/internal/domain"

type EventRepo interface {
	Create(event domain.Event) error
	ReadAll() ([]domain.Event, error)
}

type EventService struct {
	eventRepo EventRepo
}

func NewEventService(eventRepo EventRepo) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

func (s *EventService) Create(event domain.Event) error {
	return s.eventRepo.Create(event)
}

func (s *EventService) ReadAll() ([]domain.Event, error) {
	return s.eventRepo.ReadAll()
}
