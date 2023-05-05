package services

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDatabase : intializes and returns mysql db
func NewDatabase() (*gorm.DB, error) {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	fmt.Sprint(USER)
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS,
		HOST, DBNAME)
	fmt.Println(URL)
	db, err := gorm.Open(mysql.Open(URL))

	if err != nil {
		// panic("Failed to connect to database!")
		return nil, err

	}
	fmt.Println("Database connection established")
	return db, nil

}

// USER := os.Getenv("DB_USER")
// PASS := os.Getenv("DB_PASSWORD")
// HOST := os.Getenv("DB_HOST")
// DBNAME := os.Getenv("DB_NAME")
