package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"fitness-club-energy/internal/dto/request"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/service"

	"github.com/gorilla/mux"
)

type WorkoutHandler struct {
	workoutService *service.WorkoutService
}

func NewWorkoutHandler(workoutService *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{workoutService: workoutService}
}

// GetWorkouts godoc
// @Summary Получить все тренировки
// @Tags workouts
// @Produce json
// @Success 200 {array} model.Workout
// @Router /api/v1/workouts [get]
func (h *WorkoutHandler) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.workoutService.GetAllWorkouts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(workouts)
}

// GetWorkoutByID godoc
// @Summary Получить тренировку по ID
// @Tags workouts
// @Produce json
// @Param id path int true "ID тренировки"
// @Success 200 {object} model.Workout
// @Failure 404 {object} map[string]string
// @Router /api/v1/workouts/{id} [get]
func (h *WorkoutHandler) GetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	workout, err := h.workoutService.GetWorkoutByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Workout not found"})
		return
	}
	json.NewEncoder(w).Encode(workout)
}

// CreateWorkout godoc
// @Summary Создать новую тренировку
// @Tags workouts
// @Accept json
// @Produce json
// @Param workout body model.Workout true "Данные тренировки"
// @Success 201 {object} model.Workout
// @Failure 400 {object} map[string]string
// @Router /api/v1/workouts [post]
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout model.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.workoutService.CreateWorkout(&workout); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
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
