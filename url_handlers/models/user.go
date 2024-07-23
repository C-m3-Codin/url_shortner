package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string `json:"name"`
	Username   string `json:"username" gorm:"unique"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"isAdmin" gorm:"defult:false"`
	ShortLinks []ShortLink
}

// hashes the password given plain text and updates the user.password to hashed password
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// checks password match - returns nil if match and error if not
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
