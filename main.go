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

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

	db, err := sql.Open("mysql", dbConnectionString)

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
				db.QueryRow("INSERT INTO todos(name, complete) VALUES (?, ?)", todo.Name, todo.Complete)

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
				var todo Todo
				todo.Id = int64(todoId)
				json.Unmarshal(b, &todo)

				res, err := db.Exec("UPDATE todos SET name = ?, complete = ? WHERE id = ?", todo.Name, todo.Complete, todo.Id)

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
			res, err := db.Exec("DELETE FROM todos WHERE id = ?", todoId)

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
