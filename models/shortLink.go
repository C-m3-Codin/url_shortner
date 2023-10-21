package models

import (
	"time"

	"gorm.io/gorm"
)

// Define a struct to represent a shortened URL.

type ShortLink struct {
	gorm.Model
	OriginalURL   string    `json:"original_url,omitempty"`
	ShortenedURL  string    `json:"shortened_url,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	OwnerUsername string    `json:"OwnerUsername,omitempty"`
	UserID        uint      `gorm:"foreignKey:ID",json:"user_id,omitempty"`
}

type ShortLink_response struct {
	OriginalURL  string    `json:"original_url"`
	ShortenedURL string    `json:"shortened_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}
