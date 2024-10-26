package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/benpsk/go-survey-api/internal/models"
	"github.com/benpsk/go-survey-api/internal/validations"
)

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	userId, ok := Auth(r)
	if !ok {
		Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.Service.GetUserById(context.Background(), userId)
	if err != nil {
		Error(w, "User not found", http.StatusNotFound)
		return
	}
	Success(w, user, http.StatusOK)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	validationErrors := validations.StoreUser(user)
	if len(validationErrors) > 0 {
		ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	res, err := h.Service.Create(context.Background(), user)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Success(w, res, http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	validationErrors := validations.UserLogin(user)
	if len(validationErrors) > 0 {
		ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	res, err := h.Service.Login(context.Background(), user)
	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Success(w, res, http.StatusOK)
}
