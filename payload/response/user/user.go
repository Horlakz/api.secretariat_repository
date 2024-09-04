package userResponse

import "github.com/horlakz/api.secretariat_repository/payload/response"

type UserResponse struct {
	response.Response
	Data UserResponseData `json:"data"`
}

type UserResponseData struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
}
