package authentication

import (
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int64  `json:"id" binding:"-"`
	Name     string `json:"name" binding:"-"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u User) Validate() error {
	panic("unimplemented")
}
func Signup(context *gin.Context) {
	var user User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "error": err.Error()})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "error saving the credentials", "error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "signed up successfully", "user_id": user.ID})
}

func Login(context *gin.Context) {
	var user User
	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid input", "error": err.Error()})
		return
	}

	err = user.Validate()
	if err != nil {
		if err.Error() == "email not found" || err.Error() == "incorrect password" {
			context.JSON(http.StatusUnauthorized, gin.H{"message": "credentials didn't meet"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error", "error": err.Error()})
		}
		return
	}
	token, err := utils.Generatetoken(user.Email, user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"mesaage": "error generating token"})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"message": "succesfully got the token", "token": token})

	context.JSON(http.StatusAccepted, gin.H{"message": "login successful"})
}
