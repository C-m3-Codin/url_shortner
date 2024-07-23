package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/c-m3-codin/url_shortner/models"
	jwt "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

// Claims represents the JWT claims.
type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserId   uint   `json:""`
	exp      time.Time
	IsAdmin  bool
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for the given username and email.
func GenerateJWT(user models.User) (string, error) {
	// Set the expiration time of the token to 1 hour from now.
	expirationTime := time.Now().Add(time.Hour * 1)

	//fmt.Println"Time now ", time.Now(), "\n Time of expiry ", expirationTime)

	// Create the JWT claims, which includes the username, email, and expiration time.
	claims := &Claims{
		UserId:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		exp:      expirationTime,
		IsAdmin:  user.IsAdmin,
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
func ValidateToken(tokenString string) (error, *Claims) {
	// Parse the token using the HS256 algorithm and the secret key.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	// Check if there was an error parsing the token.
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return errors.New("invalid token signature"), nil
		}
		return err, nil
	}

	// Check if the token is valid.
	if !token.Valid {
		return errors.New("invalid token"), nil
	}

	claims := token.Claims.(*Claims)
	//fmt.Println"claims is \n", claims.Username, claims.StandardClaims.IssuedAt)

	//fmt.Println"Time now ", time.Now(), "\n token is  ", token.Header)

	err = token.Claims.Valid()
	if err != nil {
		fmt.Printf("Token not valid")
		return err, nil
	}

	return nil, claims
}
