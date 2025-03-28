package main

import (
	"Contacts/db"
	"Contacts/internals/contacts"
	"Contacts/internals/middlewares"
	"Contacts/internals/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	// Users
	r.POST("/api/users/register", users.Register)
	r.POST("/api/users/login", users.Login)

	r.GET("/", middlewares.CheckAuth(), func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id не в контексте"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// Contacts
	r.POST("/api/contacts/create", middlewares.CheckAuth(), contacts.Create)
	r.GET("/api/contacts/all", middlewares.CheckAuth(), contacts.GetAll)

	r.Run()
}
