package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"fitness-club-energy/internal/dto/request"
	"fitness-club-energy/internal/dto/response"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/service"
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
// @Success 200 {object} response.UserListResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	// Преобразование model → response DTO
	var userResponses []response.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, response.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	json.NewEncoder(w).Encode(response.UserListResponse{
		Users: userResponses,
		Total: len(userResponses),
	})
}

// GetUser godoc
// @Summary Получить пользователя по ID
// @Tags users
// @Param id path int true "ID пользователя"
// @Produce json
// @Success 200 {object} response.UserResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Пользователь не найден"})
		return
	}

	// Преобразование model → response DTO
	resp := response.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	json.NewEncoder(w).Encode(resp)
}

// CreateUser godoc
// @Summary Создать нового пользователя
// @Tags users
// @Accept json
// @Param request body request.CreateUserRequest true "Данные пользователя"
// @Produce json
// @Success 201 {object} response.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Преобразование request DTO → model
	user := model.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := h.userService.CreateUser(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	// Преобразование model → response DTO
	resp := response.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
