package userResponse

import (
	"github.com/horlakz/api.secretariat_repository/dto"
	"github.com/horlakz/api.secretariat_repository/payload/response"
)

type LoginResponse struct {
	response.Response

	Data dto.LoginResponseDTO `json:"data"`
}
