package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	db "main/db/sqlc"
	"main/utils"
)

func main() {
	dbURL := "postgresql://root:Nikhil@123k@localhost:5432/project3?sslmode=disable"
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer conn.Close()
	queries := db.New(conn)

	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("queries", queries)
		ctx.Next()
	})
	utils.RegisterRoutes(router)

	log.Fatal(router.Run(":8080"))
}
