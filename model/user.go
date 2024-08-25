package model

import (
	"github.com/google/uuid"

	"github.com/horlakz/api.secretariat_repository/lib/database"
)

type User struct {
	database.BaseModel

	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Password        string `json:"password"`
	Role            string `json:"role"`
}

type VerificationCode struct {
	database.BaseModel

	UserID  uuid.UUID `json:"user_id"`
	Code    string    `json:"code"`
	Purpose string    `json:"purpose"`
}
