package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(server *gin.Engine, adapter Adapter) {
	server.GET("/healthcheck", healthcheck(adapter))
	server.GET("/user-event", getUserEvent(adapter))
}

func healthcheck(adapter Adapter) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"healthcheck": "pass"})
	}
}

func getUserEvent(adapter Adapter) func(*gin.Context) {
	return func(c *gin.Context) {
		userEvent, err := adapter.api.GetUserEvent()
		if err != nil {
			message := fmt.Sprintf("failed to generate a user event: %s", err)

			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": message},
			)
			return
		}

		if err != nil {
			message := fmt.Sprintf("failed to generate a user event: %s", err)

			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": message},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			userEvent,
		)
	}
}
