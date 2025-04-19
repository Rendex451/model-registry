package repositories

import (
	"gorm.io/gorm"

	"model-registry/internal/models"
)

type VersionRepository interface {
	CreateVersion(modelID int, versionName, filePath string) (*models.ModelVersion, error)
	GetVersion(modelID int, versionName string) (*models.ModelVersion, error)
	GetLatestVersion(modelID int) (*models.ModelVersion, error)
	ListVersions(modelID *int) ([]models.ModelVersion, error)
}

type versionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) VersionRepository {
	return &versionRepository{db: db}
}

func (r *versionRepository) CreateVersion(modelID int, versionName, filePath string) (*models.ModelVersion, error) {
	version := &models.ModelVersion{
		ModelID:     modelID,
		VersionName: versionName,
		FilePath:    filePath,
	}
	if err := r.db.Create(version).Error; err != nil {
		return nil, err
	}

	return version, nil
}

func (r *versionRepository) GetVersion(modelID int, versionName string) (*models.ModelVersion, error) {
	var version models.ModelVersion
	if result := r.db.Where("model_id = ? AND version_name = ?", modelID, versionName).First(&version); result.Error != nil {
		return nil, result.Error
	}

	return &version, nil
}

func (r *versionRepository) GetLatestVersion(modelID int) (*models.ModelVersion, error) {
	var version models.ModelVersion
	if result := r.db.Where("model_id = ?", modelID).Order("date_added DESC").First(&version); result.Error != nil {
		return nil, result.Error
	}

	return &version, nil
}

func (r *versionRepository) ListVersions(modelID *int) ([]models.ModelVersion, error) {
	var versions []models.ModelVersion
	query := r.db
	if modelID != nil {
		query = query.Where("model_id = ?", *modelID)
	}
	if err := query.Find(&versions).Error; err != nil {
		return nil, err
	}

	return versions, nil
}
