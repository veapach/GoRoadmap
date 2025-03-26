package main

import (
	"Contacts/db"
	"Contacts/internals/users"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	r := gin.Default()

	// Users
	r.POST("/api/users/register", users.Register)
	r.POST("/api/users/login", users.Login)

	r.Run()
}
