package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"url_shortner/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Declare a global variable to hold the database connection.
var db *gorm.DB

// Define a struct to represent a shortened URL.
type ShortLink struct {
	ID           int       `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortenedURL string    `json:"shortened_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func main() {
	// Declare a variable to hold the error returned from NewDatabase().
	var err error

	// Create a database connection and store it in the global variable db.
	db, err = services.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Create the necessary tables in the database.
	err = db.AutoMigrate(&ShortLink{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a goroutine to periodically delete expired links from the database.
	go func() {
		for {
			fmt.Print("Checking for expiry")
			err := db.Where("expires_at < ?", time.Now()).Delete(&ShortLink{}).Error
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Minute)
		}
	}()

	// Create a new Gin router.
	r := gin.Default()

	// Define a route for creating shortened URLs.
	r.POST("/create", createShortLink)

	// Define a route for redirecting to a shortened URL.
	r.GET("/:shortenedUrl", redirectShortLink)

	// Start the server on port 8000.
	err = r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

// Handler function for creating a shortened URL.
func createShortLink(c *gin.Context) {
	// Declare a variable to hold the incoming request body.
	var shortLink ShortLink

	// Bind the request body to the shortLink variable.
	err := c.BindJSON(&shortLink)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check that the OriginalURL field is not empty.
	if shortLink.OriginalURL == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("missing OriginalURL field"))
		return
	}

	// Count the number of active short URLs in the database.
	var count int64
	err = db.Model(&ShortLink{}).Where("expires_at >= ?", time.Now()).Count(&count).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check if the number of active short URLs is less than 20000.
	if count >= 20000 {
		c.AbortWithError(http.StatusTooManyRequests, errors.New("maximum number of active short URLs reached"))
		return
	}

	// Generate a random shortened URL.
	shortenedUrl := generateShortenedURL()

	// Set the expiry time to 24 hours from now.
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert the new shortened URL into the database.
	err = db.Create(&ShortLink{
		OriginalURL:  shortLink.OriginalURL,
		ShortenedURL: shortenedUrl,
		ExpiresAt:    expiresAt,
	}).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Update the shortLink variable with the shortened URL and expiry time.
	shortLink.ShortenedURL = shortenedUrl
	shortLink.ExpiresAt = expiresAt

	// Return the new shortened URL as a JSON response.
	c.JSON(http.StatusOK, shortLink)
}

func redirectShortLink(c *gin.Context) {
	shortenedUrl := c.Param("shortenedUrl")

	// Query the database for the original URL associated with the shortened URL
	var shortLink ShortLink
	err := db.Where("shortened_url = ?", shortenedUrl).First(&shortLink).Error
	if err != nil {
		// Return a 404 error if the shortened URL is not found in the database
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if time.Now().After(shortLink.ExpiresAt) {
		// Return a 404 error if the shortened URL has expired
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("link has expired"))
		return
	}

	// Redirect the user to the original URL
	fmt.Printf("Redirecting to ", shortLink.OriginalURL)
	c.Redirect(http.StatusPermanentRedirect, shortLink.OriginalURL)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortenedURL() string {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 8)
	for i := range b {
		// Pick a random character from the character set and add it to the byte slice
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
