package http_handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smart7even/golang-do/internal/domain"
)

func (h *Handler) makeDeviceRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		devices, err := h.Services.Device.ReadAll(user.Id)

		if err != nil {
			c.String(400, "Can't read devices")
			return
		}

		c.JSON(200, devices)
	})

	r.GET("/:id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(400, "Invalid id")
			return
		}

		device, err := h.Services.Device.Read(id, user.Id)

		if err != nil {
			c.String(400, "Can't read device")
			return
		}

		c.JSON(200, device)
	})

	r.POST("/", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		var device domain.Device

		if err := c.ShouldBindJSON(&device); err != nil {
			c.String(400, "Invalid request")
			return
		}

		device.UserId = user.Id

		err := h.Services.Device.Create(&device)

		if err != nil {
			c.String(400, "Can't create device")
			return
		}

		c.JSON(200, device)
	})

	r.PUT("/:id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		_, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(400, "Invalid id")
			return
		}

		var device domain.Device

		if err := c.ShouldBindJSON(&device); err != nil {
			c.String(400, "Invalid request")
			return
		}

		device.UserId = user.Id

		err = h.Services.Device.Update(&device)

		if err != nil {
			c.String(400, "Can't update device")
			return
		}

		c.String(200, "Device updated")
	})

	r.DELETE("/:id", func(c *gin.Context) {
		user := c.MustGet("user").(*domain.User)

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.String(400, "Invalid id")
			return
		}

		err = h.Services.Device.Delete(id, user.Id)

		if err != nil {
			c.String(400, "Can't delete device")
			return
		}

		c.String(200, "Device deleted")
	})
}
