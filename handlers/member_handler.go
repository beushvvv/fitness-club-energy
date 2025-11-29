package handlers

import (
	"fitness-club/database"
	"fitness-club/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateMember(c *gin.Context) {
	var input struct {
		FirstName   string `json:"first_name" binding:"required"`
		LastName    string `json:"last_name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Phone       string `json:"phone"`
		DateOfBirth string `json:"date_of_birth"` // Принимаем как строку
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, нет ли уже пользователя с таким email
	var existingMember models.Member
	if err := database.DB.Where("email = ?", input.Email).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Member with this email already exists"})
		return
	}

	// Парсим дату рождения (если предоставлена)
	var dob time.Time
	if input.DateOfBirth != "" {
		parsedDob, err := time.Parse(time.RFC3339, input.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 format"})
			return
		}
		dob = parsedDob
	}

	member := models.Member{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		Phone:       input.Phone,
		DateOfBirth: dob,
	}

	result := database.DB.Create(&member)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func GetMembers(c *gin.Context) {
	var members []models.Member
	result := database.DB.Find(&members)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func GetMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member

	result := database.DB.First(&member, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	c.JSON(http.StatusOK, member)
}

func UpdateMember(c *gin.Context) {
	id := c.Param("id")
	var member models.Member

	// Проверяем существование пользователя
	if err := database.DB.First(&member, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email" binding:"omitempty,email"`
		Phone       string `json:"phone"`
		DateOfBirth string `json:"date_of_birth"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем поля если они предоставлены
	if input.FirstName != "" {
		member.FirstName = input.FirstName
	}
	if input.LastName != "" {
		member.LastName = input.LastName
	}
	if input.Email != "" {
		member.Email = input.Email
	}
	if input.Phone != "" {
		member.Phone = input.Phone
	}
	if input.DateOfBirth != "" {
		parsedDob, err := time.Parse(time.RFC3339, input.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use RFC3339 format"})
			return
		}
		member.DateOfBirth = parsedDob
	}

	result := database.DB.Save(&member)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

func DeleteMember(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Member{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member deleted successfully"})
}

// ResetMembers - только для тестирования! Удаляет всех членов
func ResetMembers(c *gin.Context) {
	if err := database.DB.Exec("DELETE FROM members").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All members deleted"})
}
