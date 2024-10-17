package user

import (
	"net/mail"
	"strings"
)

func validateStore(user User) map[string]string {
	msg := make(map[string]string)
	if strings.TrimSpace(user.Password) == "" {
		msg["password"] = "Password is required"
	}
	if strings.TrimSpace(user.Email) == "" {
		msg["email"] = "Email is required"
	} else if _, err := mail.ParseAddress(user.Email); err != nil {
		msg["email"] = "Invalid email format"
	}
	return msg
}
