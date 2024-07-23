package models

import (
	"time"

	"gorm.io/gorm"
)

// Define a struct to represent a shortened URL.
type RedirectRequests struct {
	gorm.Model
	ShortenedURL string `json:"shortened_url"`
	RemoteAddr   string
	UrlHitTime   time.Time
}
