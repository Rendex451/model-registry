package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"model-registry/internal/models"
	"model-registry/internal/repositories"
	"model-registry/internal/utils"
)

var (
	modelRepo   repositories.ModelRepository
	versionRepo repositories.VersionRepository
	logger      *slog.Logger
)

func InitHandlers(mRepo repositories.ModelRepository, vRepo repositories.VersionRepository, log *slog.Logger) {
	modelRepo = mRepo
	versionRepo = vRepo
	logger = log
}

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello from model-registry service",
		"version": "v1",
	})
}

func AddModel(c *gin.Context) {
	name := c.Query("name")
	description := c.Query("description")
	modelType := c.Query("type")

	if name == "" || modelType == "" {
		logger.Warn("AddModel failed: missing required parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and type are required"})
		return
	}

	model, err := modelRepo.CreateModel(name, description, modelType)
	if err != nil {
		logger.Error("AddModel failed: unable to create model", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"model_id": model.ID})
}

func ListModels(c *gin.Context) {
	models, err := modelRepo.ListModels()
	if err != nil {
		logger.Error("ListModels failed: unable to fetch models", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch models"})
		return
	}

	response := []gin.H{}
	for _, m := range models {
		response = append(response, gin.H{
			"model_id":    m.ID,
			"name":        m.Name,
			"description": m.Description,
			"type":        m.Type,
		})
	}

	c.JSON(http.StatusOK, response)
}

func SaveModelFile(c *gin.Context) {
	modelID, err := strconv.Atoi(c.Query("model_id"))
	if err != nil {
		logger.Warn("SaveModelFile failed: invalid model_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model_id"})
		return
	}

	versionName := c.Query("version_name")
	if versionName == "" {
		logger.Warn("SaveModelFile failed: version_name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "version_name is required"})
		return
	}

	if versionName == "latest" {
		logger.Warn("SaveModelFile failed: version_name cannot be 'latest'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "version_name cannot be 'latest'"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Warn("SaveModelFile failed: file is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if filepath.Ext(file.Filename) != ".model" {
		logger.Warn("SaveModelFile failed: invalid file extension")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	model, err := modelRepo.GetModel(modelID)
	if err != nil || model == nil {
		logger.Warn("SaveModelFile failed: model not found", slog.Int("model_id", modelID))
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	existingVersion, err := versionRepo.GetVersion(modelID, versionName)
	if err == nil && existingVersion != nil {
		logger.Warn("SaveModelFile failed: version already exists", slog.Int("model_id", modelID), slog.String("version_name", versionName))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Version already exists"})
		return
	}

	filePath := utils.GenerateFilePath(modelID, versionName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		logger.Error("SaveModelFile failed: unable to save file", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	version, err := versionRepo.CreateVersion(modelID, versionName, filePath)
	if err != nil {
		logger.Error("SaveModelFile failed: unable to create version", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create version"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"version_id": version.ID})
}

func DownloadModelFile(c *gin.Context) {
	modelID, err := strconv.Atoi(c.Query("model_id"))
	if err != nil {
		logger.Warn("DownloadModelFile failed: invalid model_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model_id"})
		return
	}

	versionName := c.Query("version_name")
	if versionName == "" {
		logger.Warn("DownloadModelFile failed: version_name is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "version_name is required"})
		return
	}

	var version *models.ModelVersion
	if versionName == "latest" {
		version, err = versionRepo.GetLatestVersion(modelID)
		if err != nil || version == nil {
			logger.Warn("DownloadModelFile failed: latest version not found", slog.Int("model_id", modelID))
			c.JSON(http.StatusNotFound, gin.H{"error": "Latest version not found"})
			return
		}
	} else {
		version, err = versionRepo.GetVersion(modelID, versionName)
		if err != nil || version == nil {
			logger.Warn("DownloadModelFile failed: version not found", slog.Int("model_id", modelID), slog.String("version_name", versionName))
			c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
			return
		}
	}

	fileContent, err := os.ReadFile(version.FilePath)
	if err != nil {
		logger.Error("DownloadModelFile failed: unable to read file", slog.String("file_path", version.FilePath), slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", fileContent)
}

func ListVersions(c *gin.Context) {
	modelIDParam := c.Query("model_id")
	var modelID *int
	if modelIDParam != "" {
		id, err := strconv.Atoi(modelIDParam)
		if err != nil {
			logger.Warn("ListVersions failed: invalid model_id", slog.String("model_id_param", modelIDParam))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model_id"})
			return
		}
		modelID = &id
	}

	versions, err := versionRepo.ListVersions(modelID)
	if err != nil {
		logger.Error("ListVersions failed: unable to fetch versions", slog.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions"})
		return
	}

	response := []gin.H{}
	for _, v := range versions {
		response = append(response, gin.H{
			"version_id":   v.ID,
			"model_id":     v.ModelID,
			"version_name": v.VersionName,
			"date_added":   v.DateAdded,
		})
	}

	c.JSON(http.StatusOK, response)
}
