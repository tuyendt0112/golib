package log

import (
	"log/slog"
	"os"
	"testing"
)

func TestSetLogHandler(t *testing.T) {
	// Save original logger
	originalLogger := slog.Default()
	defer slog.SetDefault(originalLogger)

	// Test case 1: Set log handler in development mode
	os.Setenv("ENV", "dev")
	defer os.Unsetenv("ENV")

	SetLogHandler()
	
	// Verify logger is set (should not panic)
	slog.Info("Test log message")
	slog.Debug("Test debug message")
	
	// Test case 2: Set log handler in production mode
	os.Setenv("ENV", "prod")
	SetLogHandler()
	
	slog.Info("Test log message in prod")
	// Debug should not appear in prod mode
	slog.Debug("This debug should not appear")
	
	// Test case 3: Set log handler with production string
	os.Setenv("ENV", "production")
	SetLogHandler()
	
	slog.Info("Test log message in production")
	
	// Test case 4: Set log handler with master string
	os.Setenv("ENV", "master")
	SetLogHandler()
	
	slog.Info("Test log message in master")
}

