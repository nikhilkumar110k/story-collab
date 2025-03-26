package utils

import (
	"context"
	"database/sql"
	"errors"
	db "main/db/sqlc"
	"main/events"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RegisterRoutes(server *gin.Engine) {
	authenticated := server.Group("/")
	authenticated.Use(Authenticate)
	server.GET("/GetAuthors", events.GetAuthors)
	server.POST("/createauthor", events.CreateAuthors)
	server.POST("/deleteauthor", events.DeleteAuthors)
	server.POST("/createstories", events.CreateStory)
	server.POST("/signup", SignUp)
	server.POST("/login", Login)
}

const Secretkey = "totallsecretkeylol"

type User struct {
	ID       int64  `json:"id" binding:"-"`
	Name     string `json:"name" binding:"-"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Validate(ctx *gin.Context) error {
	queries, exists := ctx.MustGet("queries").(*db.Queries)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return errors.New("database connection not found")
	}

	user, err := queries.GetAuthorsByEmail(context.Background(), u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("email not found")
		}
		return err
	}

	u.ID = user.ID
	return nil
}

func SignUp(ctx *gin.Context) {
	queries := ctx.MustGet("queries").(*db.Queries)

	var req struct {
		Name     string `json:"name"`
		Bio      string `json:"bio"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	author, err := queries.CreateAuthor(ctx, db.CreateAuthorParams{
		Name:     req.Name,
		Bio:      req.Bio,
		Email:    req.Email,
		Password: req.Password, // Store hashed password, not plaintext
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}

	ctx.JSON(http.StatusOK, author)
}

func Login(ctx *gin.Context) {
	var user User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	if err := user.Validate(ctx); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Generatetoken(email string, userid int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email":  email,
		"userid": userid,
		"exp":    time.Now().Add(time.Hour * 20).Unix(),
	})
	return token.SignedString([]byte(Secretkey))
}

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatus(http.StatusNonAuthoritativeInfo)
		return
	}
	err, userid := VerifyToken(token)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	context.Set("userid", userid)
	context.Next()

}

func VerifyToken(token string) (error, int64) {
	tokenparsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Wrong token")
		}
		return []byte(Secretkey), nil
	})
	if err != nil {
		return errors.New("couldn't parse"), 0
	}
	validornot := tokenparsed.Valid
	if !validornot {
		return errors.New("the token is not valid or just expired"), 0
	}
	claims, ok := tokenparsed.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Invalid Claims"), 0
	}
	useridRaw, exists := claims["userid"]
	if !exists {
		return errors.New("userid not found in claims"), 0
	}
	var Userid int64
	switch v := useridRaw.(type) {
	case float64:
		Userid = int64(v)
	case int64:
		Userid = v
	default:
		return errors.New("unexpected type for userid in claims"), 0
	}

	return nil, Userid
}
