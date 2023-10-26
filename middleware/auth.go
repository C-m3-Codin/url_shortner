package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"
	"github.com/c-m3-codin/url_shortner/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// middleware function to check authentication calls next function if valid auth else aborts
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")

		userID := ctx.DefaultQuery("user_id", "")

		fmt.Println("user is ", userID)

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

		if claims.IsAdmin {
			fmt.Println("Admin Useer Logged In   ")

			if userID != "" {

				fmt.Println("Admin Useer Logged In as userID : ", userID)

				num, err := strconv.Atoi(userID)
				if err != nil {
					fmt.Println("Conversion error:", err)
					ctx.Header("Content-Type", "application/json")
					ctx.JSON(http.StatusNonAuthoritativeInfo, "Erron: User Id incorrect ")
					ctx.Abort()
					return

				} else {
					userId_uint := uint(num)
					fmt.Printf("Converted value as uint: %d\n", userId_uint)
					services.DB.Where("ID=?", userId_uint)

					var user models.User
					services.DB.Where("ID=?", userId_uint).First(&user)

					if user.Username == "" {

						fmt.Println("No user found with the user id mentioned : ", userId_uint)
						fmt.Println("Conversion error:", err)
						ctx.Header("Content-Type", "application/json")
						ctx.JSON(http.StatusNonAuthoritativeInfo, "Erron: User Id incorrect ")
						ctx.Abort()
						return
					}

					fmt.Println("Found user deets : ", user)

					ctx.Set("email", user.Email)
					ctx.Set("username", user.Username)
					ctx.Set("userID", user.ID)
					ctx.Set("isAdmin", user.IsAdmin)
					ctx.Next()

				}

			}

		}

		ctx.Set("email", claims.Email)
		ctx.Set("username", claims.Username)
		ctx.Set("userID", claims.UserId)
		ctx.Set("isAdmin", claims.IsAdmin)

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
