package repositories

import (
	"model-registry/internal/models"

	"gorm.io/gorm"
)

type ModelRepository interface {
	CreateModel(name, description, modelType string) (*models.Model, error)
	GetModel(id int) (*models.Model, error)
	ListModels() ([]models.Model, error)
}

type modelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) ModelRepository {
	return &modelRepository{db: db}
}

func (r *modelRepository) CreateModel(name, description, modelType string) (*models.Model, error) {
	model := &models.Model{
		Name:        name,
		Description: description,
		Type:        modelType,
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

func (r *modelRepository) GetModel(modelID int) (*models.Model, error) {
	var model models.Model
	if err := r.db.First(&model, modelID).Error; err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *modelRepository) ListModels() ([]models.Model, error) {
	var models []models.Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return models, nil
}
