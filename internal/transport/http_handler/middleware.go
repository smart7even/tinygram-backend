package http_handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func (h *Handler) authMiddleware(c *gin.Context) {
	token := c.GetHeader("token")

	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := h.Services.User.ReadByToken(token)

	if err != nil {
		fmt.Printf("Error while reading user: %v", err)
		c.String(400, "Can't read user")
		return
	}

	c.Set("user", user)
	c.Next()
}
