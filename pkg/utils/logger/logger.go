package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New creates a new logger instance with default configuration
func New() *logrus.Logger {
	logger := logrus.New()
	
	// Set default configuration
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	
	return logger
}

// NewWithConfig creates a new logger with custom configuration
func NewWithConfig(level string, format string, output string) *logrus.Logger {
	logger := logrus.New()
	
	// Set log level
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	
	// Set formatter
	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}
	
	// Set output
	switch output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		if output != "" {
			// Try to open file
			if file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
				logger.SetOutput(file)
			} else {
				logger.SetOutput(os.Stdout)
				logger.WithError(err).Warn("Failed to open log file, using stdout")
			}
		} else {
			logger.SetOutput(os.Stdout)
		}
	}
	
	return logger
}