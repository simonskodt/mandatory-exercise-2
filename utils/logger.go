package utils

import (
	"log"
	"os"
)

// logPath is the path to the output directory for the Logger.
const logPath = "..\\logs\\"

// Logger is a log used to write log calls to a file.
// It consists of three prefix loggers InfoLogger, WarningLogger and ErrorLogger,
// which sets their respective prefix to the log file.
type Logger struct {
	InfoLogger    *log.Logger	// InfoLogger adds the prefix "INFO" to the log string at each call.
	WarningLogger *log.Logger	// WarningLogger adds the prefix "WARNING" to the log string at each call.
	ErrorLogger   *log.Logger	// ErrorLogger adds the prefix "ERROR" to the log string at each call.
}

// InfoPrintln Prints both to the log file and to the console.
// Equals to calling Logger.InfoLogger.Println() and a normal log.println()
func (l *Logger) InfoPrintln(v ...interface{}) {
	l.InfoLogger.Println(v...)
	log.Println(v...)
}

func (l *Logger) WarningPrintln(v ...interface{}) {
	l.WarningLogger.Println(v...)
	log.Println(v...)
}

func (l *Logger) WarningPrintf(format string, v ...interface{}) {
	l.WarningLogger.Printf(format, v...)
	log.Printf(format, v...)
}

// NewLogger creates a new Logger and binds it to a file with the given filename.
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


