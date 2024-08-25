package model

import (
	"github.com/google/uuid"
	"github.com/horlakz/api.secretariat_repository/lib/database"
)

type File struct {
	database.BaseModel

	Name     string    `json:"name"`
	Key      string    `json:"key"`
	MimeType string    `json:"mime_type"`
	Size     int64     `json:"size"`
	UserId   uuid.UUID `json:"user_id"`
}

type Transfer struct {
	database.BaseModel

	FileId     uuid.UUID `json:"file_id"`
	FromUserId uuid.UUID `json:"from_user_id"`
	ToUserId   uuid.UUID `json:"to_user_id"`
}
