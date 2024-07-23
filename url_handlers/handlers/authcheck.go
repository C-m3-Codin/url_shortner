package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAuth(c *gin.Context) {
	username := c.GetString("username")
	response := "pong " + username
	c.JSON(http.StatusOK, gin.H{"message": response})
}
