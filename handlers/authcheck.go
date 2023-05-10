package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Pong"})
}
