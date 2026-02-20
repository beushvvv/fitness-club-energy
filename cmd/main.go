package main

import (
	"log"
	"net/http"

	_ "fitness-club-energy/docs"
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/config"
	"fitness-club-energy/internal/handler"
	"fitness-club-energy/internal/repository"
	"fitness-club-energy/internal/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Fitness Club Energy API
// @version 1.0
// @description API для фитнес-клуба Energy с PostgreSQL и Redis
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
	log.Println("✅ Connected to PostgreSQL")

	// Инициализация Redis
	redisClient := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redisClient.Ping(); err != nil {
		log.Printf("⚠️ Redis not available: %v", err)
	} else {
		log.Println("✅ Connected to Redis")
	}
	defer redisClient.Close()

	// Создаём CacheWrapper (один экземпляр для всех сервисов)
	cacheWrapper := cache.NewCacheWrapper(redisClient)

	// Инициализация репозиториев
	db := repository.GetDB()
	userRepo := repository.NewUserRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	workoutRepo := repository.NewWorkoutRepository(db) // ← теперь существует

	// Инициализация сервисов с CacheWrapper
	userService := service.NewUserService(userRepo, cacheWrapper)
	membershipService := service.NewMembershipService(membershipRepo, cacheWrapper)
	workoutService := service.NewWorkoutService(workoutRepo, cacheWrapper)

	// Инициализация обработчиков
	userHandler := handler.NewUserHandler(userService)
	membershipHandler := handler.NewMembershipHandler(membershipService)
	workoutHandler := handler.NewWorkoutHandler(workoutService)

	// Настройка роутера
	r := mux.NewRouter()

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")

	api.HandleFunc("/memberships", membershipHandler.GetMemberships).Methods("GET")
	api.HandleFunc("/memberships/{id}", membershipHandler.GetMembershipByID).Methods("GET")
	api.HandleFunc("/memberships", membershipHandler.CreateMembership).Methods("POST")

	api.HandleFunc("/workouts", workoutHandler.GetWorkouts).Methods("GET")
	api.HandleFunc("/workouts/{id}", workoutHandler.GetWorkoutByID).Methods("GET")
	api.HandleFunc("/workouts", workoutHandler.CreateWorkout).Methods("POST")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ Fitness Club Energy API is running"))
	}).Methods("GET")

	log.Printf("🚀 Server started on :%s", cfg.ServerPort)
	log.Printf("📚 Swagger: http://localhost:%s/swagger/index.html", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatal("❌ Server error:", err)
	}
}
