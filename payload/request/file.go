package request

type CreateFileRequest struct {
	Name     string `json:"name"`
	Key      string `json:"key"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
}

type TransferRequest struct {
	ToEmail string `json:"to_email"`
	FileId  string `json:"file_id"`
}
