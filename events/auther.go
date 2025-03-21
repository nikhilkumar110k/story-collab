package events

import (
	"context"
	db "main/db/sqlc"

	"github.com/gin-gonic/gin"
)

func getAuthors(ctx *gin.Context) {

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
