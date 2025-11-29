package main

import (
	"fitness-club/database"
	"fitness-club/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к БД
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграции
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Создание роутера
	router := gin.Default()

	// Корневой маршрут
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Fitness Club API is running!",
			"endpoints": map[string]string{
				"GET /members":        "Get all members",
				"POST /members":       "Create new member",
				"GET /members/:id":    "Get member by ID",
				"PUT /members/:id":    "Update member",
				"DELETE /members/:id": "Delete member",
			},
		})
	})

	// Маршруты для членов
	router.POST("/members", handlers.CreateMember)
	router.GET("/members", handlers.GetMembers)
	router.GET("/members/:id", handlers.GetMember)
	router.PUT("/members/:id", handlers.UpdateMember)
	router.DELETE("/members/:id", handlers.DeleteMember)
	// Только для тестирования - удалить в продакшене
	router.DELETE("/members/reset", handlers.ResetMembers)

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
