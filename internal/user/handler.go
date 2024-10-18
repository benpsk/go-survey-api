package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/benpsk/go-survey-api/pkg"
)

type UserHandler struct {
	Service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) User(w http.ResponseWriter, r *http.Request) {
	userId, ok := pkg.Auth(r)
	if !ok {
		pkg.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.Service.GetUserById(context.Background(), userId)
	if err != nil {
		pkg.Error(w, "User not found", http.StatusNotFound)
		return
	}
	pkg.Success(w, user, http.StatusOK)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		pkg.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	validationErrors := validateStore(user)
	if len(validationErrors) > 0 {
		pkg.ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	res, err := h.Service.Create(context.Background(), user)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkg.Success(w, res, http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		pkg.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	validationErrors := validateLogin(user)
	if len(validationErrors) > 0 {
		pkg.ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	res, err := h.Service.Login(context.Background(), user)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pkg.Success(w, res, http.StatusOK)
}
