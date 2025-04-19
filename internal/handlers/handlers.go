package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"model-registry/internal/models"
	"model-registry/internal/repositories"
	"model-registry/internal/utils"
)

var modelRepo repositories.ModelRepository
var versionRepo repositories.VersionRepository

func InitHandlers(mRepo repositories.ModelRepository, vRepo repositories.VersionRepository) {
	modelRepo = mRepo
	versionRepo = vRepo
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and type are required"})
		return
	}

	model, err := modelRepo.CreateModel(name, description, modelType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"model_id": model.ID})
}

func ListModels(c *gin.Context) {
	models, err := modelRepo.ListModels()
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model_id"})
		return
	}

	versionName := c.Query("version_name")
	if versionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version_name is required"})
		return
	}

	if versionName == "latest" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version_name cannot be 'latest'"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if filepath.Ext(file.Filename) != ".model" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file extension"})
		return
	}

	model, err := modelRepo.GetModel(modelID)
	if err != nil || model == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	existingVersion, err := versionRepo.GetVersion(modelID, versionName)
	if err == nil && existingVersion != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Version already exists"})
		return
	}

	filePath := utils.GenerateFilePath(modelID, versionName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	version, err := versionRepo.CreateVersion(modelID, versionName, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create version"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"version_id": version.ID})
}

func DownloadModelFile(c *gin.Context) {
	modelID, err := strconv.Atoi(c.Query("model_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model_id"})
		return
	}

	versionName := c.Query("version_name")
	if versionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "version_name is required"})
		return
	}

	var version *models.ModelVersion
	if versionName == "latest" {
		version, err = versionRepo.GetLatestVersion(modelID)
		if err != nil || version == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Latest version not found"})
			return
		}
	} else {
		version, err = versionRepo.GetVersion(modelID, versionName)
		if err != nil || version == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
			return
		}
	}

	fileContent, err := os.ReadFile(version.FilePath)
	if err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model_id"})
			return
		}
		modelID = &id
	}

	versions, err := versionRepo.ListVersions(modelID)
	if err != nil {
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
