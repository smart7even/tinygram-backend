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
	"strconv"
	"strings"

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

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		todoId, err := strconv.Atoi(strings.Split(r.RequestURI, "/")[2])

		if err != nil {
			fmt.Printf("Error while parsing request URL: %v", err)
			return
		}

		if r.Method == "PUT" {
			if b, err := io.ReadAll(r.Body); err == nil {
				var todo Todo
				todo.Id = int64(todoId)
				json.Unmarshal(b, &todo)

				res, err := db.Exec("UPDATE todos SET name = $1, complete = $2 WHERE id = $3", todo.Name, todo.Complete, todo.Id)

				if err != nil {
					fmt.Printf("Error while editing todo: %v", err)
					return
				}

				rowsAffected, err := res.RowsAffected()

				if err != nil {
					fmt.Printf("Error while getting affected rows: %v", err)
					return
				}

				if rowsAffected == 1 {
					fmt.Fprint(w, "Task edited")
				} else {
					fmt.Fprintf(w, "There is no task with id %v", todoId)
				}
			}
		} else if r.Method == "DELETE" {
			res, err := db.Exec("DELETE FROM todos WHERE id = $1", todoId)

			if err != nil {
				fmt.Printf("Error while deleting todo: %v", err)
				return
			}

			rowsAffected, err := res.RowsAffected()

			if err != nil {
				fmt.Printf("Error while getting affected rows: %v", err)
				return
			}

			if rowsAffected == 1 {
				fmt.Fprint(w, "Task deleted")
			} else {
				fmt.Fprintf(w, "There is no task with id %v", todoId)
			}
		}

	})

	log.Fatal(http.ListenAndServe(address, nil))
}
