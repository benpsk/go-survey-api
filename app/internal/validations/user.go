package validations

import (
	"net/mail"
	"strings"

	"github.com/benpsk/go-survey-api/internal/models"
)

func StoreUser(user models.User) map[string]string {
	msg := make(map[string]string)
	if strings.TrimSpace(user.Name) == "" {
		msg["name"] = "Name is required"
	}
	if len(user.Name) > 100 {
		msg["name"] = "Name must not exceed 100 characters"
	}
	return validate(user, msg)
}

func UserLogin(user models.User) map[string]string {
	msg := make(map[string]string)
	return validate(user, msg)
}

func validate(user models.User, msg map[string]string) map[string]string {
	if strings.TrimSpace(user.Password) == "" {
		msg["password"] = "Password is required"
	}
	if len(user.Password) > 100 {
		msg["password"] = "Password must not exceed 100 characters"
	}
	if strings.TrimSpace(user.Email) == "" {
		msg["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(user.Email); err != nil {
		msg["email"] = "Invalid email format"
	}
	if len(user.Email) > 100 {
		msg["email"] = "Email must not exceed 100 characters"
	}
	return msg
}
