package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Todo struct {
	Id       int64
	Name     string
	Complete bool
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	address := os.Getenv("ADRESS")

	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		fmt.Printf("Can't connect to db: %v", err)
		return
	}

	db.Ping()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			rows, err := db.Query("SELECT id, name, complete FROM todos")

			if err != nil {
				fmt.Printf("Error while requesting todos: %v", err)
				fmt.Fprint(w, "Error", html.EscapeString(r.URL.Path))
				return
			}

			defer rows.Close()

			var todos []Todo

			for rows.Next() {
				var todo Todo
				rows.Scan(&todo.Id, &todo.Name, &todo.Complete)
				todos = append(todos, todo)
			}

			todosJson, err := json.Marshal(todos)

			if err != nil {
				fmt.Printf("Error while encoding todos to JSON: %v", err)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(todosJson))
		} else if r.Method == "POST" {
			if b, err := io.ReadAll(r.Body); err == nil {
				var todo Todo
				json.Unmarshal(b, &todo)
				db.QueryRow("INSERT INTO todos(name, complete) VALUES ($1, $2)", todo.Name, todo.Complete)
				fmt.Fprint(w, "Task added")
			}
		}
	})

	log.Fatal(http.ListenAndServe(address, nil))
}
