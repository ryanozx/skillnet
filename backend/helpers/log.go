package helpers

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LoggerStruct struct {
	LogDirectory string
}

func InitialiseLoggerFile(LogsSavePath string) (*LoggerStruct, error) {
	pathErr := os.Mkdir(LogsSavePath, 0666)
	if pathErr != nil {
		return nil, pathErr
	}
	return &LoggerStruct{
		LogDirectory: LogsSavePath,
	}, nil
}

func (logger *LoggerStruct) setLogFile() *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", day, month.String(), year)
	filePath, _ := os.OpenFile(logger.LogDirectory+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	return filePath
}

func createLogger(filePath *os.File, prefix string) *log.Logger {
	return log.New(filePath, prefix, log.Ldate|log.Ltime|log.Lshortfile)
}

func (logger *LoggerStruct) Info() *log.Logger {
	getFilePath := logger.setLogFile()
	return createLogger(getFilePath, "INFO: ")
}

func (logger *LoggerStruct) Warning() *log.Logger {
	getFilePath := logger.setLogFile()
	return createLogger(getFilePath, "WARNING: ")
}

func (logger *LoggerStruct) Error() *log.Logger {
	getFilePath := logger.setLogFile()
	return createLogger(getFilePath, "ERROR: ")
}

func (logger *LoggerStruct) Fatal() *log.Logger {
	getFilePath := logger.setLogFile()
	return createLogger(getFilePath, "FATAL: ")
}
