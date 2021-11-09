package utils

import (
	"log"
	"os"
)

const logPath = "..\\logs\\"

type Logger struct {
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
}

func NewLogger(filename string) *Logger {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		_ = os.Mkdir(logPath, os.ModeDir)
	}

	f, err := os.OpenFile(logPath+filename+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return &Logger{
		InfoLogger:    log.New(f, "INFO: ", log.LstdFlags),
		WarningLogger: log.New(f, "WARNING: ", log.LstdFlags),
		ErrorLogger:   log.New(f, "ERROR: ", log.LstdFlags|log.Lshortfile),
	}
}


