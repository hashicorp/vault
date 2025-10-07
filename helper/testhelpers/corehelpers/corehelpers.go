// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package corehelpers contains testhelpers that don't depend on package vault,
// and thus can be used within vault (as well as elsewhere.)
package corehelpers

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

var externalPlugins = []string{"transform", "kmip", "keymgmt"}

// RetryUntil runs f until it returns a nil result or the timeout is reached.
// If a nil result hasn't been obtained by timeout, calls t.Fatal.
func RetryUntil(t testing.TB, timeout time.Duration, f func() error) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var err error
	for time.Now().Before(deadline) {
		if err = f(); err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("did not complete before deadline, err: %v", err)
}

// MakeTestPluginDir creates a temporary directory suitable for holding plugins.
// This helper also resolves symlinks to make tests happy on OS X.
func MakeTestPluginDir(t testing.TB) string {
	t.Helper()

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}

	// OSX tempdir are /var, but actually symlinked to /private/var
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
	})

	return dir
}

func NewMockBuiltinRegistry() *mockBuiltinRegistry {
	return &mockBuiltinRegistry{
		forTesting: map[string]mockBackend{
			"mysql-database-plugin":      {PluginType: consts.PluginTypeDatabase},
			"postgresql-database-plugin": {PluginType: consts.PluginTypeDatabase},
			"approle":                    {PluginType: consts.PluginTypeCredential},
			"pending-removal-test-plugin": {
				PluginType:        consts.PluginTypeCredential,
				DeprecationStatus: consts.PendingRemoval,
			},
			"aws":    {PluginType: consts.PluginTypeCredential},
			"consul": {PluginType: consts.PluginTypeSecrets},
		},
	}
}

type mockBackend struct {
	consts.PluginType
	consts.DeprecationStatus
}

type mockBuiltinRegistry struct {
	forTesting map[string]mockBackend
}

func toFunc(f logical.Factory) func() (interface{}, error) {
	return func() (interface{}, error) {
		return f, nil
	}
}

func (m *mockBuiltinRegistry) Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool) {
	testBackend, ok := m.forTesting[name]
	if !ok {
		return nil, false
	}
	testPluginType := testBackend.PluginType
	if pluginType != testPluginType {
		return nil, false
	}

	switch name {
	case "approle", "pending-removal-test-plugin":
		return toFunc(approle.Factory), true
	case "aws":
		return toFunc(func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
			b := new(framework.Backend)
			b.Setup(ctx, config)
			b.BackendType = logical.TypeCredential
			return b, nil
		}), true
	case "postgresql-database-plugin":
		return toFunc(func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
			b := new(framework.Backend)
			b.Setup(ctx, config)
			b.BackendType = logical.TypeLogical
			return b, nil
		}), true
	case "mysql-database-plugin":
		return mysql.New(mysql.DefaultUserNameTemplate), true
	case "consul":
		return toFunc(func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
			b := new(framework.Backend)
			b.Setup(ctx, config)
			b.BackendType = logical.TypeLogical
			return b, nil
		}), true
	default:
		return nil, false
	}
}

// Keys only supports getting a realistic list of the keys for database plugins,
// and approle
func (m *mockBuiltinRegistry) Keys(pluginType consts.PluginType) []string {
	switch pluginType {
	case consts.PluginTypeDatabase:
		// This is a hard-coded reproduction of the db plugin keys in
		// helper/builtinplugins/registry.go. The registry isn't directly used
		// because it causes import cycles.
		return []string{
			"mysql-database-plugin",
			"mysql-aurora-database-plugin",
			"mysql-rds-database-plugin",
			"mysql-legacy-database-plugin",

			"cassandra-database-plugin",
			"couchbase-database-plugin",
			"elasticsearch-database-plugin",
			"hana-database-plugin",
			"influxdb-database-plugin",
			"mongodb-database-plugin",
			"mongodbatlas-database-plugin",
			"mssql-database-plugin",
			"postgresql-database-plugin",
			"redis-elasticache-database-plugin",
			"redshift-database-plugin",
			"redis-database-plugin",
			"snowflake-database-plugin",
		}
	case consts.PluginTypeCredential:
		return []string{
			"pending-removal-test-plugin",
			"approle",
		}

	case consts.PluginTypeSecrets:
		return append(externalPlugins, "kv")
	}

	return []string{}
}

func (r *mockBuiltinRegistry) IsBuiltinEntPlugin(name string, pluginType consts.PluginType) bool {
	for _, i := range externalPlugins {
		if i == name {
			return true
		}
	}
	return false
}

func (m *mockBuiltinRegistry) Contains(name string, pluginType consts.PluginType) bool {
	for _, key := range m.Keys(pluginType) {
		if key == name {
			return true
		}
	}
	return false
}

func (m *mockBuiltinRegistry) DeprecationStatus(name string, pluginType consts.PluginType) (consts.DeprecationStatus, bool) {
	if m.Contains(name, pluginType) {
		return m.forTesting[name].DeprecationStatus, true
	}

	return consts.Unknown, false
}

type TestLogger struct {
	hclog.InterceptLogger
	Path string
	File *os.File
	sink hclog.SinkAdapter
}

func NewTestLogger(t testing.TB) *TestLogger {
	var logFile *os.File
	var logPath string
	output := os.Stderr

	logDir := os.Getenv("VAULT_TEST_LOG_DIR")
	if logDir != "" {
		logPath = filepath.Join(logDir, t.Name()+".log")
		// t.Name may include slashes.
		dir, _ := filepath.Split(logPath)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		logFile, err = os.Create(logPath)
		if err != nil {
			t.Fatal(err)
		}
		output = logFile
	}

	// We send nothing on the regular logger, that way we can later deregister
	// the sink to stop logging during cluster cleanup.
	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Output:            io.Discard,
		IndependentLevels: true,
		Name:              t.Name(),
	})
	sink := hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Output:            output,
		Level:             hclog.Trace,
		IndependentLevels: true,
	})
	logger.RegisterSink(sink)

	testLogger := &TestLogger{
		Path:            logPath,
		File:            logFile,
		InterceptLogger: logger,
		sink:            sink,
	}

	t.Cleanup(func() {
		testLogger.StopLogging()
		if t.Failed() {
			_ = testLogger.File.Close()
		} else {
			_ = os.Remove(testLogger.Path)
		}
	})
	return testLogger
}

func (tl *TestLogger) StopLogging() {
	tl.InterceptLogger.DeregisterSink(tl.sink)
}
