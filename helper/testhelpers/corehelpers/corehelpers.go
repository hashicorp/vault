// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package corehelpers contains testhelpers that don't depend on package vault,
// and thus can be used within vault (as well as elsewhere.)
package corehelpers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/go-testing-interface"
)

// RetryUntil runs f until it returns a nil result or the timeout is reached.
// If a nil result hasn't been obtained by timeout, calls t.Fatal.
func RetryUntil(t testing.T, timeout time.Duration, f func() error) {
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
func MakeTestPluginDir(t testing.T) (string, func(t testing.T)) {
	if t != nil {
		t.Helper()
	}

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		if t == nil {
			panic(err)
		}
		t.Fatal(err)
	}

	// OSX tempdir are /var, but actually symlinked to /private/var
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		if t == nil {
			panic(err)
		}
		t.Fatal(err)
	}

	return dir, func(t testing.T) {
		if err := os.RemoveAll(dir); err != nil {
			if t == nil {
				panic(err)
			}
			t.Fatal(err)
		}
	}
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
	}
	return []string{}
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

func TestNoopAudit(t testing.T, config map[string]string) *NoopAudit {
	n, err := NewNoopAudit(config)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func NewNoopAudit(config map[string]string) (*NoopAudit, error) {
	view := &logical.InmemStorage{}
	err := view.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	if err != nil {
		return nil, err
	}

	n := &NoopAudit{
		Config: &audit.BackendConfig{
			SaltView: view,
			SaltConfig: &salt.Config{
				HMAC:     sha256.New,
				HMACType: "hmac-sha256",
			},
			Config: config,
		},
	}
	n.formatter.AuditFormatWriter = &audit.JSONFormatWriter{
		SaltFunc: n.Salt,
	}
	return n, nil
}

func NoopAuditFactory(records **[][]byte) audit.Factory {
	return func(_ context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		n, err := NewNoopAudit(config.Config)
		if err != nil {
			return nil, err
		}
		if records != nil {
			*records = &n.records
		}
		return n, nil
	}
}

type NoopAudit struct {
	Config         *audit.BackendConfig
	ReqErr         error
	ReqAuth        []*logical.Auth
	Req            []*logical.Request
	ReqHeaders     []map[string][]string
	ReqNonHMACKeys []string
	ReqErrs        []error

	RespErr            error
	RespAuth           []*logical.Auth
	RespReq            []*logical.Request
	Resp               []*logical.Response
	RespNonHMACKeys    [][]string
	RespReqNonHMACKeys [][]string
	RespErrs           []error

	formatter audit.AuditFormatter
	records   [][]byte
	l         sync.RWMutex
	salt      *salt.Salt
	saltMutex sync.RWMutex
}

func (n *NoopAudit) LogRequest(ctx context.Context, in *logical.LogInput) error {
	n.l.Lock()
	defer n.l.Unlock()
	if n.formatter.AuditFormatWriter != nil {
		var w bytes.Buffer
		err := n.formatter.FormatRequest(ctx, &w, audit.FormatterConfig{}, in)
		if err != nil {
			return err
		}
		n.records = append(n.records, w.Bytes())
	}

	n.ReqAuth = append(n.ReqAuth, in.Auth)
	n.Req = append(n.Req, in.Request)
	n.ReqHeaders = append(n.ReqHeaders, in.Request.Headers)
	n.ReqNonHMACKeys = in.NonHMACReqDataKeys
	n.ReqErrs = append(n.ReqErrs, in.OuterErr)

	return n.ReqErr
}

func (n *NoopAudit) LogResponse(ctx context.Context, in *logical.LogInput) error {
	n.l.Lock()
	defer n.l.Unlock()

	if n.formatter.AuditFormatWriter != nil {
		var w bytes.Buffer
		err := n.formatter.FormatResponse(ctx, &w, audit.FormatterConfig{}, in)
		if err != nil {
			return err
		}
		n.records = append(n.records, w.Bytes())
	}

	n.RespAuth = append(n.RespAuth, in.Auth)
	n.RespReq = append(n.RespReq, in.Request)
	n.Resp = append(n.Resp, in.Response)
	n.RespErrs = append(n.RespErrs, in.OuterErr)

	if in.Response != nil {
		n.RespNonHMACKeys = append(n.RespNonHMACKeys, in.NonHMACRespDataKeys)
		n.RespReqNonHMACKeys = append(n.RespReqNonHMACKeys, in.NonHMACReqDataKeys)
	}

	return n.RespErr
}

func (n *NoopAudit) LogTestMessage(ctx context.Context, in *logical.LogInput, config map[string]string) error {
	n.l.Lock()
	defer n.l.Unlock()
	var w bytes.Buffer
	tempFormatter := audit.NewTemporaryFormatter(config["format"], config["prefix"])
	err := tempFormatter.FormatResponse(ctx, &w, audit.FormatterConfig{}, in)
	if err != nil {
		return err
	}
	n.records = append(n.records, w.Bytes())
	return nil
}

func (n *NoopAudit) Salt(ctx context.Context) (*salt.Salt, error) {
	n.saltMutex.RLock()
	if n.salt != nil {
		defer n.saltMutex.RUnlock()
		return n.salt, nil
	}
	n.saltMutex.RUnlock()
	n.saltMutex.Lock()
	defer n.saltMutex.Unlock()
	if n.salt != nil {
		return n.salt, nil
	}
	salt, err := salt.NewSalt(ctx, n.Config.SaltView, n.Config.SaltConfig)
	if err != nil {
		return nil, err
	}
	n.salt = salt
	return salt, nil
}

func (n *NoopAudit) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := n.Salt(ctx)
	if err != nil {
		return "", err
	}
	return salt.GetIdentifiedHMAC(data), nil
}

func (n *NoopAudit) Reload(ctx context.Context) error {
	return nil
}

func (n *NoopAudit) Invalidate(ctx context.Context) {
	n.saltMutex.Lock()
	defer n.saltMutex.Unlock()
	n.salt = nil
}

type TestLogger struct {
	hclog.Logger
	Path string
	File *os.File
	sink hclog.SinkAdapter
}

func NewTestLogger(t testing.T) *TestLogger {
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
		Output:            ioutil.Discard,
		IndependentLevels: true,
	})
	sink := hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Output:            output,
		Level:             hclog.Trace,
		IndependentLevels: true,
	})
	logger.RegisterSink(sink)
	return &TestLogger{
		Path:   logPath,
		File:   logFile,
		Logger: logger,
		sink:   sink,
	}
}

func (tl *TestLogger) StopLogging() {
	tl.Logger.(hclog.InterceptLogger).DeregisterSink(tl.sink)
}
