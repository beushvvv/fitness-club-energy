package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	ServerPort    string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func Load() *Config {
	// Загружаем .env файл
	godotenv.Load(".env")

	return &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres123"),
		DBName:        getEnv("DB_NAME", "fitness_club"),
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", "0"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key, defaultValue string) int {
	val := getEnv(key, defaultValue)
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return intVal
}
