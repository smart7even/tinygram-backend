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
		fmt.Printf("Token is required\n")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userId, err := h.Services.Auth.Verify(token)

	if err != nil {
		fmt.Printf("Error while verifying token: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := h.Services.User.Read(userId)

	if err != nil {
		fmt.Printf("Error while reading user: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user", &user)
	c.Next()
}
