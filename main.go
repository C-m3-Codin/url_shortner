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

var db *gorm.DB

type ShortLink struct {
	ID           int       `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortenedURL string    `json:"shortened_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func main() {

	var err error

	db, err = services.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&ShortLink{})
	if err != nil {
		log.Fatal(err)
	}

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
	r := gin.Default()
	r.POST("/create", createShortLink)
	r.GET("/:shortenedUrl", redirectShortLink)

	err = r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
func createShortLink(c *gin.Context) {
	var shortLink ShortLink
	err := c.BindJSON(&shortLink)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if shortLink.OriginalURL == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("missing OriginalURL field"))
		return
	}

	shortenedUrl := generateShortenedURL()
	expiresAt := time.Now().Add(24 * time.Hour)

	err = db.Create(&ShortLink{
		OriginalURL:  shortLink.OriginalURL,
		ShortenedURL: shortenedUrl,
		ExpiresAt:    expiresAt,
	}).Error
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	shortLink.ShortenedURL = shortenedUrl
	shortLink.ExpiresAt = expiresAt

	c.JSON(http.StatusOK, shortLink)
}

func redirectShortLink(c *gin.Context) {
	shortenedUrl := c.Param("shortenedUrl")

	var shortLink ShortLink
	err := db.Where("shortened_url = ?", shortenedUrl).First(&shortLink).Error
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if time.Now().After(shortLink.ExpiresAt) {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("link has expired"))
		return
	}

	fmt.Printf("Redirecting to ", shortLink.OriginalURL)
	c.Redirect(http.StatusPermanentRedirect, shortLink.OriginalURL)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortenedURL() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
