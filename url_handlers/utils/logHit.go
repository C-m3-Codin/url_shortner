package utils

import (
	"fmt"
	"time"

	"github.com/c-m3-codin/url_shortner/models"
	"github.com/c-m3-codin/url_shortner/services"
)

func LogHit(redReq *models.ShortLink, remoteIP string) {

	fmt.Printf("Request Logging active")

	err := services.DB.Create(&models.RedirectRequests{
		RemoteAddr:   remoteIP,
		ShortenedURL: redReq.ShortenedURL,
		UrlHitTime:   time.Now(),
	}).Error

	if err != nil {
		panic("Couldnt log Request")

	}

}
