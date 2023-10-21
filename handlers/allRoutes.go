package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"
	"github.com/gin-gonic/gin"
)

func GetUrls(c *gin.Context) {

	userId := c.GetUint("userID")

	var shortLinks []models.ShortLink

	services.DB.Where("user_id = ?", userId).Find(&shortLinks)

	var shortLink_response []models.ShortLink_response
	for _, link := range shortLinks {
		shortLink_response = append(shortLink_response, models.ShortLink_response{
			OriginalURL:  link.OriginalURL,
			ShortenedURL: link.ShortenedURL,
			ExpiresAt:    link.ExpiresAt,
		})
	}

	fmt.Println(shortLink_response)

	jsonResponse, err := json.Marshal(shortLink_response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal response"})
		return
	}
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(jsonResponse))
	// db.Where("name <> ?", "jinzhu").Find(&users)

}
