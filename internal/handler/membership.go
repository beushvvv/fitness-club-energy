package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"fitness-club-energy/internal/dto/request"
	"fitness-club-energy/internal/dto/response"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/service"
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
// @Success 200 {array} response.MembershipResponse
// @Router /api/v1/memberships [get]
func (h *MembershipHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	memberships, err := h.membershipService.GetAllMemberships()
	if err != nil {
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	// Преобразование model → response DTO
	var membershipResponses []response.MembershipResponse
	for _, membership := range memberships {
		membershipResponses = append(membershipResponses, response.MembershipResponse{
			ID:        membership.ID,
			UserID:    membership.UserID,
			Type:      membership.Type,
			Price:     membership.Price,
			StartDate: membership.StartDate,
			EndDate:   membership.EndDate,
			IsActive:  membership.IsActive,
			CreatedAt: membership.CreatedAt,
		})
	}

	json.NewEncoder(w).Encode(membershipResponses)
}

// CreateMembership godoc
// @Summary Создать новый абонемент
// @Tags memberships
// @Accept json
// @Param request body request.CreateMembershipRequest true "Данные абонемента"
// @Produce json
// @Success 201 {object} response.MembershipResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/memberships [post]
func (h *MembershipHandler) CreateMembership(w http.ResponseWriter, r *http.Request) {
	var req request.CreateMembershipRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Парсим даты
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат start_date. Используйте YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат end_date. Используйте YYYY-MM-DD"})
		return
	}

	// Преобразование request DTO → model
	membership := model.Membership{
		UserID:    req.UserID,
		Type:      req.Type,
		Price:     req.Price,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}

	if err := h.membershipService.CreateMembership(&membership); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	// Преобразование model → response DTO
	resp := response.MembershipResponse{
		ID:        membership.ID,
		UserID:    membership.UserID,
		Type:      membership.Type,
		Price:     membership.Price,
		StartDate: membership.StartDate,
		EndDate:   membership.EndDate,
		IsActive:  membership.IsActive,
		CreatedAt: membership.CreatedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
