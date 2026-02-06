package handler

import (
	"encoding/json"
	"fitness-club-1/internal/models"
	"net/http"
	"strconv"
	"time"
)

// GetWorkouts godoc
// @Summary Получить тренировки с динамическими фильтрами
// @Description Получение тренировок с фильтрацией по типу, дате, пользователю
// @Tags workouts
// @Param type query string false "Тип тренировки"
// @Param date query string false "Дата тренировки (YYYY-MM-DD)"
// @Param user_id query int false "ID пользователя"
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/v1/workouts [get]
func GetWorkouts(w http.ResponseWriter, r *http.Request) {
	// Получение динамических параметров от клиента
	workoutType := r.URL.Query().Get("type")
	date := r.URL.Query().Get("date")
	userIDStr := r.URL.Query().Get("user_id")

	// Динамическое формирование ответа
	workouts := getDynamicWorkouts(workoutType, date, userIDStr)

	json.NewEncoder(w).Encode(workouts)
}

func getDynamicWorkouts(workoutType, date, userIDStr string) []map[string]interface{} {
	// Mock данные (в реальности из БД)
	workouts := []map[string]interface{}{
		{
			"id":          1,
			"type":        "Йога",
			"date":        "2024-01-15",
			"user_id":     1,
			"duration":    60,
			"calories":    300,
			"trainer":     "Анна",
			"description": "Утренняя практика йоги",
		},
		{
			"id":          2,
			"type":        "Кардио",
			"date":        "2024-01-16",
			"user_id":     2,
			"duration":    45,
			"calories":    500,
			"trainer":     "Максим",
			"description": "Интервальная тренировка",
		},
		{
			"id":          3,
			"type":        "Силовая",
			"date":        "2024-01-17",
			"user_id":     1,
			"duration":    90,
			"calories":    600,
			"trainer":     "Иван",
			"description": "Тренировка с весом",
		},
		{
			"id":          4,
			"type":        "Йога",
			"date":        "2024-01-18",
			"user_id":     3,
			"duration":    75,
			"calories":    350,
			"trainer":     "Анна",
			"description": "Вечерняя релаксация",
		},
	}

	// Динамическая фильтрация
	var filtered []map[string]interface{}
	for _, w := range workouts {
		match := true

		if workoutType != "" && w["type"] != workoutType {
			match = false
		}
		if date != "" && w["date"] != date {
			match = false
		}
		if userIDStr != "" {
			userID, err := strconv.Atoi(userIDStr)
			if err == nil && w["user_id"] != userID {
				match = false
			}
		}

		if match {
			filtered = append(filtered, w)
		}
	}

	return filtered
}

// CreateWorkout godoc
// @Summary Создать тренировку с динамическими данными
// @Description Принимает динамические данные о тренировке от клиента
// @Tags workouts
// @Accept json
// @Param request body models.WorkoutRequest true "Динамические данные тренировки"
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/workouts [post]
func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var req models.WorkoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Валидация динамических данных
	if req.UserID <= 0 {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "user_id должен быть положительным числом"})
		return
	}

	if req.DurationMinutes <= 0 {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "duration_minutes должен быть положительным числом"})
		return
	}

	// Обработка динамических данных
	workout := map[string]interface{}{
		"id":               100, // В реальности из БД
		"user_id":          req.UserID,
		"type":             req.Type,
		"duration_minutes": req.DurationMinutes,
		"calories_burned":  req.CaloriesBurned,
		"notes":            req.Notes,
		"workout_date":     time.Now().Format("2006-01-02"),
		"created_at":       time.Now().Format(time.RFC3339),
		"message":          "Тренировка успешно создана",
		"dynamic_fields": map[string]interface{}{
			"estimated_calories": req.DurationMinutes * 8, // Пример динамического расчета
			"intensity_level":    "medium",
			"recommended_next":   getRecommendedNextWorkout(req.Type),
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// Вспомогательная функция для динамических рекомендаций
func getRecommendedNextWorkout(currentType string) string {
	switch currentType {
	case "Йога":
		return "Кардио"
	case "Кардио":
		return "Силовая"
	case "Силовая":
		return "Йога"
	default:
		return "Кардио"
	}
}
