package events

import (
	"context"
	db "main/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetAuthors(ctx *gin.Context) {

	queries, exists := ctx.MustGet("queries").(*db.Queries)
	if !exists {
		ctx.JSON(500, gin.H{"error": "Database connection not found"})
		return
	}
	defer queries.Close()

	authors, err := queries.ListAuthors(context.Background())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, authors)
}

func CreateAuthors(ctx *gin.Context) {
	queriesInterface, exists := ctx.Get("queries")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error retrieving queries"})
		return
	}

	queries, ok := queriesInterface.(*db.Queries)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "invalid queries type"})
		return
	}

	var req struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	authorParams := db.CreateAuthorParams{
		Name: req.Name,
		Bio:  pgtype.Text{String: req.Bio, Valid: true}.String,
	}

	createdAuthor, err := queries.CreateAuthor(ctx, authorParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, createdAuthor)
}

func DeleteAuthors(ctx *gin.Context) {
	queriesInterface, exists := ctx.Get("queries")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error retrieving queries"})
		return
	}
	queries, ok := queriesInterface.(*db.Queries)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "invalid queries type"})
		return
	}
	var authortodelete struct {
		ID int `json:"id"`
	}
	if err := ctx.ShouldBindJSON(&authortodelete); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	err := queries.DeleteAuthor(ctx, int64(authortodelete.ID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "unable to delete the author"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "author successfully deleted"})
}
