package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"url_shortner/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// const (
// 	dbUser     = "root"
// 	dbPassword = "password"
// 	dbName     = "shortlinks"
// )

var db *gorm.DB

type ShortLink struct {
	ID           int       `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortenedURL string    `json:"shortened_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func main() {

	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbHost := os.Getenv("DB_HOST")
	// dbName := os.Getenv("DB_NAME")
	var err error
	// fmt.Sprintf("\n\n\\n", dbHost)
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:3305)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)
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
			err := db.Where("expires_at < ?", time.Now()).Delete(&ShortLink{}).Error
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Hour)
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
