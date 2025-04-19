package utils

import (
	"path/filepath"
	"strconv"

	"model-registry/internal/config"
)

func GenerateFilePath(modelID int, versionName string) string {
	baseDir, _ := filepath.Abs(config.GetConfig().ModelDir)
	return filepath.Join(baseDir, strconv.Itoa(modelID), versionName+".model")
}
