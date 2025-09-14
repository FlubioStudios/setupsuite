package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Logger provides verbose logging functionality
type Logger struct {
	enabled bool
	logFile *os.File
	logger  *log.Logger
}

// VerboseLogger is the global logger instance
var VerboseLogger *Logger

// InitLogger initializes the verbose logger
func InitLogger(enabled bool) error {
	VerboseLogger = &Logger{enabled: enabled}

	if !enabled {
		return nil
	}

	// Create logs directory if it doesn't exist
	logDir := "/var/log/setupsuite"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// Fallback to local directory if /var/log is not writable
		logDir = "./logs"
		os.MkdirAll(logDir, 0755)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(logDir, fmt.Sprintf("setupsuite_%s.log", timestamp))

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	VerboseLogger.logFile = file
	VerboseLogger.logger = log.New(file, "", log.LstdFlags|log.Lmicroseconds)

	VerboseLogger.LogInfo("=== SetupSuite Verbose Logging Started ===")
	VerboseLogger.LogInfo("Log file: %s", logPath)

	fmt.Printf("Verbose logging enabled: %s\n", logPath)
	return nil
}

// CloseLogger closes the log file
func CloseLogger() {
	if VerboseLogger != nil && VerboseLogger.logFile != nil {
		VerboseLogger.LogInfo("=== SetupSuite Verbose Logging Ended ===")
		VerboseLogger.logFile.Close()
	}
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(format string, args ...interface{}) {
	if !l.enabled || l.logger == nil {
		return
	}
	l.logger.Printf("[INFO] "+format, args...)
}

// LogFileOperation logs file operations (create, write, read, etc.)
func (l *Logger) LogFileOperation(operation, path string) {
	if !l.enabled {
		return
	}
	l.LogInfo("FILE_OP: %s -> %s", operation, path)
}

// LogFileContent logs file content being written
func (l *Logger) LogFileContent(path, content string) {
	if !l.enabled {
		return
	}
	l.LogInfo("FILE_CONTENT: %s\n--- CONTENT START ---\n%s\n--- CONTENT END ---", path, content)
}

// LogCommand logs command execution
func (l *Logger) LogCommand(cmd string, args []string) {
	if !l.enabled {
		return
	}
	l.LogInfo("COMMAND: %s %s", cmd, strings.Join(args, " "))
}

// LogCommandOutput logs command output
func (l *Logger) LogCommandOutput(cmd string, args []string, output []byte, err error) {
	if !l.enabled {
		return
	}
	if err != nil {
		l.LogInfo("COMMAND_ERROR: %s %s -> ERROR: %v", cmd, strings.Join(args, " "), err)
	}
	if len(output) > 0 {
		l.LogInfo("COMMAND_OUTPUT: %s %s\n--- OUTPUT START ---\n%s\n--- OUTPUT END ---",
			cmd, strings.Join(args, " "), string(output))
	}
}

// LogError logs an error
func (l *Logger) LogError(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	l.LogInfo("[ERROR] "+format, args...)
}

// LogWarning logs a warning
func (l *Logger) LogWarning(format string, args ...interface{}) {
	if !l.enabled {
		return
	}
	l.LogInfo("[WARNING] "+format, args...)
}

// VerboseOpenFile wraps os.OpenFile with logging
func VerboseOpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	operation := "OPEN"
	if flag&os.O_CREATE != 0 {
		operation = "CREATE"
	}
	if flag&os.O_WRONLY != 0 || flag&os.O_RDWR != 0 {
		operation += "_WRITE"
	}

	VerboseLogger.LogFileOperation(operation, name)
	return os.OpenFile(name, flag, perm)
}

// VerboseCreate wraps os.Create with logging
func VerboseCreate(name string) (*os.File, error) {
	VerboseLogger.LogFileOperation("CREATE", name)
	return os.Create(name)
}

// VerboseMkdirAll wraps os.MkdirAll with logging
func VerboseMkdirAll(path string, perm os.FileMode) error {
	VerboseLogger.LogFileOperation("MKDIR_ALL", path)
	return os.MkdirAll(path, perm)
}

// VerboseWriteFile writes content to file with logging
func VerboseWriteFile(filename, content string) error {
	VerboseLogger.LogFileOperation("WRITE_FILE", filename)
	VerboseLogger.LogFileContent(filename, content)

	file, err := VerboseCreate(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// VerboseWrite writes data to a file and logs the operation
func VerboseWrite(file *os.File, data []byte) (int, error) {
	n, err := file.Write(data)
	if VerboseLogger.enabled {
		VerboseLogger.LogInfo("Wrote %d bytes to file: %s", n, file.Name())
		if err != nil {
			VerboseLogger.LogError("Write error: %s", err)
		}
	}
	return n, err
}

// VerboseWriteString writes a string to a file and logs the operation
func VerboseWriteString(file *os.File, data string) (int, error) {
	n, err := file.WriteString(data)
	if VerboseLogger.enabled {
		VerboseLogger.LogInfo("Wrote string (%d bytes) to file: %s", n, file.Name())
		if err != nil {
			VerboseLogger.LogError("WriteString error: %s", err)
		}
	}
	return n, err
}

// VerboseCommand wraps exec.Command with logging
func VerboseCommand(name string, args ...string) *exec.Cmd {
	VerboseLogger.LogCommand(name, args)
	return exec.Command(name, args...)
}

// VerboseCommandRun runs a command and logs its output
func VerboseCommandRun(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	VerboseLogger.LogCommand(name, args)

	output, err := cmd.CombinedOutput()
	VerboseLogger.LogCommandOutput(name, args, output, err)

	return err
}

// VerboseCommandOutput runs a command and returns output with logging
func VerboseCommandOutput(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	VerboseLogger.LogCommand(name, args)

	output, err := cmd.Output()
	VerboseLogger.LogCommandOutput(name, args, output, err)

	return output, err
}
