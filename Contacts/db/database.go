package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Name     string `gorm:"not null" json:"name"`
	Phone    string `gorm:"not null" json:"phone"    binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required"`
}

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("contacts.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&User{})
}
