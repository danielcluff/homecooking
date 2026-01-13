package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/homecooking/backend/internal/services"
)

type UserInviteHandler struct {
	inviteService *services.UserInviteService
}

func NewUserInviteHandler(inviteService *services.UserInviteService) *UserInviteHandler {
	return &UserInviteHandler{
		inviteService: inviteService,
	}
}

type CreateUserInviteRequest struct {
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (h *UserInviteHandler) CreateInvite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req CreateUserInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	invite, err := h.inviteService.CreateInvite(req.Email, req.Role, userID, req.ExpiresAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invite)
}

func (h *UserInviteHandler) GetInvite(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		http.Error(w, "Code required", http.StatusBadRequest)
		return
	}

	invite, err := h.inviteService.GetInvite(code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invite)
}

func (h *UserInviteHandler) ListInvites(w http.ResponseWriter, r *http.Request) {
	invites, err := h.inviteService.ListInvites()
	if err != nil {
		http.Error(w, "Failed to fetch invites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invites)
}

func (h *UserInviteHandler) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "ID required", http.StatusBadRequest)
		return
	}

	if err := h.inviteService.DeleteInvite(id); err != nil {
		http.Error(w, "Failed to delete invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UseInviteRequest struct {
	Code string `json:"code"`
}

func (h *UserInviteHandler) UseInvite(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req UseInviteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	invite, err := h.inviteService.UseInvite(req.Code, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invite)
}
