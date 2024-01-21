package api

import (
	"fmt"
	"hexagonal-architecture/internal/consensus"
	"hexagonal-architecture/internal/core"
	"hexagonal-architecture/internal/database"
	"hexagonal-architecture/internal/logger"
	"net/http"

	"go.etcd.io/raft/v3"

	"github.com/gin-gonic/gin"
)

func FileUploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
		return
	}
	fileName := file.Filename
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer fileContent.Close()

	// Initialize raft, database, and logger adapters
	raftAdapter := consensus.NewRaftConsensusAdapter(1, []raft.Peer{{ID: uint64(2)}, {ID: uint64(3)}, {ID: uint64(4)}, {ID: uint64(5)}})
	redisAdapter := database.NewRedisAdapter()
	hclogAdapter := logger.NewHclogAdapter()

	// Create instances of core with injected dependencies
	fileProcessor := core.NewFileProcessor(raftAdapter, redisAdapter, hclogAdapter)

	// Validate file and calculate size
	fileSize, err := fileProcessor.ProcessFile(fileName, fileContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("File %s uploaded successfully with size:%d", fileName, fileSize)})

}
