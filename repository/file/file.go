package file_repository

import (
	"strings"

	"github.com/google/uuid"

	"github.com/horlakz/api.secretariat_repository/lib/database"
	"github.com/horlakz/api.secretariat_repository/model"
	"github.com/horlakz/api.secretariat_repository/repository"
)

type FilePageable struct {
	repository.Pageable

	UserId uuid.UUID `json:"user_id"`
}

type FileRepositoryInterface interface {
	Create(file model.File) (model.File, error)
	FindAllFiles(pageable FilePageable) ([]model.File, repository.Pagination, error)
	FindFileById(uuid uuid.UUID) (model.File, error)
	FindFileByKey(key string) (model.File, error)
	UpdateFile(file model.File) (model.File, error)
	DeleteFile(id uuid.UUID) error
}

type fileRepository struct {
	database database.DatabaseInterface
}

func NewFileRepository(database database.DatabaseInterface) FileRepositoryInterface {
	return &fileRepository{database: database}
}

// Create implements FileRepositoryInterface.
func (f *fileRepository) Create(file model.File) (model.File, error) {
	file.Prepare()

	if err := f.database.Connection().Create(&file).Error; err != nil {
		return model.File{}, err
	}

	return file, nil
}

// FindAllFiles implements FileRepositoryInterface.
func (f *fileRepository) FindAllFiles(pageable FilePageable) (files []model.File, pagination repository.Pagination, err error) {
	var file model.File

	pagination.CurrentPage = int64(pageable.Page)
	pagination.TotalItems = 0
	pagination.TotalPages = 1

	search := strings.TrimSpace(pageable.Search)
	offset := (pageable.Page - 1) * pageable.Size
	model := f.database.Connection().Model(&file)

	if len(search) > 0 {
		model = model.Where("original_name LIKE ?", "%"+search+"%")
	}

	if pageable.UserId != uuid.Nil {
		model = model.Where("user_id = ?", pageable.UserId)
	}

	if err = model.Count(&pagination.TotalItems).Error; err != nil {
		return nil, pagination, err
	}

	// apply pagination
	paginatedQuery := model.
		Offset(offset).
		Limit(pageable.Size).
		Order(pageable.SortBy + " " + pageable.SortDirection)

	if err = paginatedQuery.Find(&files).Error; err != nil {
		return nil, pagination, err
	}

	if pagination.TotalItems > 0 {
		pagination.TotalPages = (pagination.TotalItems + int64(pageable.Size) - 1) / int64(pageable.Size)
	} else {
		pagination.TotalPages = 1
	}

	return files, pagination, err
}

// FindFileById implements FileRepositoryInterface.
func (f *fileRepository) FindFileById(uuid uuid.UUID) (model.File, error) {
	var file model.File

	if err := f.database.Connection().Where("id = ?", uuid).First(&file).Error; err != nil {
		return model.File{}, err
	}

	return file, nil
}

// FindFileByKey implementes FileRepositoryInterface
func (f *fileRepository) FindFileByKey(key string) (file model.File, err error) {
	if err = f.database.Connection().Where("key = ?", key).First(&file).Error; err != nil {
		return model.File{}, err
	}

	return file, nil
}

// UpdateFile implements FileRepositoryInterface.
func (f *fileRepository) UpdateFile(file model.File) (model.File, error) {
	// Retrieve the existing file from the database.
	var existingFile model.File
	if err := f.database.Connection().Where("id = ?", file.ID).First(&existingFile).Error; err != nil {
		return model.File{}, err
	}

	// Apply updates to the fields of the existing file record.
	// Only update the fields that are present in the `file` parameter.
	if file.Name != "" {
		existingFile.Name = file.Name
	}
	if file.Key != "" {
		existingFile.Key = file.Key
	}

	// Save the updated file record back to the database.
	if err := f.database.Connection().Save(&existingFile).Error; err != nil {
		return model.File{}, err
	}

	return existingFile, nil
}

// DeleteFile implements FileRepositoryInterface.
func (f *fileRepository) DeleteFile(id uuid.UUID) error {
	file, err := f.FindFileById(id)

	if err != nil {
		return err
	}

	if err := f.database.Connection().Delete(&file).Error; err != nil {
		return err
	}

	return nil
}
