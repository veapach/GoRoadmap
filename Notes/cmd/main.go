package main

import (
	"Notes/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	r := gin.Default()

	r.POST("/user", func(c *gin.Context) {
		var user db.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
			return
		}

		if err := db.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно создан"})

	})

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
