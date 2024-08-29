package file_service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/horlakz/api.secretariat_repository/dto"
	"github.com/horlakz/api.secretariat_repository/model"
	"github.com/horlakz/api.secretariat_repository/repository"
	file_repository "github.com/horlakz/api.secretariat_repository/repository/file"
	user_repository "github.com/horlakz/api.secretariat_repository/repository/user"
	"gorm.io/gorm"
)

type fileService struct {
	fileRepository    file_repository.FileRepositoryInterface
	transferRepostory file_repository.TransferRepositoryInterface
	userRepository    user_repository.UserRepositoryInterface
}

type FileServiceInterface interface {
	CreateFile(dto dto.FileDTO) (dto.FileDTO, error)
	FindAllFiles(pageable file_repository.FilePageable) ([]dto.FileDTO, repository.Pagination, error)
	FindFileById(fileId uuid.UUID) (dto.FileDTO, error)
	UpdateFile(dto dto.FileDTO) (dto.FileDTO, error)
	DeleteFile(fileId uuid.UUID) error
	TransferFile(dto dto.TransferDTO, toEmail string) (dto.TransferDTO, error)
}

func NewFileService(
	fileRepository file_repository.FileRepositoryInterface,
	transferRepostory file_repository.TransferRepositoryInterface,
	userRepository user_repository.UserRepositoryInterface,
) FileServiceInterface {
	return &fileService{
		fileRepository:    fileRepository,
		transferRepostory: transferRepostory,
		userRepository:    userRepository,
	}
}

func (fs *fileService) ConvertToDTO(file model.File) (dto dto.FileDTO) {

	dto.ID = file.ID
	dto.Name = file.Name
	dto.Key = file.Key
	dto.MimeType = file.MimeType
	dto.Ext = file.Ext
	dto.Size = file.Size
	dto.UserId = file.UserId
	dto.CreatedAt = file.CreatedAt

	return dto
}

func (fs *fileService) ConvertToModel(dto dto.FileDTO) (file model.File) {

	file.ID = dto.ID
	file.Name = dto.Name
	file.Key = dto.Key
	file.MimeType = dto.MimeType
	file.Ext = dto.Ext
	file.Size = dto.Size
	file.UserId = dto.UserId

	return file
}

func (fs *fileService) CreateFile(dto dto.FileDTO) (dto.FileDTO, error) {
	file := fs.ConvertToModel(dto)

	createdFile, err := fs.fileRepository.Create(file)
	if err != nil {
		return dto, err
	}

	return fs.ConvertToDTO(createdFile), nil
}

func (fs *fileService) FindAllFiles(pageable file_repository.FilePageable) ([]dto.FileDTO, repository.Pagination, error) {
	files, pagination, err := fs.fileRepository.FindAllFiles(pageable)
	if err != nil {
		return nil, pagination, err
	}

	dtos := []dto.FileDTO{}
	for _, file := range files {
		dtos = append(dtos, fs.ConvertToDTO(file))
	}

	return dtos, pagination, nil
}

func (fs *fileService) FindFileById(fileId uuid.UUID) (dto.FileDTO, error) {
	file, err := fs.fileRepository.FindFileById(fileId)

	if err != nil {
		return dto.FileDTO{}, err
	}

	return fs.ConvertToDTO(file), nil
}

func (fs *fileService) UpdateFile(dto dto.FileDTO) (dto.FileDTO, error) {
	file := fs.ConvertToModel(dto)

	updatedFile, err := fs.fileRepository.UpdateFile(file)
	if err != nil {
		return dto, err
	}

	return fs.ConvertToDTO(updatedFile), nil
}

func (fs *fileService) TransferFile(dto dto.TransferDTO, toEmail string) (dto.TransferDTO, error) {
	toUser, err := fs.userRepository.FindUserByEmail(toEmail)

	if err == gorm.ErrRecordNotFound {
		return dto, errors.New("user not found")
	}

	if err != nil {
		return dto, err
	}

	transfer := model.Transfer{
		FileId:     dto.FileId,
		FromUserId: dto.FromUserId,
		ToUserId:   toUser.ID,
	}

	if _, err := fs.transferRepostory.Create(transfer); err != nil {
		return dto, err
	}

	file, err := fs.fileRepository.FindFileById(dto.FileId)
	if err != nil {
		return dto, err
	}

	file.UserId = toUser.ID

	if _, err := fs.fileRepository.Create(file); err != nil {
		return dto, err
	}

	return dto, nil
}

func (fs *fileService) DeleteFile(fileId uuid.UUID) error {
	return fs.fileRepository.DeleteFile(fileId)
}
