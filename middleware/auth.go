package middleware

import (
	"net/http"

	"github.com/c-m3-codin/url_shortner/utils"
	"github.com/gin-gonic/gin"
)

// middleware function to check authentication calls next function if valid auth else aborts
func Auth(context *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := context.GetHeader("Authorization")
		if tokenStr == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "No Auth token in header"})
			context.Abort()
			return
		}
		err := utils.ValidateToken(tokenStr)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			context.Abort()
			return
		}
		context.Next()
	}
}
