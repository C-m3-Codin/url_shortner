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

// handler to Redirect to the original link from the shortened link
func RedirectShortLink(c *gin.Context) {

	shortenedUrl := c.Param("shortenedUrl")
	// shortenedUrl = c.Request.Host + "/" + shortenedUrl

	// Query the database for the original URL associated with the shortened URL
	var shortLink models.ShortLink
	err := services.DB.Where("shortened_url = ?", shortenedUrl).FirstOrCreate(&shortLink).Error
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
	go utils.LogHit(&shortLink, c.Request.RemoteAddr)
	c.Redirect(http.StatusPermanentRedirect, shortLink.OriginalURL)
}
