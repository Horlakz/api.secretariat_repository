package dto

type UserDTO struct {
	DTO

	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phonenumber"`
	IsEmailVerified bool   `json:"isEmailVerified"`
	Password        string `json:"password"`
	Role            string `json:"role"`
}

type VerificationCodeDTO struct {
	DTO

	Code   string  `json:"code"`
	UserID string  `json:"user_id"`
	User   UserDTO `json:"user"`
}
