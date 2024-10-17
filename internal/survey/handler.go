package survey

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/benpsk/go-survey-api/pkg"
)

type SurveyHandler struct {
	Service *SurveyService
}

func NewSurveyHandler(service *SurveyService) *SurveyHandler {
	return &SurveyHandler{Service: service}
}

// api lists
// 1. create
// 2. listings by users
// 3. update

func (h *SurveyHandler) Store(w http.ResponseWriter, r *http.Request) {
	var input SurveyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkg.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	validationErrors := validateStore(input)
	if len(validationErrors) > 0 {
		pkg.ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	userId, ok := pkg.Auth(r)
	if !ok {
		pkg.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	input.UserId = userId
	res, err := h.Service.Store(context.Background(), input)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkg.Success(w, res, http.StatusCreated)
}

func (h *SurveyHandler) Get(w http.ResponseWriter, r *http.Request) {
	userId, ok := pkg.Auth(r)
	if !ok {
		pkg.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
	res, err := h.Service.Get(context.Background(), userId)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkg.Success(w, res, http.StatusCreated)
}

func (h *SurveyHandler) GetById(w http.ResponseWriter, r *http.Request) {
	userId, ok := pkg.Auth(r)
	if !ok {
		pkg.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, "Invalid id", http.StatusUnauthorized)
	}
	res, err := h.Service.GetById(context.Background(), id)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.UserId != userId {
		pkg.Error(w, "Survey Unauthorized", http.StatusUnauthorized)
		return
	}
	pkg.Success(w, res, http.StatusCreated)
}

func (h *SurveyHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, "Invalid id", http.StatusUnauthorized)
	}

	var input SurveyInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		pkg.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}
	validationErrors := validateStore(input)
	if len(validationErrors) > 0 {
		pkg.ValidationError(w, validationErrors, http.StatusBadRequest)
		return
	}
	userId, ok := pkg.Auth(r)
	if !ok {
		pkg.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ctx := context.Background()
	survey, err := h.Service.Repo.getById(ctx, id)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if survey.UserId != userId {
		pkg.Error(w, "Survey Unauthorized", http.StatusUnauthorized)
		return
	}
	input.UserId = userId
	res, err := h.Service.Update(ctx, id, input)
	if err != nil {
		pkg.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pkg.Success(w, res, http.StatusOK)

}
