package middleware

import (
	"net/http"

	"github.com/c-m3-codin/url_shortner/utils"
	"github.com/gin-gonic/gin"
)

// middleware function to check authentication calls next function if valid auth else aborts
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")

		if tokenStr == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No Auth token in header"})
			ctx.Abort()
			return
		}
		err, claims := utils.ValidateToken(tokenStr)
		// fmt.Println(err.Error())
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
