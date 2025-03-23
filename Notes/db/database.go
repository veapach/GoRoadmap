package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Note struct {
	gorm.Model
	Title  string `json:"title"`
	UserId uint   `json:"user_id"`
}

type User struct {
	gorm.Model
	Name  string `gorm:"not null" json:"name" binding:"required"`
	Notes []Note
}

func InitDB() {
	dsn, err := loadEnv()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	var err2 error
	DB, err2 = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err2 != nil {
		log.Fatal("Failed to connect to DB:", err2)
	}

	DB.AutoMigrate(&User{}, &Note{})

}

func loadEnv() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	host := os.Getenv("HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, pass, dbname, port)

	return dsn, nil

}
