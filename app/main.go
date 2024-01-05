package main

import (
	"hexagonal-architecture/internal/api"
	"hexagonal-architecture/internal/logger"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	logs := logger.NewHclogAdapter()
	// Register the file upload endpoint
	router.POST("/upload", api.FileUploadHandler)

	// Start the server
	serverAddr := ":8080"
	logs.Debug("Server listening on", serverAddr, "address")
	if err := router.Run(serverAddr); err != nil {
		logs.Error("Failed to start server:", err)
	}
}
