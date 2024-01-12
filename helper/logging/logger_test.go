// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logging

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger_SetupBasic(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info

	logger, err := Setup(cfg, nil)
	require.NoError(t, err)
	require.NotNil(t, logger)
	require.Equal(t, logger.Name(), "test-system")
	require.True(t, logger.IsInfo())
}

func TestLogger_SetupInvalidLogLevel(t *testing.T) {
	cfg := newTestLogConfig(t)

	_, err := Setup(cfg, nil)
	assert.Containsf(t, err.Error(), "invalid log level", "expected error %s", err)
}

func TestLogger_SetupLoggerErrorLevel(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Error

	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Error("test error msg")
	logger.Info("test info msg")

	output := buf.String()

	require.Contains(t, output, "[ERROR] test-system: test error msg")
	require.NotContains(t, output, "[INFO] test-system: test info msg")
}

func TestLogger_SetupLoggerDebugLevel(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Debug
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("test info msg")
	logger.Debug("test debug msg")

	output := buf.String()

	require.Contains(t, output, "[INFO]  test-system: test info msg")
	require.Contains(t, output, "[DEBUG] test-system: test debug msg")
}

func TestLogger_SetupLoggerWithoutName(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.Name = ""
	cfg.LogLevel = hclog.Info
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Warn("test warn msg")

	require.Contains(t, buf.String(), "[WARN]  test warn msg")
}

func TestLogger_SetupLoggerWithJSON(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Debug
	cfg.LogFormat = JSONFormat
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

func TestLogger_SetupLoggerWithValidLogPathMissingFileName(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info
	cfg.LogFilePath = tmpDir + "/" // add the trailing slash to the temp dir
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("juan?")

	m, err := filepath.Glob(cfg.LogFilePath + "*")
	require.NoError(t, err)
	require.Truef(t, len(m) == 1, "no files were found")
}

func TestLogger_SetupLoggerWithValidLogPathFileName(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info
	cfg.LogFilePath = filepath.Join(tmpDir, "juan.log")
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("juan?")
	f, err := os.Stat(cfg.LogFilePath)
	require.NoError(t, err)
	require.NotNil(t, f)
}

func TestLogger_SetupLoggerWithValidLogPathFileNameRotate(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info
	cfg.LogFilePath = filepath.Join(tmpDir, "juan.log")
	cfg.LogRotateBytes = 1 // set a tiny number of bytes to force rotation
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Info("juan?")
	logger.Info("john?")
	f, err := os.Stat(cfg.LogFilePath)
	require.NoError(t, err)
	require.NotNil(t, f)
	m, err := filepath.Glob(tmpDir + "/juan-*") // look for juan-{timestamp}.log
	require.NoError(t, err)
	require.Truef(t, len(m) == 1, "no files were found")
}

func TestLogger_SetupLoggerWithValidLogPath(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info
	cfg.LogFilePath = tmpDir + "/" // add the trailing slash to the temp dir
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestLogger_SetupLoggerWithInValidLogPath(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info
	cfg.LogLevel = hclog.Info
	cfg.LogFilePath = "nonexistentdir/"
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

	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Info
	cfg.LogFilePath = tmpDir + "/"
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
		cfg := newTestLogConfig(t)
		cfg.LogLevel = hclog.Info
		cfg.LogFilePath = tc.path

		_, err := Setup(cfg, &bytes.Buffer{})
		assert.Error(t, err, "%s: expected error due to *", name)
		assert.Contains(t, err.Error(), tc.message, "%s: error message does not match: %s", name, err.Error())
	}
}

func TestLogger_ChangeLogLevels(t *testing.T) {
	cfg := newTestLogConfig(t)
	cfg.LogLevel = hclog.Debug
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	assert.Equal(t, hclog.Debug, logger.GetLevel())

	// Create new named loggers from the base logger and change the levels
	logger2 := logger.Named("test2")
	logger3 := logger.Named("test3")

	logger2.SetLevel(hclog.Info)
	logger3.SetLevel(hclog.Error)

	assert.Equal(t, hclog.Debug, logger.GetLevel())
	assert.Equal(t, hclog.Info, logger2.GetLevel())
	assert.Equal(t, hclog.Error, logger3.GetLevel())
}

func newTestLogConfig(t *testing.T) *LogConfig {
	t.Helper()

	cfg, err := NewLogConfig("test")
	require.NoError(t, err)
	cfg.Name = "test-system"

	return cfg
}
