package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/smart7even/golang-do/internal/transport/grpc_handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "127.0.0.1:8081", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTodoServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetTodos(ctx, &pb.GetTodosParams{})
	if err != nil {
		log.Fatalf("could not get todos: %v", err)
	}
	log.Printf("Todos: %v\n", r.GetTodos())

	lastTodo := r.Todos[len(r.Todos)-1]
	fmt.Printf("Last todo: %v\n", lastTodo)

	addTodoResponse, addTodoErr := c.AddTodo(ctx, &pb.AddTodoParams{
		Todo: &pb.Todo{Name: fmt.Sprintf("New todo added at %s\n", time.Now()), Complete: false},
	})

	if addTodoErr != nil {
		log.Fatalf("could not add todo: %v", addTodoErr)
	}
	log.Printf("New todo added: %v", addTodoResponse.GetAdded())

	editTodoResponse, editTodoErr := c.EditTodo(ctx, &pb.EditTodoParams{
		Todo: &pb.Todo{Id: lastTodo.Id, Name: fmt.Sprintf("Todo %v edited", lastTodo.Id), Complete: false},
	})

	if editTodoErr != nil {
		log.Fatalf("could not add todo: %v", editTodoErr)
	}
	log.Printf("Todo edited: %v", editTodoResponse.Edited)

	deleteTodoResponse, deleteTodoErr := c.DeleteTodo(ctx, &pb.DeleteTodoParams{
		TodoId: lastTodo.Id,
	})

	if deleteTodoErr != nil {
		log.Fatalf("could not add todo: %v", deleteTodoErr)
	}

	log.Printf("Todo deleted: %v", deleteTodoResponse.Deleted)
}
