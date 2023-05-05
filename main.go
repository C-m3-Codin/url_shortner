package main

import (
	"net/http"
	"url_shortner/services"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		services.LoadEnv()     //loading env
		services.NewDatabase() //new database connection
		context.JSON(http.StatusOK, gin.H{"data": "Hello World !"})
	})
	router.Run(":8000")
}
