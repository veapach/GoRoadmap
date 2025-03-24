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

	r.POST("/user", users.Register) 

	r.GET("/users", func(c *gin.Context) {
		var users []db.User
		if err := db.DB.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
			return
		}

		c.JSON(http.StatusOK, users)
	})

	r.Run()
}
