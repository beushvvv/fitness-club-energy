package handler

import (
	"encoding/json"
	"fitness-club-1/internal/models"
	"fitness-club-1/internal/service"
	"net/http"
	"time"
)

type MembershipHandler struct {
	membershipService *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{membershipService: membershipService}
}

// GetMemberships godoc
// @Summary Получить все абонементы
// @Tags memberships
// @Produce json
// @Success 200 {array} models.Membership
// @Router /api/v1/memberships [get]
func (h *MembershipHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	memberships, err := h.membershipService.GetAllMemberships()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(memberships)
}

// CreateMembership godoc
// @Summary Создать новый абонемент
// @Tags memberships
// @Accept json
// @Param request body models.CreateMembershipRequest true "Данные абонемента"
// @Produce json
// @Success 201 {object} models.Membership
// @Failure 400 {object} models.ErrorResponse
// @Router /api/v1/memberships [post]
func (h *MembershipHandler) CreateMembership(w http.ResponseWriter, r *http.Request) {
	var req models.CreateMembershipRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Парсим даты
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат start_date. Используйте YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат end_date. Используйте YYYY-MM-DD"})
		return
	}

	// Создаем membership
	membership := models.Membership{
		UserID:    req.UserID,
		Type:      req.Type,
		Price:     req.Price,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}

	if err := h.membershipService.CreateMembership(&membership); err != nil {
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(membership)
}
