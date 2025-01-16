package logger

import (
	"fmt"
	"path/filepath"
	"strings"
)

// formatPath ensures consistent path formatting across different operating systems
func formatPath(path string) string {
	// Remove file extension if present
	path = strings.TrimSuffix(path, filepath.Ext(path))
	
	// Convert to platform-specific path format
	return filepath.FromSlash(path)
}

// NewLoggerConfig creates and initializes the logging configuration
func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		paths: GetLoggerPaths(),
	}
}

// UpdatePathStatus enables or disables a specific logging path
// Returns an error if the specified path is not found
func (c *LoggerConfig) UpdatePathStatus(pathName string, enabled bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, path := range c.paths {
		if path.Path == pathName {
			path.Enabled = enabled
			// Update enabled/disabled log IDs based on the new status
			if enabled {
				path.EnabledLogIDs = path.LogIDs
				path.DisabledLogIDs = []int{}
			} else {
				path.EnabledLogIDs = []int{}
				path.DisabledLogIDs = path.LogIDs
			}
			return nil
		}
	}
	return fmt.Errorf("path not found: %s", pathName)
}