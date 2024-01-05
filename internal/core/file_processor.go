package core

import (
	"errors"
	"hexagonal-architecture/internal"
	"io"
	"path/filepath"
	"strings"
)

type FileProcessor struct {
	ConsensusPort internal.ConsensusPort
	DatabasePort  internal.DatabasePort
	LoggerPort    internal.LoggerPort
}

func NewFileProcessor(consensusPort internal.ConsensusPort, databasePort internal.DatabasePort, loggerPort internal.LoggerPort) internal.FileProcessorPort {
	return &FileProcessor{
		ConsensusPort: consensusPort,
		DatabasePort:  databasePort,
		LoggerPort:    loggerPort,
	}
}

func (fp *FileProcessor) ProcessFile(fileName string, file io.Reader) (int64, error) {
	//check if file is already present in database
	fileSize, err := fp.DatabasePort.Retrieve(fileName)
	if err == nil {
		fp.LoggerPort.Info("File already present in database")
		return fileSize, nil
	}
	fileInfo, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	if len(fileInfo) == 0 {
		return 0, errors.New("uploaded file is empty")
	}
	// Check if the file has a valid extension
	if !isValidFileExtension(fileName) {
		fp.LoggerPort.Error("Invalid file extension", "fileName", fileName)
		return 0, errors.New("invalid file extension")
	}
	fileSize = int64(len(fileInfo))

	// Update the database using consensus mechanism
	fp.ConsensusPort.UpdateDatabase(fileName, fileSize)

	err = fp.DatabasePort.Store(fileName, fileSize)
	if err != nil {
		fp.LoggerPort.Error("Error storing file in the database", "fileName", fileName, "error", err.Error())
		return 0, err
	}

	fp.LoggerPort.Info("File processed successfully")

	return fileSize, nil
}
func isValidFileExtension(filename string) bool {
	allowedExtensions := map[string]bool{
		".txt":  true,
		".pdf":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}
