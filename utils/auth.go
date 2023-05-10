package utils

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

// Claims represents the JWT claims.
type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for the given username and email.
func GenerateJWT(email string, username string) (string, error) {
	// Set the expiration time of the token to 1 hour from now.
	expirationTime := time.Now().Add(time.Hour)

	// Create the JWT claims, which includes the username, email, and expiration time.
	claims := &Claims{
		Username: username,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "my_app",
			Subject:   "auth",
		},
	}

	// Create the token using the HS256 algorithm and the JWT claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key and return the signed token as a string.
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateToken validates the given JWT token and returns an error if the token is invalid.
func ValidateToken(tokenString string) error {
	// Parse the token using the HS256 algorithm and the secret key.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("\n\n error in sign method \n\n")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	// Check if there was an error parsing the token.
	if err != nil {
		fmt.Println("\n\n error parsing \n\n")
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("\n\n invalid token signature  \n\n")
			return errors.New("invalid token signature")
		}
		return err
	}

	// Check if the token is valid.
	if !token.Valid {
		fmt.Println("\n\n token not valid \n\n")
		return errors.New("invalid token")
	}

	fmt.Println("\n\n all well \n\n")

	return nil
}
