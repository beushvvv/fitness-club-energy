package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

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

// GetMemberships - получение всех абонементов
func (h *MembershipHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	memberships, err := h.membershipService.GetAllMemberships()
	if err != nil {
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(memberships)
}

// GetMembership - получение абонемента по ID (НОВЫЙ МЕТОД)
func (h *MembershipHandler) GetMembership(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	membership, err := h.membershipService.GetMembershipByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Абонемент не найден"})
		return
	}

	json.NewEncoder(w).Encode(membership)
}

// CreateMembership - создание абонемента
func (h *MembershipHandler) CreateMembership(w http.ResponseWriter, r *http.Request) {
	var req request.CreateMembershipRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	membership := &model.Membership{
		UserID:    req.UserID,
		Type:      req.Type,
		Price:     req.Price,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}

	if err := h.membershipService.CreateMembership(membership); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(membership)
}

// UpdateMembership - обновление абонемента (НОВЫЙ МЕТОД)
func (h *MembershipHandler) UpdateMembership(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var req request.CreateMembershipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	membership := &model.Membership{
		ID:        id,
		UserID:    req.UserID,
		Type:      req.Type,
		Price:     req.Price,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  true,
	}

	if err := h.membershipService.UpdateMembership(membership); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.ErrorResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(membership)
}
