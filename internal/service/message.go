package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/smart7even/golang-do/internal/domain"
)

type MessageRepo interface {
	Create(message domain.Message, userId string) error
	ReadAll(chatId string, userId string) ([]domain.Message, error)
	Update(message domain.Message, userId string) error
	Delete(id string, userId string) error
}

func NewMessageService(messageRepo MessageRepo, eventsRepo EventRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		eventsRepo:  eventsRepo,
	}
}

type MessageService struct {
	messageRepo MessageRepo
	eventsRepo  EventRepo
}

func (s *MessageService) Create(message domain.Message, userId string) error {
	if message.Text == "" {
		return errors.New("message text is required")
	}

	event := domain.Event{
		Id:          uuid.New().String(),
		Name:        "message_created",
		Description: "message created",
		CreatedAt:   time.Now(),
		Payload: map[string]interface{}{
			"message": message,
		},
	}

	err := s.eventsRepo.Create(event)

	if err != nil {
		return err
	}

	return s.messageRepo.Create(message, userId)
}

func (s *MessageService) ReadAll(chatId string, userId string) ([]domain.Message, error) {
	return s.messageRepo.ReadAll(chatId, userId)
}

func (s *MessageService) Update(message domain.Message, userId string) error {
	return s.messageRepo.Update(message, userId)
}

func (s *MessageService) Delete(id string, userId string) error {
	return s.messageRepo.Delete(id, userId)
}
