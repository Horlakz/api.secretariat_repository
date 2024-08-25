package dto

import "github.com/google/uuid"

type FileDTO struct {
	DTO

	Name     string    `json:"name"`
	Key      string    `json:"key"`
	MimeType string    `json:"mime_type"`
	Size     int64     `json:"size"`
	UserId   uuid.UUID `json:"user_id"`
}

type TransferDTO struct {
	DTO

	FileId     uuid.UUID `json:"file_id"`
	FromUserId uuid.UUID `json:"from_user_id"`
	ToUserId   uuid.UUID `json:"to_user_id"`
}
