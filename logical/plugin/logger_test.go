package plugin

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/rpc"
	"strings"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/logging"
)

func TestLogger_levels(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	l := logging.NewVaultLoggerWithWriter(writer, hclog.Trace)

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	expected := "foobar"
	testLogger := &deprecatedLoggerClient{client: client}

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

	l := logging.NewVaultLoggerWithWriter(ioutil.Discard, hclog.Trace)

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	testLogger := &deprecatedLoggerClient{client: client}

	if !testLogger.IsDebug() || !testLogger.IsInfo() || !testLogger.IsTrace() || !testLogger.IsWarn() {
		t.Fatal("expected logger to return true for all logger level checks")
	}
}

func TestLogger_log(t *testing.T) {
	client, server := plugin.TestRPCConn(t)
	defer client.Close()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	l := logging.NewVaultLoggerWithWriter(writer, hclog.Trace)

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	expected := "foobar"
	testLogger := &deprecatedLoggerClient{client: client}

	// Test trace 6 = logxi.LevelInfo
	testLogger.Log(6, expected, nil)
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

	l := hclog.New(&hclog.LoggerOptions{Output: ioutil.Discard})

	server.RegisterName("Plugin", &LoggerServer{
		logger: l,
	})

	testLogger := &deprecatedLoggerClient{client: client}
	testLogger.SetLevel(4) // 4 == logxi.LevelWarn

	if !testLogger.IsWarn() {
		t.Fatal("expected logger to support warn level")
	}
}

type deprecatedLoggerClient struct {
	client *rpc.Client
}

func (l *deprecatedLoggerClient) Trace(msg string, args ...interface{}) {
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	l.client.Call("Plugin.Trace", cArgs, &struct{}{})
}

func (l *deprecatedLoggerClient) Debug(msg string, args ...interface{}) {
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	l.client.Call("Plugin.Debug", cArgs, &struct{}{})
}

func (l *deprecatedLoggerClient) Info(msg string, args ...interface{}) {
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	l.client.Call("Plugin.Info", cArgs, &struct{}{})
}
func (l *deprecatedLoggerClient) Warn(msg string, args ...interface{}) error {
	var reply LoggerReply
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	err := l.client.Call("Plugin.Warn", cArgs, &reply)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}

	return nil
}
func (l *deprecatedLoggerClient) Error(msg string, args ...interface{}) error {
	var reply LoggerReply
	cArgs := &LoggerArgs{
		Msg:  msg,
		Args: args,
	}
	err := l.client.Call("Plugin.Error", cArgs, &reply)
	if err != nil {
		return err
	}
	if reply.Error != nil {
		return reply.Error
	}

	return nil
}

func (l *deprecatedLoggerClient) Fatal(msg string, args ...interface{}) {
	// NOOP since it's not actually used within vault
	return
}

func (l *deprecatedLoggerClient) Log(level int, msg string, args []interface{}) {
	cArgs := &LoggerArgs{
		Level: level,
		Msg:   msg,
		Args:  args,
	}
	l.client.Call("Plugin.Log", cArgs, &struct{}{})
}

func (l *deprecatedLoggerClient) SetLevel(level int) {
	l.client.Call("Plugin.SetLevel", level, &struct{}{})
}

func (l *deprecatedLoggerClient) IsTrace() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsTrace", new(interface{}), &reply)
	return reply.IsTrue
}
func (l *deprecatedLoggerClient) IsDebug() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsDebug", new(interface{}), &reply)
	return reply.IsTrue
}

func (l *deprecatedLoggerClient) IsInfo() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsInfo", new(interface{}), &reply)
	return reply.IsTrue
}

func (l *deprecatedLoggerClient) IsWarn() bool {
	var reply LoggerReply
	l.client.Call("Plugin.IsWarn", new(interface{}), &reply)
	return reply.IsTrue
}
