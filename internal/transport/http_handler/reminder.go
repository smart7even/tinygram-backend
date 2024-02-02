package http_handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smart7even/golang-do/internal/domain"
)

func (h *Handler) makeReminderRoutes(r *gin.RouterGroup) {
	// Handle GET requests

	r.GET("/", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		reminders, err := h.Services.Reminder.ReadAll(user.Id)

		if err != nil {
			c.String(400, "Can't read reminders")
			return
		}

		c.JSON(200, reminders)
	})

	r.GET("/:id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(400, "Invalid id")
			return
		}

		reminder, err := h.Services.Reminder.Read(id, user.Id)

		if err != nil {
			c.String(400, "Can't read reminder")
			return
		}

		c.JSON(200, reminder)
	})

	// Handle POST requests

	r.POST("/", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		var reminder domain.Reminder

		if err := c.ShouldBindJSON(&reminder); err != nil {
			c.String(400, "Invalid request")
			return
		}

		reminder.UserId = user.Id

		err := h.Services.Reminder.Create(reminder)

		if err != nil {
			c.String(400, "Can't create reminder")
			return
		}

		c.String(200, "Reminder created")
	})

	// Handle PUT requests

	r.PUT("/:id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(400, "Invalid id")
			return
		}

		var reminder domain.Reminder

		if err := c.ShouldBindJSON(&reminder); err != nil {
			c.String(400, "Invalid request")
			return
		}

		reminder.Id = id
		reminder.UserId = user.Id

		err = h.Services.Reminder.Update(reminder)

		if err != nil {
			c.String(400, "Can't update reminder")
			return
		}

		c.String(200, "Reminder updated")
	})

	// Handle DELETE requests

	r.DELETE("/:id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(400, "Invalid id")
			return
		}

		err = h.Services.Reminder.Delete(id, user.Id)

		if err != nil {
			c.String(400, "Can't delete reminder")
			return
		}

		c.String(200, "Reminder deleted")
	})
}
