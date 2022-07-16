package http_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smart7even/golang-do/internal/domain"
	"github.com/smart7even/golang-do/internal/service"
)

type Handler struct {
	Services service.Services
}

func (h *Handler) InitAPI() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware)

	r.GET("/tasks", func(c *gin.Context) {
		todos, err := h.Services.Todo.ReadAll()

		if err != nil {
			response := fmt.Sprintf("Can't read todos: %v", err)
			fmt.Println(response)
			c.String(400, "Can't read todos")
			return
		}

		c.JSON(200, todos)
	})

	r.POST("/tasks", func(c *gin.Context) {
		if b, err := io.ReadAll(c.Request.Body); err == nil {
			var todo domain.Todo
			json.Unmarshal(b, &todo)
			err := h.Services.Todo.Create(todo)

			if err != nil {
				fmt.Printf("Error while creating todos: %v", err)
				c.String(500, "Error while creating todos")
				return
			}

			c.String(201, "Task added")
		} else {
			fmt.Printf("Error while parsing request body: %v", err)
		}
	})

	r.PUT("/tasks/:id", func(c *gin.Context) {
		todoId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			fmt.Printf("Error while parsing todoId from path: %v", err)
			c.String(400, "Error while parsing todoId from path")
		}

		if b, err := io.ReadAll(c.Request.Body); err == nil {
			var todo domain.Todo
			todo.Id = int64(todoId)
			json.Unmarshal(b, &todo)

			err := h.Services.Todo.Update(todo)

			if err != nil {
				if errors.Is(err, service.TodoDoesNotExist{}) {
					response := fmt.Sprintf("Todo with id %v doesn't exist", err.(service.TodoDoesNotExist).TodoId)
					fmt.Println(response)
					c.String(400, response)
				}
				return
			}

			c.String(200, "Task with id %v edited", todoId)
		}
	})

	r.DELETE("/tasks/:id", func(c *gin.Context) {
		todoId, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			fmt.Printf("Error while parsing todoId from path: %v", err)
			c.String(400, "Error while parsing todoId from path")
		}

		err = h.Services.Todo.Delete(int64(todoId))

		if err != nil {
			c.String(400, "There is no task with id %v", todoId)
			return
		}

		c.String(200, "Task deleted")
	})

	r.POST("/users", func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.String(400, "Token is required")
			return
		}

		err := h.Services.User.Create(token)

		if err != nil {
			fmt.Printf("Error while creating user: %v", err)
			c.String(400, "Can't create user")
			return
		}

		c.String(200, "User created")
	})

	r.GET("/users", func(c *gin.Context) {
		users, err := h.Services.User.ReadAll()

		if err != nil {
			response := fmt.Sprintf("Can't read users: %v", err)
			fmt.Println(response)
			c.String(400, "Can't read users")
			return
		}

		c.JSON(200, users)
	})

	r.DELETE("/users", func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.String(400, "Token is required")
			return
		}

		err := h.Services.User.Delete(token)

		if err != nil {
			fmt.Printf("Error while deleting user: %v", err)
			c.String(400, "Can't delete user")
			return
		}

		c.String(200, "User deleted")
	})

	return r
}
