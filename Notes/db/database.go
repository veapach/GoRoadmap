package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Note struct {
	gorm.Model
	Title  string `gorm:"not null" json:"title"`
	Text   string `gorm:"not null" json:"text" binding:"required"`
	UserId uint   `json:"user_id"`
}

type User struct {
	gorm.Model
	Name     string `gorm:"not null" json:"name" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required"`
	Notes    []Note
}

func InitDB() {
	dsn, err := loadEnv()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	var err2 error
	DB, err2 = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err2 != nil {
		log.Fatal("Failed to connect to DB:", err2)
	}

	DB.AutoMigrate(&User{}, &Note{})

}

func loadEnv() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}

	dbName := os.Getenv("DB_NAME") + ".db"
	if dbName == "" {
		dbName = "notes.db"
	}

	return dbName, nil
}
