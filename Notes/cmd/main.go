package main

import (
	"Notes/db"
	"Notes/internal/notes"
	"Notes/internal/users"

	"github.com/gin-gonic/gin"
)

func main() {

	db.InitDB()

	r := gin.Default()

  // Users
	r.POST("/api/users/register", users.Register)
	r.POST("/api/users/login", users.Login)

  // Notes
  r.POST("/api/notes/create", users.AuthMiddleware(), notes.CreateNote)
  r.GET("/api/notes/get-all", users.AuthMiddleware(), notes.GetAllNotes)

	r.Run()
}
