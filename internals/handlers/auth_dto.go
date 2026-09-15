package handlers

import (
	"net/mail"
	"strings"
	"time"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (req SignupRequest) Validate() error {
	if strings.TrimSpace(req.Name) == "" {
		return &ValidationError{
			Field: "name",
			Msg:   "must not be empty",
		}
	}

	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return &ValidationError{
			Field: "email",
			Msg:   "must be valid email address",
		}
	}

	if len(req.Password) < 8 {
		return &ValidationError{
			Field: "password",
			Msg:   "must be at least 8 characters"}
	}

	return nil
}
