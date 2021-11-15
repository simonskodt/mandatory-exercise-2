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
	file *os.File
}

// InfoPrintln Prints both to the log file and to the console.
// Equals to calling logger.InfoLogger.Println() and a normal log.println()
func (l *Logger) InfoPrintln(v ...interface{}) {
	l.InfoLogger.Println(v...)
	log.Println(v...)
}

func (l *Logger) InfoPrintf(format string, v ...interface{}) {
	l.InfoLogger.Printf(format, v...)
	log.Printf(format, v...)
}

func (l *Logger) WarningPrintln(v ...interface{}) {
	l.WarningLogger.Println(v...)
	log.Println(v...)
}

func (l *Logger) WarningPrintf(format string, v ...interface{}) {
	l.WarningLogger.Printf(format, v...)
	log.Printf(format, v...)
}

func (l *Logger) ErrorPrintf(format string, v ...interface{}) {
	l.ErrorLogger.Printf(format, v...)
	log.Printf(format, v...)
}

func (l *Logger) ErrorFatalf(format string, v ...interface{}) {
	l.ErrorLogger.Printf(format, v...)
	log.Fatalf(format, v...)
}

// DeleteLog deletes the log file associated with this logger.
func (l *Logger) DeleteLog() {
	err := l.file.Close()
	if err != nil {
		l.ErrorLogger.Printf("Could not close log. Deletion terminated. :: %v", err)
		return
	}
	_ = os.Remove(l.file.Name())
}

// NewLogger creates a new Logger and binds it to a file with the given filename.
func NewLogger(filename string) *Logger {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		_ = os.Mkdir(logPath, os.ModeDir)
	}

	file, err := os.OpenFile(logPath+filename+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return &Logger{
		InfoLogger:    log.New(file, "INFO: ", log.LstdFlags),
		WarningLogger: log.New(file, "WARNING: ", log.LstdFlags),
		ErrorLogger:   log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile),
		file: file,
	}
}


