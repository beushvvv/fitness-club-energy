package handler

import (
	"encoding/json"
	"fitness-club-1/internal/models"
	"fitness-club-1/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers godoc
// @Summary Получить всех пользователей
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// GetUser godoc
// @Summary Получить пользователя по ID
// @Tags users
// @Param id path int true "ID пользователя"
// @Produce json
// @Success 200 {object} models.User
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// CreateUser godoc
// @Summary Создать нового пользователя
// @Tags users
// @Accept json
// @Param request body models.CreateUserRequest true "Данные пользователя"
// @Produce json
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Преобразование DTO в Entity
	user := models.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := h.userService.CreateUser(&user); err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	// Преобразование Entity в Response DTO
	response := models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
