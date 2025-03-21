package utils

import (
	"main/events"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.GET("/events", events.GetAuthors)
}
