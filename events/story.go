package events

import (
	storydb "main/storydb/sqlcstory"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateStory(ctx *gin.Context) {
	queryinterface, exists := ctx.Get("queries1")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error: couldn't fetch query interface"})
		return
	}

	query, ok := queryinterface.(*storydb.Queries)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "query not found"})
		return
	}

	var story struct {
		Originalstory  string `json:"originalstory"`
		Pulledrequests int    `json:"pulledrequests"`
		Updatedstory   string `json:"updatedstory"`
		AuthorID       int    `json:"authorid"`
	}

	if err := ctx.ShouldBindJSON(&story); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	storyvalues := storydb.CreateStoryParams{
		Originalstory:  story.Originalstory,
		Pulledrequests: story.Pulledrequests,
		Updatedstory:   story.Updatedstory,
		AuthorID:       story.AuthorID,
	}

	createdstory, err := query.CreateStory(ctx, storyvalues)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error creating the story"})
		return
	}

	ctx.JSON(http.StatusAccepted, createdstory)
}
