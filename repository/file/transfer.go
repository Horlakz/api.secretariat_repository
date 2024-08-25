package file_repository

import (
	"github.com/google/uuid"

	"github.com/horlakz/api.secretariat_repository/lib/database"
	"github.com/horlakz/api.secretariat_repository/model"
)

type TransferRepositoryInterface interface {
	Create(transfer model.Transfer) (model.Transfer, error)
	FindById(id uuid.UUID) (model.Transfer, error)
	FindAllByFromUserId(userId uuid.UUID) ([]model.Transfer, error)
	Delete(id uuid.UUID) error
}

type transferRepository struct {
	database database.DatabaseInterface
}

func NewTransferRepository(database database.DatabaseInterface) TransferRepositoryInterface {
	return &transferRepository{database: database}
}

// Create implements TransferRepositoryInterface.
func (t *transferRepository) Create(transfer model.Transfer) (model.Transfer, error) {
	// Prepare the transfer model if needed
	transfer.Prepare()

	// e.g., transfer.Prepare() (if the model has a Prepare method)
	if err := t.database.Connection().Create(&transfer).Error; err != nil {
		return model.Transfer{}, err
	}
	return transfer, nil
}

// FindById implements TransferRepositoryInterface.
func (t *transferRepository) FindById(id uuid.UUID) (model.Transfer, error) {
	var transfer model.Transfer
	if err := t.database.Connection().Where("id = ?", id).First(&transfer).Error; err != nil {
		return model.Transfer{}, err
	}
	return transfer, nil
}

// FindAllByFromUserId implements TransferRepositoryInterface.
func (t *transferRepository) FindAllByFromUserId(userId uuid.UUID) ([]model.Transfer, error) {
	var transfers []model.Transfer
	if err := t.database.Connection().Where("user_id = ?", userId).Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// Delete implements TransferRepositoryInterface.
func (t *transferRepository) Delete(id uuid.UUID) error {
	transfer, err := t.FindById(id)
	if err != nil {
		return err
	}

	if err := t.database.Connection().Delete(&transfer).Error; err != nil {
		return err
	}

	return nil
}
