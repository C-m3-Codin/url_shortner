package utils

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortenedURL() string {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 8)
	for i := range b {
		// Pick a random character from the character set and add it to the byte slice
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
