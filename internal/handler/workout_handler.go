package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fitness-club-energy/internal/dto/request"
	"fitness-club-energy/internal/dto/response"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/service"
)

type WorkoutHandler struct {
	workoutService *service.WorkoutService
}

func NewWorkoutHandler(workoutService *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{workoutService: workoutService}
}

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
func (h *WorkoutHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	// Собираем фильтры из запроса
	filters := make(map[string]interface{})

	if workoutType := r.URL.Query().Get("type"); workoutType != "" {
		filters["type"] = workoutType
	}
	if date := r.URL.Query().Get("date"); date != "" {
		filters["date"] = date
	}
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			filters["user_id"] = userID
		}
	}

	// Получаем данные через сервис с кешированием
	workouts, err := h.workoutService.GetWorkoutsWithFilters(filters)
	if err != nil {
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(workouts)
}

// CreateWorkout godoc
// @Summary Создать тренировку
// @Tags workouts
// @Accept json
// @Param request body request.CreateWorkoutRequest true "Данные тренировки"
// @Produce json
// @Success 201 {object} response.WorkoutResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/workouts [post]
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var req request.CreateWorkoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Валидация
	if req.UserID <= 0 || req.DurationMinutes <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверные параметры"})
		return
	}

	// Создаём модель
	workout := &model.Workout{
		UserID:          req.UserID,
		Type:            req.Type,
		DurationMinutes: req.DurationMinutes,
		CaloriesBurned:  req.CaloriesBurned,
		Notes:           req.Notes,
		WorkoutDate:     time.Now(),
		CreatedAt:       time.Now(),
	}

	// Сохраняем
	if err := h.workoutService.CreateWorkout(workout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}
