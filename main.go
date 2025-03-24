package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	db "main/db/sqlc"
	storiesdb "main/storydb/sqlcstory"
	"main/utils"
)

func main() {
	mainDBURL := "postgresql://root:Nikhil@123k@localhost:5432/project3?sslmode=disable"
	mainConn, err := pgxpool.New(context.Background(), mainDBURL)
	if err != nil {
		log.Fatal("Cannot connect to main database:", err)
	}
	defer mainConn.Close()
	mainQueries := db.New(mainConn)

	storiesDBURL := "postgresql://root:Nikhil@123k@localhost:5432/storiesdb?sslmode=disable"
	storiesConn, err := pgxpool.New(context.Background(), storiesDBURL)
	if err != nil {
		log.Fatal("Cannot connect to stories database:", err)
	}
	defer storiesConn.Close()
	storiesQueries := storiesdb.New(storiesConn)

	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("queries", mainQueries)
		ctx.Set("queries1", storiesQueries)
		ctx.Next()
	})

	utils.RegisterRoutes(router)

	log.Fatal(router.Run(":8080"))
}
