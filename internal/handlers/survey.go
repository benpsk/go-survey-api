package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benpsk/go-survey-api/internal/models"
	"github.com/benpsk/go-survey-api/internal/validations"
)

// api lists
// 1. create
// 2. listings by users
// 3. update

func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	var input models.SurveyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	validationErrors := validations.StoreSurvey(input)
	if len(validationErrors) > 0 {
		ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	userId, ok := Auth(r)
	if !ok {
		Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	input.UserId = userId
	res, err := h.Service.Store(context.Background(), input)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Success(w, res, http.StatusCreated)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userId, ok := Auth(r)
	if !ok {
		Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	res, err := h.Service.Get(context.Background(), userId)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Success(w, res, http.StatusCreated)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	userId, ok := Auth(r)
	if !ok {
		Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, "Invalid id", http.StatusUnauthorized)
	}
	res, err := h.Service.GetById(context.Background(), id)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.UserId != userId {
		Error(w, "Survey Unauthorized", http.StatusUnauthorized)
		return
	}
	Success(w, res, http.StatusCreated)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(w, "Invalid id", http.StatusUnauthorized)
	}

	var input models.SurveyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}
	validationErrors := validations.StoreSurvey(input)
	if len(validationErrors) > 0 {
		ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	userId, ok := Auth(r)
	if !ok {
		Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	survey, err := h.Service.Repo.GetById(ctx, id)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if survey.UserId != userId {
		Error(w, "Survey Unauthorized", http.StatusUnauthorized)
		return
	}
	input.UserId = userId
	res, err := h.Service.Update(ctx, id, input)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	Success(w, res, http.StatusOK)
}
