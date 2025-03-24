package main

import (
	"Notes/db"
	"Notes/internal/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	r := gin.Default()

	r.POST("/api/users/register", users.Register)
	r.POST("/api/users/login", users.Login)

	r.GET("/users", users.AuthMiddleware(), func(c *gin.Context) {
		var users []db.User
		if err := db.DB.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
			return
		}

		c.JSON(http.StatusOK, users)
	})

	r.Run()
}
