// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_SetupBasic(t *testing.T) {
	cfg := &LogConfig{Name: "test-system", LogLevel: log.Info}

	logger, err := Setup(cfg, nil)
	require.NoError(t, err)
	require.NotNil(t, logger)
	require.Equal(t, logger.Name(), "test-system")
	require.True(t, logger.IsInfo())
}

func TestLogger_SetupInvalidLogLevel(t *testing.T) {
	cfg := &LogConfig{}

	_, err := Setup(cfg, nil)
	assert.Containsf(t, err.Error(), "invalid log level", "expected error %s", err)
}

func TestLogger_SetupLoggerErrorLevel(t *testing.T) {
	cfg := &LogConfig{
		LogLevel: log.Error,
	}

	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Error("test error msg")
	logger.Info("test info msg")

	output := buf.String()

	require.Contains(t, output, "[ERROR] test error msg")
	require.NotContains(t, output, "[INFO]  test info msg")
}

func TestLogger_SetupLoggerDebugLevel(t *testing.T) {
	cfg := LogConfig{LogLevel: log.Debug}
	var buf bytes.Buffer

	logger, err := Setup(&cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("test info msg")
	logger.Debug("test debug msg")

	output := buf.String()

	require.Contains(t, output, "[INFO]  test info msg")
	require.Contains(t, output, "[DEBUG] test debug msg")
}

func TestLogger_SetupLoggerWithName(t *testing.T) {
	cfg := &LogConfig{
		LogLevel: log.Debug,
		Name:     "test-system",
	}
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Warn("test warn msg")

	require.Contains(t, buf.String(), "[WARN]  test-system: test warn msg")
}

func TestLogger_SetupLoggerWithJSON(t *testing.T) {
	cfg := &LogConfig{
		LogLevel:  log.Debug,
		LogFormat: JSONFormat,
		Name:      "test-system",
	}
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Warn("test warn msg")

	var jsonOutput map[string]string
	err = json.Unmarshal(buf.Bytes(), &jsonOutput)
	require.NoError(t, err)
	require.Contains(t, jsonOutput, "@level")
	require.Equal(t, jsonOutput["@level"], "warn")
	require.Contains(t, jsonOutput, "@message")
	require.Equal(t, jsonOutput["@message"], "test warn msg")
}

func TestLogger_SetupLoggerWithValidLogPath(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := &LogConfig{
		LogLevel:    log.Info,
		LogFilePath: tmpDir, //+ "/",
	}
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestLogger_SetupLoggerWithInValidLogPath(t *testing.T) {
	cfg := &LogConfig{
		LogLevel:    log.Info,
		LogFilePath: "nonexistentdir/",
	}
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.Error(t, err)
	require.True(t, errors.Is(err, os.ErrNotExist))
	require.Nil(t, logger)
}

func TestLogger_SetupLoggerWithInValidLogPathPermission(t *testing.T) {
	tmpDir := "/tmp/" + t.Name()

	err := os.Mkdir(tmpDir, 0o000)
	assert.NoError(t, err, "unexpected error testing with invalid log path permission")
	defer os.RemoveAll(tmpDir)

	cfg := &LogConfig{
		LogLevel:    log.Info,
		LogFilePath: tmpDir + "/",
	}
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.Error(t, err)
	require.True(t, errors.Is(err, os.ErrPermission))
	require.Nil(t, logger)
}

func TestLogger_SetupLoggerWithInvalidLogFilePath(t *testing.T) {
	cases := map[string]struct {
		path    string
		message string
	}{
		"file name *": {
			path:    "/this/isnt/ok/juan*.log",
			message: "file name contains globbing character",
		},
		"file name ?": {
			path:    "/this/isnt/ok/juan?.log",
			message: "file name contains globbing character",
		},
		"file name [": {
			path:    "/this/isnt/ok/[juan].log",
			message: "file name contains globbing character",
		},
		"directory path *": {
			path:    "/this/isnt/ok/*/qwerty.log",
			message: "directory contains glob character",
		},
		"directory path ?": {
			path:    "/this/isnt/ok/?/qwerty.log",
			message: "directory contains glob character",
		},
		"directory path [": {
			path:    "/this/isnt/ok/[foo]/qwerty.log",
			message: "directory contains glob character",
		},
	}

	for name, tc := range cases {
		name := name
		tc := tc
		cfg := &LogConfig{
			LogLevel:    log.Info,
			LogFilePath: tc.path,
		}
		_, err := Setup(cfg, &bytes.Buffer{})
		assert.Error(t, err, "%s: expected error due to *", name)
		assert.Contains(t, err.Error(), tc.message, "%s: error message does not match: %s", name, err.Error())
	}
}

func TestLogger_ChangeLogLevels(t *testing.T) {
	cfg := &LogConfig{
		LogLevel: log.Debug,
		Name:     "test-system",
	}
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, log.Debug, logger.GetLevel())

	// Create new named loggers from the base logger and change the levels
	logger2 := logger.Named("test2")
	logger3 := logger.Named("test3")

	logger2.SetLevel(log.Info)
	logger3.SetLevel(log.Error)

	assert.Equal(t, log.Debug, logger.GetLevel())
	assert.Equal(t, log.Info, logger2.GetLevel())
	assert.Equal(t, log.Error, logger3.GetLevel())
}
