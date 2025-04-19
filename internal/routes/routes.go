package routes

import (
	"github.com/gin-gonic/gin"

	"model-registry/internal/handlers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", handlers.Root)

	modelGroup := router.Group("/model")
	{
		modelGroup.POST("/", handlers.AddModel)
		modelGroup.GET("/", handlers.ListModels)
	}

	router.POST("/save", handlers.SaveModelFile)
	router.GET("/download", handlers.DownloadModelFile)
	router.GET("/versions", handlers.ListVersions)

	return router
}
