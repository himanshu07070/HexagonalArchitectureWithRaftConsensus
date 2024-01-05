package core

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockConsensusPort struct {
	UpdateDatabaseFunc func(fileName string, fileSize int64)
}

func (mcp *MockConsensusPort) UpdateDatabase(fileName string, fileSize int64) {
	if mcp.UpdateDatabaseFunc != nil {
		mcp.UpdateDatabaseFunc(fileName, fileSize)
	}
}

type MockDatabasePort struct {
	StoredFileName   string
	StoredFileSize   int64
	RetrievedFile    string
	RetrievedFileErr error
}

func (m *MockDatabasePort) Store(fileName string, fileSize int64) error {
	m.StoredFileName = fileName
	m.StoredFileSize = fileSize
	return nil
}

func (m *MockDatabasePort) Retrieve(fileName string) (int64, error) {
	if m.RetrievedFileErr != nil {
		return 0, m.RetrievedFileErr
	}
	if fileName == m.RetrievedFile {
		return 42, nil
	}
	return 0, errors.New("file not found")
}

type MockLoggerPort struct {
	InfoLogged  string
	ErrorLogged string
	TraceLogged string
	DebugLogged string
}

func (m *MockLoggerPort) Info(message string, args ...interface{}) {
	m.InfoLogged = message
}

func (m *MockLoggerPort) Error(message string, args ...interface{}) {
	m.ErrorLogged = message
}

func (m *MockLoggerPort) Trace(message string, args ...interface{}) {
	m.TraceLogged = message
}

func (m *MockLoggerPort) Debug(message string, args ...interface{}) {
	m.DebugLogged = message
}

func TestFileProcessor_ProcessFile(t *testing.T) {

	t.Run("FileAlreadyPresent", func(t *testing.T) {
		// Setup: Mock dependencies
		mockConsensus := &MockConsensusPort{}
		mockDatabase := &MockDatabasePort{}
		mockLogger := &MockLoggerPort{}
		fileProcessor := NewFileProcessor(mockConsensus, mockDatabase, mockLogger)
		// Test Case: File already present in the database
		mockDatabase.RetrievedFile = "existing_file.txt"
		size, err := fileProcessor.ProcessFile("existing_file.txt", nil)

		assert.NoError(t, err)
		assert.Equal(t, int64(42), size)
		assert.Equal(t, "File already present in database", mockLogger.InfoLogged)
	})

	t.Run("EmptyFile", func(t *testing.T) {
		// Setup: Mock dependencies
		mockConsensus := &MockConsensusPort{}
		mockDatabase := &MockDatabasePort{}
		mockLogger := &MockLoggerPort{}
		fileProcessor := NewFileProcessor(mockConsensus, mockDatabase, mockLogger)
		// Test Case: Empty file
		_, err := fileProcessor.ProcessFile("empty_file.txt", bytes.NewReader([]byte{}))

		assert.Error(t, err)
		assert.Equal(t, "uploaded file is empty", err.Error())
		assert.Empty(t, mockLogger.InfoLogged)
	})

	t.Run("InvalidFileExtension", func(t *testing.T) {
		// Setup: Mock dependencies
		mockConsensus := &MockConsensusPort{}
		mockDatabase := &MockDatabasePort{}
		mockLogger := &MockLoggerPort{}
		fileProcessor := NewFileProcessor(mockConsensus, mockDatabase, mockLogger)
		// Test Case: Invalid file extension
		_, err := fileProcessor.ProcessFile("invalid_file.exe", bytes.NewReader([]byte{1, 2, 3}))

		assert.Error(t, err)
		assert.Equal(t, "invalid file extension", err.Error())
		assert.Empty(t, mockLogger.InfoLogged)
	})

	t.Run("FileProcessingSuccessful", func(t *testing.T) {
		// Setup: Mock dependencies
		mockConsensus := &MockConsensusPort{}
		mockDatabase := &MockDatabasePort{}
		mockLogger := &MockLoggerPort{}
		fileProcessor := NewFileProcessor(mockConsensus, mockDatabase, mockLogger)
		// Test Case: File processing successful
		size, err := fileProcessor.ProcessFile("valid_file.txt", bytes.NewReader([]byte{1, 2, 3}))

		assert.NoError(t, err)
		assert.Equal(t, int64(3), size)
		assert.Equal(t, "File processed successfully", mockLogger.InfoLogged)
		assert.Equal(t, "valid_file.txt", mockDatabase.StoredFileName)
		assert.Equal(t, int64(3), mockDatabase.StoredFileSize)
	})

}
