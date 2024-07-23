package handlers

import (
	"fmt"
	"net/http"

	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"
	"github.com/gin-gonic/gin"
)

func GetHits(c *gin.Context) {

	userID := c.GetUint("userID") // Replace with the actual user's ID.

	shortUrl_received := c.Param("shortenedUrl")
	fmt.Println("Hit all hits with userID: ", userID, " and shortURL: ", shortUrl_received)

	if shortUrl_received == "" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, "Erron: Request with no shortURL ")
		return
	}

	var shorturl models.ShortLink
	services.DB.Where("shortened_url = ?", shortUrl_received).Find(&shorturl)

	// check if authorised to check the given url
	if shorturl.UserID != userID {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusUnauthorized, "Erron: Not Authorised to access ")
		return

	}
	type RedirectRequests_list []models.RedirectRequests
	var allHits RedirectRequests_list
	services.DB.Where("shortened_url = ?", shortUrl_received).Find(&allHits)
	fmt.Println(allHits)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, allHits)
	return

}
