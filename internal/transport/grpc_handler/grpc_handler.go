package grpc_handler

import (
	context "context"
	"log"

	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type TodoGrpcServer struct {
	services service.Services
	UnimplementedTodoServiceServer
}

func NewTodoGrpcServer(services service.Services) *TodoGrpcServer {
	return &TodoGrpcServer{services: services}
}

func (s *TodoGrpcServer) GetTodos(ctx context.Context, in *GetTodosParams) (*Todos, error) {
	log.Printf("Received get todos request")

	todos, err := s.services.Todo.ReadAll()

	if err != nil {
		return nil, err
	}

	responseTodos := make([]*Todo, len(todos))

	for _, todo := range todos {
		responseTodos = append(responseTodos, &Todo{Id: todo.Id, Name: todo.Name, Complete: todo.Complete})
	}

	return &Todos{Todos: responseTodos}, nil
}

func (s *TodoGrpcServer) AddTodo(ctx context.Context, in *AddTodoParams) (*AddTodoResponse, error) {
	log.Printf("Received add todo request")

	todo := in.GetTodo()

	err := s.services.Todo.Create(domain.Todo{Id: todo.Id, Name: todo.Name, Complete: todo.Complete})

	if err != nil {
		return &AddTodoResponse{Added: false}, err
	}

	return &AddTodoResponse{Added: true}, nil
}

func (s *TodoGrpcServer) EditTodo(ctx context.Context, in *EditTodoParams) (*EditTodoResponse, error) {
	log.Printf("Received edit todo request")

	todo := in.GetTodo()

	err := s.services.Todo.Update(domain.Todo{Id: todo.Id, Name: todo.Name, Complete: todo.Complete})

	if err != nil {
		return &EditTodoResponse{Edited: false}, err
	}

	return &EditTodoResponse{Edited: true}, nil
}

func (s *TodoGrpcServer) DeleteTodo(ctx context.Context, in *DeleteTodoParams) (*DeleteTodoResponse, error) {
	log.Printf("Received delete todo request")

	todoId := in.GetTodoId()

	err := s.services.Todo.Delete(todoId)

	if err != nil {
		return &DeleteTodoResponse{Deleted: false}, err
	}

	return &DeleteTodoResponse{Deleted: true}, nil
}
