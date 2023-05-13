package main

import (
	"fmt"
	"log"
	"time"

	"github.com/c-m3-codin/url_shortner/handlers"
	"github.com/c-m3-codin/url_shortner/middleware"
	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"

	"github.com/gin-gonic/gin"
)

// Declare a global variable to hold the database connection.
// var db *gorm.DB

func main() {
	// Declare a variable to hold the error returned from NewDatabase().
	var err error

	// Create a database connection and store it in the global variable db.
	_, err = services.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	// Create a goroutine to periodically delete expired links from the database.
	go func() {
		for {
			fmt.Print("Checking for expiry")
			err := services.DB.Where("expires_at < ?", time.Now()).Delete(&models.ShortLink{}).Error
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Minute)
		}
	}()

	// Create a new Gin router.
	r := gin.Default()

	// Define a route for redirecting to a shortened URL.
	r.GET("/:shortenedUrl", handlers.RedirectShortLink)

	// register a user with creds
	r.POST("/user/register", handlers.RegisterUser)

	// get token by logging in
	r.POST("/token", handlers.GenerateToken)

	secured := r.Group("/sec").Use(middleware.Auth())
	{
		secured.GET("/ping", handlers.CheckAuth)

		// Define a route for creating shortened URLs.
		r.POST("/create", handlers.CreateShortLink)

	}

	// Start the server on port 8000.
	err = r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
