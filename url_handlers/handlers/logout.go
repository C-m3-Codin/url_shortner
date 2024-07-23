package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/c-m3-codin/url_shortner/services"
	"github.com/gin-gonic/gin"
)

func Logout(ctx *gin.Context) {
	tokenStr := ctx.GetHeader("Authorization")
	err := BlacklistToken(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": "User Logged out",
	})
	// return nil

}

func BlacklistToken(token string) error {
	// Set the token as a key in Redis with an expiration time (e.g., token's expiration time)
	// This ensures the token is automatically removed from the blacklist when it expires.
	expirationTime := time.Hour
	err := services.RedisClient.Set(context.Background(), "blacklist:"+token, "1", expirationTime).Err()
	if err != nil {
		return err
	}
	return nil
}
