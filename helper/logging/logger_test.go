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
	cfg := NewLogConfig("test-system", log.Info, StandardFormat, t.TempDir()+"test.log")

	logger, err := Setup(cfg, nil)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestLogger_SetupInvalidLogLevel(t *testing.T) {
	cfg := NewLogConfig("test-system", 999, StandardFormat, t.TempDir()+"test.log")

	_, err := Setup(cfg, nil)
	assert.Containsf(t, err.Error(), "invalid log level", "expected error %s", err)
}

func TestLogger_SetupLoggerErrorLevel(t *testing.T) {
	cfg := NewLogConfig("test-system", log.Error, StandardFormat, t.TempDir()+"test.log")
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
	cfg := NewLogConfig("test-system", log.Debug, StandardFormat, t.TempDir()+"test.log")
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

func TestLogger_SetupLoggerWithName(t *testing.T) {
	cfg := NewLogConfig("test-system", log.Debug, StandardFormat, t.TempDir()+"test.log")
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)

	logger.Warn("test warn msg")

	require.Contains(t, buf.String(), "[WARN]  test-system: test warn msg")
}

func TestLogger_SetupLoggerWithJSON(t *testing.T) {
	cfg := NewLogConfig("test-system", log.Debug, JSONFormat, t.TempDir()+"test.log")
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
	cfg := NewLogConfig("test-system", log.Info, StandardFormat, tmpDir+"/")
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.NoError(t, err)
	require.NotNil(t, logger)
}

func TestLogger_SetupLoggerWithInValidLogPath(t *testing.T) {
	cfg := NewLogConfig("test-system", log.Info, StandardFormat, "nonexistentdir/")
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.Error(t, err)
	require.True(t, errors.Is(err, os.ErrNotExist))
	require.Nil(t, logger)
}

func TestLogger_SetupLoggerWithInValidLogPathPermission(t *testing.T) {
	tmpDir := "/tmp/" + t.Name()

	os.Mkdir(tmpDir, 0o000)
	defer os.RemoveAll(tmpDir)
	cfg := NewLogConfig("test-system", log.Info, StandardFormat, tmpDir+"/")
	var buf bytes.Buffer

	logger, err := Setup(cfg, &buf)
	require.Error(t, err)
	require.True(t, errors.Is(err, os.ErrPermission))
	require.Nil(t, logger)
}
