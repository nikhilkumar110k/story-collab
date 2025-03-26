package middlewares

import (
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")
	if token == "" {
		context.AbortWithStatus(http.StatusNonAuthoritativeInfo)
		return
	}
	err, userid := utils.VerifyToken(token)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	context.Set("userid", userid)
	context.Next()

}
