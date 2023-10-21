package handlers

import (
	"fmt"
	"net/http"

	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"
	"github.com/gin-gonic/gin"
)

type ShortLinks []models.ShortLink

func GetUrls(c *gin.Context) {

	userId := c.GetUint("userID")

	var shortLinks ShortLinks

	services.DB.Where("user_id = ?", userId).Find(&shortLinks)

	shortLink_response := shortLinks.generateResponsense()

	fmt.Println(shortLink_response)

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, shortLink_response)
	// db.Where("name <> ?", "jinzhu").Find(&users)

}

func (shortLinks ShortLinks) generateResponsense() (shortLink_response []models.ShortLink_response) {

	// []models.ShortLink_response
	for _, link := range shortLinks {
		shortLink_response = append(shortLink_response, models.ShortLink_response{
			OriginalURL:  link.OriginalURL,
			ExpiresAt:    link.ExpiresAt,
			ShortenedURL: link.ShortenedURL,
		})
	}

	return

}
