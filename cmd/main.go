package main

import (
	_ "fitness-club-energy/docs"
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/config"
	"fitness-club-energy/internal/handler"
	"fitness-club-energy/internal/repository"
	"fitness-club-energy/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Fitness Club Energy API
// @version 1.0
// @description API для фитнес-клуба Energy с PostgreSQL
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация БД
	if err := repository.InitDB(cfg); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer repository.CloseDB()

	redisClient := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redisClient.Ping(); err != nil {
		log.Printf("⚠️  Redis not available: %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}
	defer redisClient.Close()

	// Инициализация репозиториев
	db := repository.GetDB()
	userRepo := repository.NewUserRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	membershipService := service.NewMembershipService(membershipRepo)

	// Инициализация обработчиков
	userHandler := handler.NewUserHandler(userService)
	membershipHandler := handler.NewMembershipHandler(membershipService)

	// Настройка роутера
	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes
	handler.SetupRoutes(r, userHandler, membershipHandler)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ Фитнес-клуб Energy работает с PostgreSQL"))
	}).Methods("GET")

	log.Printf("🚀 Сервер запущен на :%s", cfg.ServerPort)
	log.Printf("📚 Swagger: http://localhost:%s/swagger/index.html", cfg.ServerPort)
	log.Printf("🌐 API: http://localhost:%s/api/v1/users", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal("❌ Ошибка сервера:", err)
	}
}
