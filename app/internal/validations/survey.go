package validations

import (
	"time"

	"github.com/benpsk/go-survey-api/internal/models"
)

func StoreSurvey(survey models.SurveyInput) map[string]string {
	msg := make(map[string]string)
	if survey.Name == "" {
		msg["name"] = "Name is required"
	}
	if len(survey.Name) > 100 {
		msg["name"] = "Name cannot exceed 100 characters"
	}
	if survey.PhoneNo == "" {
		msg["phone_no"] = "phone number is required"
	}
	if len(survey.PhoneNo) > 100 {
		msg["phone_no"] = "phone number cannot exceed 100 characters"
	}
	if survey.Gender != "Male" && survey.Gender != "Female" {
		msg["gender"] = "gender must be either 'Male' or 'Female'"
	}

	if survey.Dob == "" {
		msg["dob"] = "dob is required"
	}
	if survey.Dob != "" {
		if !isValidDate(survey.Dob) {
			msg["dob"] = "dob must be in the format YYYY-MM-DD"
		}
	}
	return msg
}

func isValidDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
