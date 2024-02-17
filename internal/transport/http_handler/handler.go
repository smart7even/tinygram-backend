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

	h.makeTodosRoutes(r)

	r.POST("/auth/token", func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.String(400, "Token is required")
			return
		}

		user, err := h.Services.User.ReadByToken(token)

		if err != nil {
			fmt.Printf("Error while reading user: %v", err)
			c.String(400, "Can't read user")
			return
		}

		appUser, err := h.Services.User.Read(user.Id)

		// create user if not exists
		if errors.Is(err, service.UserDoesNotExist{UserId: user.Id}) {
			err := h.Services.User.Create(token)

			if err != nil {
				fmt.Printf("Error while creating user: %v", err)
				c.String(400, "Can't create user")
				return
			}

			appUser, err = h.Services.User.Read(user.Id)

			if err != nil {
				fmt.Printf("Error while reading user: %v", err)
				c.String(400, "Can't read user")
				return
			}
		}

		appToken, err := h.Services.Auth.Sign(user.Id)

		if err != nil {
			fmt.Printf("Error while signing token: %v", err)
			c.String(400, "Can't sign token")
			return
		}

		c.JSON(200, gin.H{
			"token": appToken,
			"user":  appUser,
		})
	})

	usersGroup := r.Group("/users")
	{
		h.makeUsersRoutes(usersGroup)
	}

	chatGroup := r.Group("/chat")
	chatGroup.Use(h.authMiddleware)
	{
		h.makeChatRoutes(chatGroup)
		h.makeMessageRoutes(chatGroup)
	}

	reminderGroup := r.Group("/reminder")
	reminderGroup.Use(h.authMiddleware)
	{
		h.makeReminderRoutes(reminderGroup)
	}

	deviceGroup := r.Group("/device")
	deviceGroup.Use(h.authMiddleware)
	{
		h.makeDeviceRoutes(deviceGroup)
	}

	return r
}

func (h *Handler) makeTodosRoutes(r *gin.Engine) {
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

}

func (h *Handler) makeUsersRoutes(r *gin.RouterGroup) {
	r.POST("/", func(c *gin.Context) {
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

	r.GET("/", func(c *gin.Context) {
		users, err := h.Services.User.ReadAll()

		if err != nil {
			response := fmt.Sprintf("Can't read users: %v", err)
			fmt.Println(response)
			c.String(400, "Can't read users")
			return
		}

		c.JSON(200, users)
	})

	r.DELETE("/", func(c *gin.Context) {
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
}

func (h *Handler) makeChatRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		chat, err := h.Services.Chat.ReadAll()

		if err != nil {
			response := fmt.Sprintf("Can't read chats: %v", err)
			fmt.Println(response)
			c.String(400, "Can't read chats")
			return
		}

		c.JSON(200, chat)
	})

	r.POST("/", func(c *gin.Context) {
		if b, err := io.ReadAll(c.Request.Body); err == nil {
			var chat domain.Chat
			json.Unmarshal(b, &chat)
			err := h.Services.Chat.Create(chat)

			if err != nil {
				fmt.Printf("Error while creating chat: %v", err)
				c.String(500, "Error while creating chat")
				return
			}

			c.String(201, "Chat added")
		} else {
			fmt.Printf("Error while parsing request body: %v", err)
		}
	})

	r.PUT("/:id", func(c *gin.Context) {
		chatId := c.Param("id")

		if chatId == "" {
			fmt.Printf("Error while parsing chatId from path")
			c.String(400, "Error while parsing chatId from path")
		}

		if b, err := io.ReadAll(c.Request.Body); err == nil {
			var chat domain.Chat
			json.Unmarshal(b, &chat)
			chat.Id = chatId

			err := h.Services.Chat.Update(chat)

			if err != nil {
				response := fmt.Sprintf("Error while updating chat: %v", err)
				fmt.Println(response)
				c.String(400, response)
				return
			}

			c.String(200, "Chat with id %v edited", chatId)
		}
	})

	r.DELETE("/:id", func(c *gin.Context) {
		chatId := c.Param("id")

		if chatId == "" {
			fmt.Printf("Error while parsing chatId from path")
			c.String(400, "Error while parsing chatId from path")
		}

		err := h.Services.Chat.Delete(chatId)

		if err != nil {
			response := fmt.Sprintf("Error while deleting chat: %v", err)
			fmt.Println(response)
			c.String(400, response)
			return
		}

		c.String(200, "Chat deleted")
	})

	r.POST("/:id/user/:user_id", func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token == "" {
			c.String(400, "Token is required")
			return
		}

		user, err := h.Services.User.ReadByToken(token)

		if err != nil {
			fmt.Printf("Error while reading user: %v", err)
			c.String(400, "Can't read user")
			return
		}

		chatId := c.Param("id")
		userId := c.Param("user_id")

		if user.Id != userId {
			c.String(400, "You can't add user to chat")
			return
		}

		err = h.Services.Chat.Join(chatId, userId)

		if err != nil {
			fmt.Printf("Error while joining chat: %v", err)
			c.String(400, "Can't join chat")
			return
		}

		c.String(200, "User added to chat")
	})
}

func (h *Handler) makeMessageRoutes(r *gin.RouterGroup) {
	r.POST("/:id/message", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)
		chatId := c.Param("id")

		if b, err := io.ReadAll(c.Request.Body); err == nil {
			var message domain.Message
			json.Unmarshal(b, &message)
			message.ChatId = chatId
			message.UserId = user.Id
			err := h.Services.Message.Create(message, user.Id)

			if err != nil {
				fmt.Printf("Error while creating message: %v", err)
				c.String(500, "Error while creating message")
				return
			}

			c.String(201, "Message added")
		} else {
			fmt.Printf("Error while parsing request body: %v", err)
		}
	})

	r.GET("/:id/message", func(c *gin.Context) {
		chatId := c.Param("id")

		if chatId == "" {
			fmt.Printf("Error while parsing chatId from path")
			c.String(400, "Error while parsing chatId from path")
		}

		user := c.MustGet("user").(*domain.User)

		messages, err := h.Services.Message.ReadAll(chatId, user.Id)

		if err != nil {
			response := fmt.Sprintf("Can't read messages: %v", err)
			fmt.Println(response)
			c.String(400, "Can't read messages")
			return
		}

		c.JSON(200, messages)
	})

	r.PUT("/:id/message", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		messageId := c.Param("id")

		if messageId == "" {
			fmt.Printf("Error while parsing messageId from path")
			c.String(400, "Error while parsing messageId from path")
		}

		if b, err := io.ReadAll(c.Request.Body); err == nil {
			var message domain.Message
			message.Id = messageId
			json.Unmarshal(b, &message)

			err := h.Services.Message.Update(message, user.Id)

			if err != nil {
				response := fmt.Sprintf("Error while updating message: %v", err)
				fmt.Println(response)
				c.String(400, response)
				return
			}

			c.String(200, "Message with id %v edited", messageId)
		}
	})

	r.DELETE("/:id/message/:message_id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		messageId := c.Param("message_id")

		if messageId == "" {
			fmt.Printf("Error while parsing messageId from path")
			c.String(400, "Error while parsing messageId from path")
		}

		err := h.Services.Message.Delete(messageId, user.Id)

		if err != nil {
			response := fmt.Sprintf("Error while deleting message: %v", err)
			fmt.Println(response)
			c.String(400, response)
			return
		}

		c.String(200, "Message deleted")
	})
}
