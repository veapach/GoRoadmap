package users

import (
	"Notes/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user db.User) (string, error) {
	expTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("supersecretkey"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Register(c *gin.Context) {

	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	user.Password, _ = HashPassword(user.Password)

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно создан",
		"token":   token,
	})

}

func Login(c *gin.Context) {

	var inputUser db.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	var dbUser db.User
	if err := db.DB.Where("name = ?", inputUser.Name).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	if !CheckPasswordHash(inputUser.Password, dbUser.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	token, err := GenerateToken(dbUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Успешный вход",
		"token":   token,
	})

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
