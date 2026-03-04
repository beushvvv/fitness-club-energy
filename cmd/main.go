package main

import (
	// "log" // Удаляем стандартный лог
	"net/http"

	_ "fitness-club-energy/docs"
	"fitness-club-energy/internal/cache"
	"fitness-club-energy/internal/config"
	"fitness-club-energy/internal/handler"
	"fitness-club-energy/internal/logger" // Добавляем наш логгер
	"fitness-club-energy/internal/repository"
	"fitness-club-energy/internal/service"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	// Импортируем zap для использования сахара
)

// @title Fitness Club Energy API
// @version 1.0
// @description API для фитнес-клуба Energy с PostgreSQL и Redis
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// 1. Инициализация логгера
	if err := logger.Init("debug"); err != nil { // Для разработки ставим debug
		panic(err)
	}
	defer logger.Sync() // Важно! Сбрасываем буферы при выходе

	// Используем SugaredLogger для удобства (fmt.Printf стиль)
	sugar := logger.Log.Sugar()

	// Инициализация БД
	if err := repository.InitDB(cfg); err != nil {
		sugar.Fatalw("❌ Failed to connect to database", "error", err) // Логируем фатальную ошибку
	}
	defer repository.CloseDB()
	sugar.Info("✅ Connected to PostgreSQL") // Информационное сообщение

	// Инициализация Redis
	redisClient := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redisClient.Ping(); err != nil {
		sugar.Warnw("⚠️ Redis not available", "error", err) // Предупреждение, но не фатал
	} else {
		sugar.Info("✅ Connected to Redis")
	}
	defer redisClient.Close()

	// Создаём CacheWrapper
	cacheWrapper := cache.NewCacheWrapper(redisClient)

	// Инициализация репозиториев
	db := repository.GetDB()
	userRepo := repository.NewUserRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)
	workoutRepo := repository.NewWorkoutRepository(db)

	// Инициализация сервисов
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

	// Логируем успешный запуск
	sugar.Infow("🚀 Server started", "port", cfg.ServerPort)
	sugar.Infof("📚 Swagger: http://localhost:%s/swagger/index.html", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		sugar.Fatalw("❌ Server error", "error", err)
	}
}
