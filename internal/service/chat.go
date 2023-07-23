package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/smart7even/golang-do/internal/domain"
)

type ChatRepo interface {
	Create(chat domain.Chat) error
	ReadAll() ([]domain.Chat, error)
	Update(chat domain.Chat) error
	Delete(id string) error
	Join(chatId string, userId string) error
}

type ChatService struct {
	chatRepo ChatRepo
}

type ChatDoesNotExist struct {
	Id int64
}

func (e ChatDoesNotExist) Error() string {
	return fmt.Sprintf("chat with id %v does not exist", e.Id)
}

func NewChatService(chatRepo ChatRepo) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (s *ChatService) Create(chat domain.Chat) error {
	if chat.Name == "" {
		return errors.New("chat name is required")
	}

	uuid := uuid.New()
	chat.Id = uuid.String()

	return s.chatRepo.Create(chat)
}

func (s *ChatService) ReadAll() ([]domain.Chat, error) {
	return s.chatRepo.ReadAll()
}

func (s *ChatService) Update(chat domain.Chat) error {
	return s.chatRepo.Update(chat)
}

func (s *ChatService) Delete(id string) error {
	return s.chatRepo.Delete(id)
}

func (s *ChatService) Join(chatId string, userId string) error {
	return s.chatRepo.Join(chatId, userId)
}
