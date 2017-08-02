package plugin

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestLogger_impl(t *testing.T) {
	var _ log.Logger = new(LoggerClient)
}

func TestLogger_levels(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	l := logformat.NewVaultLoggerWithWriter(writer, log.LevelTrace)

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	expected := "foobar"
	testLogger := &LoggerClient{client: client}

	// Test trace
	testLogger.Trace(expected)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result := buf.String()
	buf.Reset()
	if !strings.Contains(result, expected) {
		t.Fatalf("expected log to contain %s, got %s", expected, result)
	}

	// Test debug
	testLogger.Debug(expected)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result = buf.String()
	buf.Reset()
	if !strings.Contains(result, expected) {
		t.Fatalf("expected log to contain %s, got %s", expected, result)
	}

	// Test debug
	testLogger.Info(expected)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result = buf.String()
	buf.Reset()
	if !strings.Contains(result, expected) {
		t.Fatalf("expected log to contain %s, got %s", expected, result)
	}

	// Test warn
	testLogger.Warn(expected)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result = buf.String()
	buf.Reset()
	if !strings.Contains(result, expected) {
		t.Fatalf("expected log to contain %s, got %s", expected, result)
	}

	// Test error
	testLogger.Error(expected)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result = buf.String()
	buf.Reset()
	if !strings.Contains(result, expected) {
		t.Fatalf("expected log to contain %s, got %s", expected, result)
	}

	// Test fatal
	testLogger.Fatal(expected)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result = buf.String()
	buf.Reset()
	if result != "" {
		t.Fatalf("expected log Fatal() to be no-op, got %s", result)
	}
}

func TestLogger_isLevels(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	l := logformat.NewVaultLoggerWithWriter(ioutil.Discard, log.LevelAll)

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	testLogger := &LoggerClient{client: client}

	if !testLogger.IsDebug() || !testLogger.IsInfo() || !testLogger.IsTrace() || !testLogger.IsWarn() {
		t.Fatal("expected logger to return true for all logger level checks")
	}
}

func TestLogger_log(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	l := logformat.NewVaultLoggerWithWriter(writer, log.LevelTrace)

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	expected := "foobar"
	testLogger := &LoggerClient{client: client}

	// Test trace
	testLogger.Log(log.LevelInfo, expected, nil)
	if err := writer.Flush(); err != nil {
		t.Fatal(err)
	}
	result := buf.String()
	if !strings.Contains(result, expected) {
		t.Fatalf("expected log to contain %s, got %s", expected, result)
	}

}

func TestLogger_setLevel(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	l := log.NewLogger(ioutil.Discard, "test-logger")

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	testLogger := &LoggerClient{client: client}
	testLogger.SetLevel(log.LevelWarn)

	if !testLogger.IsWarn() {
		t.Fatal("expected logger to support warn level")
	}
}
