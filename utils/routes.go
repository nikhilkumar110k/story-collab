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
	authenticated.GET("/GetAuthors", events.GetAuthors)
	authenticated.POST("/createauthor", events.CreateAuthors)
	authenticated.POST("/deleteauthor", events.DeleteAuthors)
	authenticated.POST("/createstories", events.CreateStory)
	server.POST("/signup", Signup)
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

	passValid := Checkpass(u.Password, user.Password)
	if !passValid {
		return errors.New("incorrect password")
	}

	u.ID = user.ID
	return nil
}

func Signup(ctx *gin.Context) {
	var user User

	// ✅ Parse request body
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "error": err.Error()})
		return
	}

	// ✅ Hash the password before saving it
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error hashing password", "error": err.Error()})
		return
	}

	// ✅ Get database connection
	queries, exists := ctx.MustGet("queries").(*db.Queries)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	// ✅ Insert user into the database
	insertedUser, err := queries.CreateAuthor(context.Background(), db.CreateAuthorParams{
		Name: user.Name,
		Bio:  hashedPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error saving user", "error": err.Error()})
		return
	}

	// ✅ Respond with success
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "signed up successfully",
		"user_id": insertedUser.ID,
	})
}

func Login(ctx *gin.Context) {
	var user User

	// ✅ Decode JSON request body
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	// ✅ Validate user
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
