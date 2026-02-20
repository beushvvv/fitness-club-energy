package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fitness-club-energy/internal/dto/response"
	"fitness-club-energy/internal/model"
	"fitness-club-energy/internal/service"

	"github.com/gorilla/mux"
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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	responseMemberships := response.FromMemberships(memberships)
	json.NewEncoder(w).Encode(responseMemberships)
}

// GetMembershipByID godoc
// @Summary Получить абонемент по ID
// @Tags memberships
// @Produce json
// @Param id path int true "ID абонемента"
// @Success 200 {object} response.MembershipResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/memberships/{id} [get]
func (h *MembershipHandler) GetMembershipByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	membership, err := h.membershipService.GetMembershipByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Membership not found"})
		return
	}

	responseMembership := response.FromMembership(membership)
	json.NewEncoder(w).Encode(responseMembership)
}

// CreateMembership godoc
// @Summary Создать новый абонемент
// @Tags memberships
// @Accept json
// @Produce json
// @Param membership body model.Membership true "Данные абонемента"
// @Success 201 {object} response.MembershipResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/memberships [post]
func (h *MembershipHandler) CreateMembership(w http.ResponseWriter, r *http.Request) {
	var membership model.Membership
	if err := json.NewDecoder(r.Body).Decode(&membership); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.membershipService.CreateMembership(&membership); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.FromMembership(&membership))
}

// GetMembership godoc
// @Summary Получить абонемент по ID
// @Tags memberships
// @Produce json
// @Param id path int true "ID абонемента"
// @Success 200 {object} response.MembershipResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/memberships/{id} [get]
func (h *MembershipHandler) GetMembership(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	membership, err := h.membershipService.GetMembershipByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Membership not found"})
		return
	}
	json.NewEncoder(w).Encode(response.FromMembership(membership))
}

// UpdateMembership godoc
// @Summary Обновить абонемент
// @Tags memberships
// @Accept json
// @Produce json
// @Param id path int true "ID абонемента"
// @Param membership body model.Membership true "Данные абонемента"
// @Success 200 {object} response.MembershipResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/memberships/{id} [put]
func (h *MembershipHandler) UpdateMembership(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var membership model.Membership
	if err := json.NewDecoder(r.Body).Decode(&membership); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}
	membership.ID = id

	if err := h.membershipService.UpdateMembership(&membership); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(response.FromMembership(&membership))
}
