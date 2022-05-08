package service

import (
	"fmt"

	"github.com/smart7even/golang-do/internal/domain"
)

type TodoRepo interface {
	Create(todo domain.Todo) error
	ReadAll() ([]domain.Todo, error)
	Update(todo domain.Todo) error
	Delete(id int64) error
}

type TodoDoesNotExist struct {
	TodoId int64
}

func (e TodoDoesNotExist) Error() string {
	return fmt.Sprintf("todo with id %v does not exist", e.TodoId)
}

type TodoService struct {
	todoRepo TodoRepo
}

func NewTodoService(todoRepo TodoRepo) *TodoService {
	return &TodoService{
		todoRepo: todoRepo,
	}
}

func (s *TodoService) Create(todo domain.Todo) error {
	return s.todoRepo.Create(todo)
}

func (s *TodoService) ReadAll() ([]domain.Todo, error) {
	return s.todoRepo.ReadAll()
}

func (s *TodoService) Update(todo domain.Todo) error {
	return s.todoRepo.Update(todo)
}

func (s *TodoService) Delete(todoId int64) error {
	return s.todoRepo.Delete(todoId)
}
