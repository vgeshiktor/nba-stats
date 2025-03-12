// pkg/logger/logger.go
package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// LogEntry defines the structure of log messages in JSON format.
type LogEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	File      string `json:"file,omitempty"`
	Message   string `json:"message"`
}

// Logger instance
var (
	infoLogger  = log.New(os.Stdout, "", log.Lshortfile)
	errorLogger = log.New(os.Stderr, "", log.Lshortfile)
)

// logJSON creates a JSON log entry.
func logJSON(level, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	_, file, line, _ := runtime.Caller(2) // Get caller file and line number

	logEntry := LogEntry{
		Level:     level,
		Timestamp: time.Now().Format(time.RFC3339),
		File:      fmt.Sprintf("%s:%d", file, line),
		Message:   message,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"level": "ERROR", "message": "Failed to log message"}`)
		return
	}

	// Print JSON logs
	if level == "ERROR" {
		errorLogger.Println(string(jsonData))
	} else {
		infoLogger.Println(string(jsonData))
	}
}

// Info logs informational messages in JSON format.
func Info(format string, v ...interface{}) {
	logJSON("INFO", format, v...)
}

// Error logs error messages in JSON format.
func Error(format string, v ...interface{}) {
	logJSON("ERROR", format, v...)
}
