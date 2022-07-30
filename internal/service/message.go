package service

import (
	"errors"

	"github.com/smart7even/golang-do/internal/domain"
)

type MessageRepo interface {
	Create(message domain.Message, userId string) error
	ReadAll(chatId string, userId string) ([]domain.Message, error)
	Update(message domain.Message, userId string) error
	Delete(id string, userId string) error
}

func NewMessageService(messageRepo MessageRepo) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

type MessageService struct {
	messageRepo MessageRepo
}

func (s *MessageService) Create(message domain.Message, userId string) error {
	if message.Text == "" {
		return errors.New("message text is required")
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
