// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul-template/version"
)

var (
	// DefaultLogFileName is the default filename if the user didn't specify one
	// which means that the user specified a directory to log to
	DefaultLogFileName = fmt.Sprintf("%s.log", version.Name)

	// DefaultLogRotateDuration is the default time taken by the agent to rotate logs
	DefaultLogRotateDuration = 24 * time.Hour
)

type LogFileConfig struct {
	// LogFilePath is the path to the file the logs get written to
	LogFilePath *string `mapstructure:"path"`

	// LogRotateBytes is the maximum number of bytes that should be written to a log
	// file
	LogRotateBytes *int `mapstructure:"log_rotate_bytes"`

	// LogRotateDuration is the time after which log rotation needs to be performed
	LogRotateDuration *time.Duration `mapstructure:"log_rotate_duration"`

	// LogRotateMaxFiles is the maximum number of log file archives to keep
	LogRotateMaxFiles *int `mapstructure:"log_rotate_max_files"`
}

// DefaultLogFileConfig returns a configuration that is populated with the
// default values.
func DefaultLogFileConfig() *LogFileConfig {
	return &LogFileConfig{}
}

// Copy returns a deep copy of this configuration.
func (c *LogFileConfig) Copy() *LogFileConfig {
	if c == nil {
		return nil
	}

	var o LogFileConfig
	o.LogFilePath = c.LogFilePath
	o.LogRotateBytes = c.LogRotateBytes
	o.LogRotateDuration = c.LogRotateDuration
	o.LogRotateMaxFiles = c.LogRotateMaxFiles
	return &o
}

// Merge combines all values in this configuration with the values in the other
// configuration, with values in the other configuration taking precedence.
// Maps and slices are merged, most other values are overwritten. Complex
// structs define their own merge functionality.
func (c *LogFileConfig) Merge(o *LogFileConfig) *LogFileConfig {
	if c == nil {
		if o == nil {
			return nil
		}
		return o.Copy()
	}

	if o == nil {
		return c.Copy()
	}

	r := c.Copy()

	if o.LogFilePath != nil {
		r.LogFilePath = o.LogFilePath
	}

	if o.LogRotateBytes != nil {
		r.LogRotateBytes = o.LogRotateBytes
	}

	if o.LogRotateDuration != nil {
		r.LogRotateDuration = o.LogRotateDuration
	}

	if o.LogRotateMaxFiles != nil {
		r.LogRotateMaxFiles = o.LogRotateMaxFiles
	}

	return r
}

// Finalize ensures there no nil pointers.
func (c *LogFileConfig) Finalize() {
	if c.LogFilePath == nil {
		c.LogFilePath = String("")
	}

	if c.LogRotateBytes == nil {
		c.LogRotateBytes = Int(0)
	}

	if c.LogRotateDuration == nil {
		c.LogRotateDuration = TimeDuration(DefaultLogRotateDuration)
	}

	if c.LogRotateMaxFiles == nil {
		c.LogRotateMaxFiles = Int(0)
	}
}

// GoString defines the printable version of this struct.
func (c *LogFileConfig) GoString() string {
	if c == nil {
		return "(*LogFileConfig)(nil)"
	}

	return fmt.Sprintf("&LogFileConfig{"+
		"LogFilePath:%s, "+
		"LogRotateBytes:%s, "+
		"LogRotateDuration:%s, "+
		"LogRotateMaxFiles:%s, "+
		"}",
		StringGoString(c.LogFilePath),
		IntGoString(c.LogRotateBytes),
		TimeDurationGoString(c.LogRotateDuration),
		IntGoString(c.LogRotateMaxFiles),
	)
}
