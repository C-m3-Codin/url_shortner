package models

import (
	"time"

	"gorm.io/gorm"
)

// Define a struct to represent a shortened URL.
type ShortLink struct {
	gorm.Model
	OriginalURL   string    `json:"original_url"`
	ShortenedURL  string    `json:"shortened_url"`
	ExpiresAt     time.Time `json:"expires_at"`
	OwnerUsername string
}
