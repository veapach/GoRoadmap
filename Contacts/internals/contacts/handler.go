package contacts

import (
	"Contacts/db"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
	var input struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone" binding:"required"`
	}

	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return

	}

	var contactUser db.User
	if err := db.DB.Where("phone = ?", input.Phone).First(&contactUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь с таким номер не найден"})
		return
	}

	var existingContact db.Contact
	result := db.DB.Where("user_id = ? AND phone = ?", currentUserID, input.Phone).
		First(&existingContact)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Контакт уже существует"})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	newContact := db.Contact{
		Name:   input.Name,
		Phone:  input.Phone,
		UserID: currentUserID.(uint),
	}

	if err := db.DB.Create(&newContact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания контакта"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Контакт успешно добавлен",
			"contact": newContact,
		},
	)
}

func GetAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}

	var contacts db.Contact
	if err := db.DB.Where("user_id = ?", userID).Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске контактов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}
