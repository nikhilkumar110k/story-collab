package utils

import (
	"main/events"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.GET("/GetAuthors", events.GetAuthors)
	authenticated.POST("/createauthor", events.CreateAuthors)
	authenticated.POST("/deleteauthor", events.DeleteAuthors)
	authenticated.POST("/createstories", events.CreateStory)
}
