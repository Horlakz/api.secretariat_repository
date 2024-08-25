package userResponse

import "github.com/horlakz/api.secretariat_repository/payload/response"

type UserResponse struct {
	response.Response
	Data UserResponseData `json:"data"`
}

type UserResponseData struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	ReferralCode string `json:"referral_code"`
}
