package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// get the secret key from the env
var jwtKeys = []byte("secretKEy")

// Generates JWT token for a given username and email expiry set to 1 hr
func GenerateJWT(email string, username string) (token string, err error) {
	expirationTime := time.Now().Add(time.Hour * 1)

	claims := &jwt.MapClaims{
		"email":    email,
		"username": username,
		"issuer":   "cp",
		"exp":      expirationTime,
		"data": map[string]string{
			"name": "JohnDoe",
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtKeys))

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// token, err = token.SignedString(jwtKey)
	return
}

// checks the validity of a signed token returns nil if valid
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKeys), nil
	})

	if err != nil {
		return
	}
	claims := token.Claims.(jwt.MapClaims)

	exp, err := claims.GetExpirationTime()
	if err != nil {
		return
	}

	if exp.Time.Local().Unix() < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return
}
