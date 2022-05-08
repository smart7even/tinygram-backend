package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/repository"
	"github.com/smart7even/golang-do/internal/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	address := os.Getenv("ADRESS")

	db, err := sql.Open("mysql", dbConnectionString)

	todoRepo := repository.NewMySQLTodoRepo(db)
	todoService := service.NewTodoService(todoRepo)

	if err != nil {
		fmt.Printf("Can't prepare driver to connect to db: %v", err)
		return
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Can't connect to db: %v", err)
		return
	}

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Handling /tasks request, method %v\n", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "GET" {
			todos, err := todoService.ReadAll()

			if err != nil {
				response := fmt.Sprintf("Can't read todos: %v", err)
				fmt.Println(response)
				fmt.Fprintln(w, "Can't read todos")
				return
			}

			todosJson, err := json.Marshal(todos)

			if err != nil {
				response := fmt.Sprintf("Error while encoding todos to JSON: %v", err)
				fmt.Println(response)
				fmt.Fprintln(w, "Error while encoding todos to JSON")
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(todosJson))
		} else if r.Method == "POST" {
			if b, err := io.ReadAll(r.Body); err == nil {
				var todo domain.Todo
				json.Unmarshal(b, &todo)
				err := todoService.Create(todo)

				if err != nil {
					response := fmt.Sprintf("Error while creating todos: %v", err)
					fmt.Println(response)
					fmt.Fprintln(w, "Error while creating todos")
					return
				}

				w.WriteHeader(201)
				w.Header().Set("Content-Type", "plain/text")
				fmt.Fprint(w, "Task added")
			} else {
				fmt.Printf("Error while parsing request body: %v", err)
			}
		} else if r.Method == "OPTIONS" {
			w.WriteHeader(200)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		todoId, err := strconv.Atoi(strings.Split(r.RequestURI, "/")[2])

		if err != nil {
			fmt.Printf("Error while parsing request URL: %v", err)
			return
		}

		fmt.Printf("Handling /tasks/{%v} request, method %v\n", todoId, r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "PUT" {
			if b, err := io.ReadAll(r.Body); err == nil {
				var todo domain.Todo
				todo.Id = int64(todoId)
				json.Unmarshal(b, &todo)

				err := todoService.Update(todo)

				if err != nil {
					if errors.Is(err, service.TodoDoesNotExist{}) {
						response := fmt.Sprintf("Todo with id %v doesn't exist", err.(service.TodoDoesNotExist).TodoId)
						fmt.Println(response)
						fmt.Fprintln(w, response)
					}
					return
				}

				fmt.Fprintf(w, "Task with id %v edited", todoId)
			}
		} else if r.Method == "DELETE" {
			err := todoService.Delete(int64(todoId))

			if err != nil {
				fmt.Fprintf(w, "There is no task with id %v", todoId)
				return
			}

			fmt.Fprint(w, "Task deleted")
		}
	})

	log.Fatal(http.ListenAndServe(address, nil))
}
