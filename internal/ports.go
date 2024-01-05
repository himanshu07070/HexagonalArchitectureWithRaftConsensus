package internal

import "io"

type ConsensusPort interface {
	UpdateDatabase(fileName string, fileSize int64)
}
type DatabasePort interface {
	Store(fileName string, fileSize int64) error
	Retrieve(fileName string) (int64, error)
}

type LoggerPort interface {
	Info(message string, args ...interface{})
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Error(message string, args ...interface{})
}

type FileProcessorPort interface {
	ProcessFile(filename string, file io.Reader) (int64, error)
}
