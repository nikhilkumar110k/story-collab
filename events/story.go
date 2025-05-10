package events

import (
	"errors"
	"net/http"
	"time"

	db "main/db/sqlc"

	"github.com/gin-gonic/gin"
)

func GetQueries(c *gin.Context) (*db.Queries, error) {
	val, exists := c.Get("queries")
	if !exists {
		return nil, errors.New("queries not found in context")
	}

	queries, ok := val.(*db.Queries)
	if !ok {
		return nil, errors.New("invalid queries type in context")
	}

	return queries, nil
}

type StoryRequest struct {
	ID            string    `json:"id" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	Description   string    `json:"description"`
	CoverImage    string    `json:"cover_image"`
	AuthorID      string    `json:"author_id" binding:"required"`
	Likes         int64     `json:"likes"`
	Views         int64     `json:"views"`
	PublishedDate time.Time `json:"published_date"`
	LastEdited    time.Time `json:"last_edited"`
	StoryType     string    `json:"story_type"`
	Status        string    `json:"status" binding:"required"`
	Genres        []string  `json:"genres"`
}

func CreateStoryHandler(ctx *gin.Context) {
	var req StoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	queries, err := GetQueries(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	story, err := queries.CreateStory(ctx, db.CreateStoryParams{
		ID:            req.ID,
		Title:         req.Title,
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		AuthorID:      req.AuthorID,
		Likes:         req.Likes,
		Views:         req.Views,
		PublishedDate: req.PublishedDate,
		LastEdited:    req.LastEdited,
		StoryType:     req.StoryType,
		Status:        db.StoryStatus(req.Status),
		Genres:        req.Genres,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create story"})
		return
	}

	ctx.JSON(http.StatusOK, story)
}

func ListStoriesHandler(ctx *gin.Context) {
}

func GetStoryByIDHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	queries, err := GetQueries(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	story, err := queries.GetStoryByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Story not found"})
		return
	}

	ctx.JSON(http.StatusOK, story)
}

func UpdateStoryHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	var req StoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	queries, err := GetQueries(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	story, err := queries.UpdateStory(ctx, db.UpdateStoryParams{
		ID:            id,
		Title:         req.Title,
		Description:   req.Description,
		CoverImage:    req.CoverImage,
		AuthorID:      req.AuthorID,
		Likes:         req.Likes,
		Views:         req.Views,
		PublishedDate: req.PublishedDate,
		LastEdited:    req.LastEdited,
		StoryType:     req.StoryType,
		Status:        db.StoryStatus(req.Status),
		Genres:        req.Genres,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update story"})
		return
	}

	ctx.JSON(http.StatusOK, story)
}

func DeleteStoryHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	queries, err := GetQueries(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := queries.DeleteStory(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete story"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Story deleted"})
}
