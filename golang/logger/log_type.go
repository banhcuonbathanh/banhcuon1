package logger


import (
	"io"
	"sync"
)

// LogLevel represents different severity levels for logging
// Higher values indicate more severe levels
type LogLevel int

const (
	NewDebugLevel LogLevel = iota // Most verbose logging level
	NewInfoLevel                 // General operational information
	NewWarnLevel                  // Warning messages for potentially harmful situations
	NewErrorLevel                 // Error messages for serious problems
)

// LogDescription defines metadata about a specific log entry
// Used to provide context and tracking for log messages
type LogDescription struct {
	Description string // Detailed description of what is being logged
	Location    string // Code location where the log is generated
	Status      string // Current status of this log entry (enabled/disabled)
}

// LogPath represents a configuration for a specific logging path
// Including its enabled state and associated log IDs
type LogPath struct {
	Path            string                  // File system path for the logs
	Enabled         bool                    // Whether logging is enabled for this path
	Description     string                  // Human-readable description of this logging path
	EnabledLogIDs   []int                  // List of currently enabled log IDs
	DisabledLogIDs  []int                  // List of currently disabled log IDs
	LogIDs          []int                  // All available log IDs for this path
	LogDescriptions map[int]LogDescription // Mapping of log IDs to their descriptions
}

// Logger represents the main logging system configuration and state
type NewLoggerType struct {
	paths        []*LogPath          // All configured logging paths
	isDev        bool               // Development mode flag
	mu           sync.RWMutex       // Mutex for thread-safe operations
	outputs      map[LogLevel]io.Writer // Output destinations per log level
	timeFormat   string             // Format string for timestamps
	minimumLevel LogLevel           // Minimum level of logs to process
}

// LoggerConfig manages the configuration state of the logging system
type LoggerConfig struct {
	paths []*LogPath   // Configured logging paths
	mu    sync.RWMutex // Mutex for thread-safe operations
}