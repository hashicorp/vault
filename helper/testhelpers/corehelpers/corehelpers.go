// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package corehelpers contains testhelpers that don't depend on package vault,
// and thus can be used within vault (as well as elsewhere.)
package corehelpers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/go-testing-interface"
)

var externalPlugins = []string{"transform", "kmip", "keymgmt"}

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

	cfg, err := audit.NewFormatterConfig()
	if err != nil {
		return nil, err
	}

	f, err := audit.NewEntryFormatter(cfg, n)
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}

	fw, err := audit.NewEntryFormatterWriter(cfg, f, &audit.JSONWriter{})
	if err != nil {
		return nil, fmt.Errorf("error creating formatter writer: %w", err)
	}

	n.formatter = fw

	n.nodeIDList = make([]eventlogger.NodeID, 2)
	n.nodeMap = make(map[eventlogger.NodeID]eventlogger.Node)

	formatterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating random NodeID for formatter node: %w", err)
	}

	n.nodeIDList[0] = formatterNodeID
	n.nodeMap[formatterNodeID] = f

	sinkNode := event.NewNoopSink()
	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating random NodeID for sink node: %w", err)
	}

	n.nodeIDList[1] = sinkNodeID
	n.nodeMap[sinkNodeID] = sinkNode

	return n, nil
}

func NoopAuditFactory(records **[][]byte) audit.Factory {
	return func(_ context.Context, config *audit.BackendConfig, _ bool, _ audit.HeaderFormatter) (audit.Backend, error) {
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

	formatter *audit.EntryFormatterWriter
	records   [][]byte
	l         sync.RWMutex
	salt      *salt.Salt
	saltMutex sync.RWMutex

	nodeIDList []eventlogger.NodeID
	nodeMap    map[eventlogger.NodeID]eventlogger.Node
}

func (n *NoopAudit) LogRequest(ctx context.Context, in *logical.LogInput) error {
	n.l.Lock()
	defer n.l.Unlock()
	if n.formatter != nil {
		var w bytes.Buffer
		err := n.formatter.FormatAndWriteRequest(ctx, &w, in)
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

	if n.formatter != nil {
		var w bytes.Buffer
		err := n.formatter.FormatAndWriteResponse(ctx, &w, in)
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

	tempFormatter, err := audit.NewTemporaryFormatter(config["format"], config["prefix"])
	if err != nil {
		return err
	}

	err = tempFormatter.FormatAndWriteResponse(ctx, &w, in)
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
	s, err := salt.NewSalt(ctx, n.Config.SaltView, n.Config.SaltConfig)
	if err != nil {
		return nil, err
	}
	n.salt = s
	return s, nil
}

func (n *NoopAudit) GetHash(ctx context.Context, data string) (string, error) {
	s, err := n.Salt(ctx)
	if err != nil {
		return "", err
	}
	return s.GetIdentifiedHMAC(data), nil
}

func (n *NoopAudit) Reload(_ context.Context) error {
	return nil
}

func (n *NoopAudit) Invalidate(_ context.Context) {
	n.saltMutex.Lock()
	defer n.saltMutex.Unlock()
	n.salt = nil
}

// RegisterNodesAndPipeline registers the nodes and a pipeline as required by
// the audit.Backend interface.
func (b *NoopAudit) RegisterNodesAndPipeline(broker *eventlogger.Broker, name string) error {
	for id, node := range b.nodeMap {
		if err := broker.RegisterNode(id, node); err != nil {
			return err
		}
	}

	pipeline := eventlogger.Pipeline{
		PipelineID: eventlogger.PipelineID(name),
		EventType:  eventlogger.EventType("audit"),
		NodeIDs:    b.nodeIDList,
	}

	return broker.RegisterPipeline(pipeline)
}

type TestLogger struct {
	hclog.InterceptLogger
	Path string
	File *os.File
	sink hclog.SinkAdapter
	// For managing temporary start-up state
	sync.RWMutex
	AllLoggers []hclog.Logger
	logging.SubloggerAdder
}

// RegisterSubloggerAdder checks to see if the provided logger interface is a
// TestLogger and re-assigns the SubloggerHook implementation if so.
func RegisterSubloggerAdder(logger hclog.Logger, adder logging.SubloggerAdder) {
	if l, ok := logger.(*TestLogger); ok {
		l.Lock()
		l.SubloggerAdder = adder
		l.Unlock()
	}
}

// AppendToAllLoggers appends the sub logger to allLoggers, or if the TestLogger
// is assigned to a SubloggerAdder implementation, it calls the underlying hook.
func (l *TestLogger) AppendToAllLoggers(sub hclog.Logger) hclog.Logger {
	l.Lock()
	defer l.Unlock()
	if l.SubloggerAdder == nil {
		l.AllLoggers = append(l.AllLoggers, sub)
		return sub
	}
	return l.SubloggerHook(sub)
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

	sink := hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Output:            output,
		Level:             hclog.Trace,
		IndependentLevels: true,
	})

	testLogger := &TestLogger{
		Path: logPath,
		File: logFile,
		sink: sink,
	}

	// We send nothing on the regular logger, that way we can later deregister
	// the sink to stop logging during cluster cleanup.
	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Output:            io.Discard,
		IndependentLevels: true,
		Name:              t.Name(),
		SubloggerHook:     testLogger.AppendToAllLoggers,
	})

	logger.RegisterSink(sink)
	testLogger.InterceptLogger = logger

	return testLogger
}

func (tl *TestLogger) StopLogging() {
	tl.InterceptLogger.DeregisterSink(tl.sink)
}
