package services

import (
	"fmt"
	"os"

	"github.com/c-m3-codin/url_shortner/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// NewDatabase : intializes and returns mysql db
func NewDatabase() (*gorm.DB, error) {

	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")
	ENV := os.Getenv("ENVIRONMENT")

	if ENV == "" {
		fmt.Println("No Env")
		USER = "user"
		PASS = "Password@123"
		HOST = "localhost"
		// PORT = 3306
		DBNAME = "golang_url_shortner"
	}

	// fmt.Println(env)
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS,
		HOST, DBNAME)
	fmt.Println(URL)
	db, err := gorm.Open(mysql.Open(URL))

	// Create the necessary tables in the database.
	db.AutoMigrate(models.ShortLink{})
	db.AutoMigrate(models.RedirectRequests{})
	db.AutoMigrate(models.User{})
	DB = db
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
