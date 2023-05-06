package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"
	"github.com/c-m3-codin/url_shortner/utils"
	"github.com/gin-gonic/gin"
)

// Handler function for creating a shortened URL.
func CreateShortLink(c *gin.Context) {
	// Declare a variable to hold the incoming request body.
	var shortLink models.ShortLink

	// Bind the request body to the shortLink variable.
	err := c.BindJSON(&shortLink)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check that the OriginalURL field is not empty.
	if shortLink.OriginalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing OriginalURL field"})
		// c.AbortWithError(http.StatusBadRequest, errors.New("missing OriginalURL field"))
		return
	}

	// Count the number of active short URLs in the database.
	var count int64
	err = services.DB.Model(&models.ShortLink{}).Where("expires_at >= ?", time.Now()).Count(&count).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Check if the number of active short URLs is less than 20000.
	if count >= 20000 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Too many active URLs"})
		// c.AbortWithError(http.StatusTooManyRequests, errors.New("maximum number of active short URLs reached"))
		return
	}

	// Generate a random shortened URL.
	shortenedUrl := utils.GenerateShortenedURL()
	fmt.Printf("\n\n" + "" + c.Request.Host + "/" + shortenedUrl + "\n\n")
	shortenedUrl = c.Request.Host + "/" + shortenedUrl
	// Set the expiry time to 24 hours from now.
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert the new shortened URL into the database.
	err = services.DB.Create(&models.ShortLink{
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
	c.JSON(http.StatusOK, gin.H{
		"original_url":  shortLink.OriginalURL,
		"shortened_url": shortenedUrl,
		"expires_at":    expiresAt,
	})
	// c.JSON(http.StatusOK, shortLink)
}
