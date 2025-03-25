package notes

import (
	"Notes/db"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateNote(c *gin.Context){
  var note db.Note
  
  userID, exists := c.Get("user_id")
  if !exists{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id не найден"})
    return
  }

  uid, ok := userID.(uint)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
    return
  }

  if err := c.ShouldBindJSON(&note); err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
    return
  }

  note.UserId = uid

  if note.Title == ""{
    runes := []rune(note.Text)
    if len(runes) > 5{
      note.Title = string(runes[:5]) + "..."
    } else {
      note.Title = note.Text
    }
  }

  if err := db.DB.Create(&note).Error; err != nil{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании заметки в БД"})
    return
  }

  c.JSON(http.StatusOK, gin.H{"message": "Заметка успешно создана"}) 
}

func GetAllNotes(c *gin.Context){

  userID, exists := c.Get("user_id")
  if !exists{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id не найден"})
    return
  }

  uid, ok := userID.(uint)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
    return
  }

  var notes []db.Note

  if err := db.DB.Where("user_id = ?", uid).Find(&notes).Error; err != nil{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске заметок"})
    return
  }

  if len(notes) == 0{
    c.JSON(http.StatusOK, gin.H{"notes": notes, "count": 0})
    return

  }

  c.JSON(http.StatusOK, gin.H{"notes": notes, "count": len(notes)})
}

func GetNoteByID(c *gin.Context){

  noteID := c.Param("note_id")

  userID, exists := c.Get("user_id")
  if !exists{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id не найден"})
    return
  }

  uid, ok := userID.(uint)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
    return
  }

  var note db.Note

  if err := db.DB.Where("user_id = ? AND id = ?", uid, noteID).Find(&note).Error; err != nil{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске заметки"})
    return
  }

  if note.ID == 0{
    c.JSON(http.StatusNotFound, gin.H{"error": "Заметка с таким id не найдена"})
    return
  }

  c.JSON(http.StatusOK, note)
}

func DeleteNoteByID(c *gin.Context){
  noteID := c.Param("note_id")

  
  userID, exists := c.Get("user_id")
  if !exists{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id не найден"})
    return
  }

  uid, ok := userID.(uint)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
    return
  }

  var note db.Note

  if err := db.DB.Where("user_id = ? AND id = ?", uid, noteID).First(&note).Error; err != nil{
    if errors.Is(err, gorm.ErrRecordNotFound){
      c.JSON(http.StatusNotFound, gin.H{"error": "Заметка не найдена"})
      return
    } else {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске заметки"})
    }
    return
  }

  if err := db.DB.Unscoped().Where("user_id = ? AND id = ?", uid, noteID).Delete(&note).Error; err != nil{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении заметки"})
    return
  }

  c.JSON(http.StatusOK, gin.H{"message": "Заметка успешно удалена"})
}

func UpdateNoteByID(c *gin.Context){
  noteID := c.Param("note_id")

  
  userID, exists := c.Get("user_id")
  if !exists{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id не найден"})
    return
  }

  uid, ok := userID.(uint)
  if !ok {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
    return
  }

  var existingNote db.Note
  if err := db.DB.Where("user_id = ? AND id = ?", uid, noteID).First(&existingNote).Error; err != nil{
    if errors.Is(err, gorm.ErrRecordNotFound){
      c.JSON(http.StatusNotFound, gin.H{"error": "Заметка не найдена"})
    } else{
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске заметки"})
    }
    return
  }

  var updatedData db.Note
  if err := c.ShouldBindJSON(&updatedData); err != nil{
    c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
    return
  }

  if err := db.DB.Model(&existingNote).Omit("id", "user_id").Updates(updatedData).Error; err != nil{
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении заметки"})
    return
  }




  c.JSON(http.StatusOK, gin.H{"message": "Заметка успешно обновлена", "note": existingNote})
}













