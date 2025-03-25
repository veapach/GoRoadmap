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
  r.GET("/api/notes/all", users.AuthMiddleware(), notes.GetAllNotes)
  r.GET("/api/notes/:note_id", users.AuthMiddleware(), notes.GetNoteByID)
  r.DELETE("/api/notes/delete/:note_id", users.AuthMiddleware(), notes.DeleteNoteByID)
  r.PUT("/api/notes/update/:note_id", users.AuthMiddleware(), notes.UpdateNoteByID)

	r.Run()
}
