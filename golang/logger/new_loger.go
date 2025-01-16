// File: logger/logger.go

package logger

import (
    "sync"
)

var (
    instance *NewLoggerType
    once     sync.Once
)

// Initialize creates a new logger instance if one doesn't exist
func Initialize(config *LoggerConfig) {
    once.Do(func() {
        instance = &NewLoggerType{
            paths:        config.paths,
            isDev:        true, // You can make this configurable
            mu:           sync.RWMutex{},
            timeFormat:   "2006-01-02 15:04:05",
            minimumLevel: NewInfoLevel,
        }
    })
}

// GetLogger returns the singleton logger instance
func GetLogger() *NewLoggerType {
    if instance == nil {
        // If not initialized, create with default config
        Initialize(NewLoggerConfig())
    }
    return instance
}

// helper methods for the logger
func (l *NewLoggerType) Info(msg string) {
    if l.minimumLevel <= NewInfoLevel {
        // Your logging implementation
    }
}

func (l *NewLoggerType) Error(msg string) {
    if l.minimumLevel <= NewErrorLevel {
        // Your logging implementation
    }
}

func (l *NewLoggerType) Debug(msg string) {
    if l.minimumLevel <= NewDebugLevel {
        // Your logging implementation
    }
}