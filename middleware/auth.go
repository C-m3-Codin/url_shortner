package middleware

import (
	"net/http"

	"github.com/c-m3-codin/url_shortner/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
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

		isBlacklisted, err := IsTokenBlacklisted(tokenStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error "})
			ctx.Abort()
			return
		}

		if isBlacklisted {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
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
		ctx.Set("userID", claims.UserId)

		ctx.Next()
	}
}

func IsTokenBlacklisted(token string) (bool, error) {
	if token == "" {
		return false, nil // No token to check
	}

	// Check if the token exists in the Redis blacklist
	val, err := utils.GetRedisValue("blacklist:" + token)

	if val == "" || err == redis.Nil {
		return false, nil // Token is not blacklisted
	}

	return true, nil // Token is blacklisted
}
