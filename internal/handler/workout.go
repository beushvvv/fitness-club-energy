package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fitness-club-energy/internal/dto/request"
	"fitness-club-energy/internal/dto/response"
	"fitness-club-energy/internal/model"
)

// GetWorkouts godoc
// @Summary Получить тренировки с динамическими фильтрами
// @Description Получение тренировок с фильтрацией по типу, дате, пользователю
// @Tags workouts
// @Param type query string false "Тип тренировки"
// @Param date query string false "Дата тренировки (YYYY-MM-DD)"
// @Param user_id query int false "ID пользователя"
// @Produce json
// @Success 200 {array} response.WorkoutResponse
// @Router /api/v1/workouts [get]
func GetWorkouts(w http.ResponseWriter, r *http.Request) {
	// Создаем DTO из query параметров
	var req request.FilterWorkoutsRequest

	req.Type = r.URL.Query().Get("type")
	req.Date = r.URL.Query().Get("date")

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			req.UserID = userID
		}
	}

	// Динамическое формирование ответа на основе DTO
	workouts := getDynamicWorkouts(req)

	// Преобразование model → response DTO
	var workoutResponses []response.WorkoutResponse
	for _, workout := range workouts {
		workoutResponses = append(workoutResponses, response.WorkoutResponse{
			ID:              workout.ID,
			UserID:          workout.UserID,
			Type:            workout.Type,
			DurationMinutes: workout.DurationMinutes,
			CaloriesBurned:  workout.CaloriesBurned,
			Notes:           workout.Notes,
			WorkoutDate:     workout.WorkoutDate,
			CreatedAt:       workout.CreatedAt,
		})
	}

	json.NewEncoder(w).Encode(workoutResponses)
}

func getDynamicWorkouts(req request.FilterWorkoutsRequest) []model.Workout {
	// Mock данные (в реальности из БД)
	mockWorkouts := []model.Workout{
		{
			ID:              1,
			UserID:          1,
			Type:            "Йога",
			DurationMinutes: 60,
			CaloriesBurned:  300,
			Notes:           "Утренняя практика йоги",
			WorkoutDate:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			CreatedAt:       time.Now(),
		},
		{
			ID:              2,
			UserID:          2,
			Type:            "Кардио",
			DurationMinutes: 45,
			CaloriesBurned:  500,
			Notes:           "Интервальная тренировка",
			WorkoutDate:     time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC),
			CreatedAt:       time.Now(),
		},
		{
			ID:              3,
			UserID:          1,
			Type:            "Силовая",
			DurationMinutes: 90,
			CaloriesBurned:  600,
			Notes:           "Тренировка с весом",
			WorkoutDate:     time.Date(2024, 1, 17, 0, 0, 0, 0, time.UTC),
			CreatedAt:       time.Now(),
		},
	}

	// Динамическая фильтрация на основе DTO
	var filtered []model.Workout
	for _, w := range mockWorkouts {
		match := true

		if req.Type != "" && w.Type != req.Type {
			match = false
		}
		if req.Date != "" {
			dateStr := w.WorkoutDate.Format("2006-01-02")
			if dateStr != req.Date {
				match = false
			}
		}
		if req.UserID > 0 && w.UserID != req.UserID {
			match = false
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
// @Param request body request.CreateWorkoutRequest true "Динамические данные тренировки"
// @Produce json
// @Success 201 {object} response.WorkoutResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/workouts [post]
func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var req request.CreateWorkoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Валидация DTO
	if req.UserID <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "user_id должен быть положительным числом"})
		return
	}

	if req.DurationMinutes <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "duration_minutes должен быть положительным числом"})
		return
	}

	// Преобразование request DTO → model
	workout := model.Workout{
		UserID:          req.UserID,
		Type:            req.Type,
		DurationMinutes: req.DurationMinutes,
		CaloriesBurned:  req.CaloriesBurned,
		Notes:           req.Notes,
		WorkoutDate:     time.Now(),
		CreatedAt:       time.Now(),
	}

	// Здесь в реальности была бы логика сохранения в БД
	// Для примера присваиваем ID
	workout.ID = 100

	// Преобразование model → response DTO
	resp := response.WorkoutResponse{
		ID:              workout.ID,
		UserID:          workout.UserID,
		Type:            workout.Type,
		DurationMinutes: workout.DurationMinutes,
		CaloriesBurned:  workout.CaloriesBurned,
		Notes:           workout.Notes,
		WorkoutDate:     workout.WorkoutDate,
		CreatedAt:       workout.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
