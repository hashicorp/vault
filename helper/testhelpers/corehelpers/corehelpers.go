// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

// Package corehelpers contains testhelpers that don't depend on package vault,
// and thus can be used within vault (as well as elsewhere.)
package corehelpers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/plugins/database/mysql"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/go-testing-interface"
)

var (
	_ audit.Backend    = (*NoopAudit)(nil)
	_ eventlogger.Node = (*noopWrapper)(nil)
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
func MakeTestPluginDir(t testing.T) string {
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

func TestNoopAudit(t testing.T, path string, config map[string]string, opts ...audit.Option) *NoopAudit {
	cfg := &audit.BackendConfig{Config: config, MountPath: path}
	n, err := NewNoopAudit(cfg, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

// NewNoopAudit should be used to create a NoopAudit as it handles creation of a
// predictable salt and wraps eventlogger nodes so information can be retrieved on
// what they've seen or formatted.
func NewNoopAudit(config *audit.BackendConfig, opts ...audit.Option) (*NoopAudit, error) {
	view := &logical.InmemStorage{}

	// Create the salt with a known key for predictable hmac values.
	se := &logical.StorageEntry{Key: "salt", Value: []byte("foo")}
	err := view.Put(context.Background(), se)
	if err != nil {
		return nil, err
	}

	// Override the salt related config settings.
	backendConfig := &audit.BackendConfig{
		SaltView: view,
		SaltConfig: &salt.Config{
			HMAC:     sha256.New,
			HMACType: "hmac-sha256",
		},
		Config:    config.Config,
		MountPath: config.MountPath,
	}

	n := &NoopAudit{Config: backendConfig}

	cfg, err := audit.NewFormatterConfig()
	if err != nil {
		return nil, err
	}

	f, err := audit.NewEntryFormatter(cfg, n, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}

	n.nodeIDList = make([]eventlogger.NodeID, 2)
	n.nodeMap = make(map[eventlogger.NodeID]eventlogger.Node, 2)

	formatterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating random NodeID for formatter node: %w", err)
	}

	// Wrap the formatting node, so we can get any bytes that were formatted etc.
	wrappedFormatter := &noopWrapper{format: "json", node: f, backend: n}

	n.nodeIDList[0] = formatterNodeID
	n.nodeMap[formatterNodeID] = wrappedFormatter

	sinkNode := event.NewNoopSink()
	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating random NodeID for sink node: %w", err)
	}

	n.nodeIDList[1] = sinkNodeID
	n.nodeMap[sinkNodeID] = sinkNode

	return n, nil
}

// NoopAuditFactory should be used when the test needs a way to access bytes that
// have been formatted by the pipeline during audit requests.
// The records parameter will be repointed to the one used within the pipeline.
func NoopAuditFactory(records **[][]byte) audit.Factory {
	return func(_ context.Context, config *audit.BackendConfig, _ bool, headerFormatter audit.HeaderFormatter) (audit.Backend, error) {
		n, err := NewNoopAudit(config, audit.WithHeaderFormatter(headerFormatter))
		if err != nil {
			return nil, err
		}
		if records != nil {
			*records = &n.records
		}

		return n, nil
	}
}

// noopWrapper is designed to wrap a formatter node in order to allow access to
// bytes formatted, headers formatted and parts of the logical.LogInput.
// Some older tests relied on being able to query this information so while those
// tests stick around we should look after them.
type noopWrapper struct {
	format  string
	node    eventlogger.Node
	backend *NoopAudit
}

type NoopAudit struct {
	Config *audit.BackendConfig

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
	records            [][]byte
	l                  sync.RWMutex
	salt               *salt.Salt
	saltMutex          sync.RWMutex

	nodeIDList []eventlogger.NodeID
	nodeMap    map[eventlogger.NodeID]eventlogger.Node
}

// Process handles the contortions required by older test code to ensure behavior.
// It will attempt to do some pre/post processing of the logical.LogInput that should
// form part of the event's payload data, as well as capturing the resulting headers
// that were formatted and track the overall bytes that a formatted event uses when
// it's ready to head down the pipeline to the sink node (a noop for us).
func (n *noopWrapper) Process(ctx context.Context, e *eventlogger.Event) (*eventlogger.Event, error) {
	n.backend.l.Lock()
	defer n.backend.l.Unlock()

	var err error

	// We're expecting audit events since this is an audit device.
	a, ok := e.Payload.(*audit.AuditEvent)
	if !ok {
		return nil, errors.New("cannot parse payload as an audit event")
	}

	in := a.Data

	// Depending on the type of the audit event (request or response) we need to
	// track different things.
	switch a.Subtype {
	case audit.RequestType:
		n.backend.ReqAuth = append(n.backend.ReqAuth, in.Auth)
		n.backend.Req = append(n.backend.Req, in.Request)
		n.backend.ReqNonHMACKeys = in.NonHMACReqDataKeys
		n.backend.ReqErrs = append(n.backend.ReqErrs, in.OuterErr)

		if n.backend.ReqErr != nil {
			return nil, n.backend.ReqErr
		}
	case audit.ResponseType:
		n.backend.RespAuth = append(n.backend.RespAuth, in.Auth)
		n.backend.RespReq = append(n.backend.RespReq, in.Request)
		n.backend.Resp = append(n.backend.Resp, in.Response)
		n.backend.RespErrs = append(n.backend.RespErrs, in.OuterErr)

		if in.Response != nil {
			n.backend.RespNonHMACKeys = append(n.backend.RespNonHMACKeys, in.NonHMACRespDataKeys)
			n.backend.RespReqNonHMACKeys = append(n.backend.RespReqNonHMACKeys, in.NonHMACReqDataKeys)
		}

		if n.backend.RespErr != nil {
			return nil, n.backend.RespErr
		}
	default:
		return nil, fmt.Errorf("unknown audit event type: %q", a.Subtype)
	}

	// Once we've taken note of the relevant properties of the event, we get the
	// underlying (wrapped) node to process it as normal.
	e, err = n.node.Process(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("error processing wrapped node: %w", err)
	}

	// Once processing has been carried out, the underlying node (a formatter node)
	// should contain the output ready for the sink node. We'll get that in order
	// to track how many bytes we formatted.
	b, ok := e.Format(n.format)
	if ok {
		n.backend.records = append(n.backend.records, b)
	}

	// Finally, the last bit of post-processing is to make sure that we track the
	// formatted headers that would have made it to the logs via the sink node.
	// They only appear in requests.
	if a.Subtype == audit.RequestType {
		reqEntry := &audit.RequestEntry{}
		err = json.Unmarshal(b, &reqEntry)
		if err != nil {
			return nil, fmt.Errorf("unable to parse formatted audit entry data: %w", err)
		}

		n.backend.ReqHeaders = append(n.backend.ReqHeaders, reqEntry.Request.Headers)
	}

	// Return the event and no error in order to let the pipeline continue on.
	return e, nil
}

func (n *noopWrapper) Reopen() error {
	return n.node.Reopen()
}

func (n *noopWrapper) Type() eventlogger.NodeType {
	return n.node.Type()
}

// Deprecated: use eventlogger.
func (n *NoopAudit) LogRequest(ctx context.Context, in *logical.LogInput) error {
	return nil
}

// Deprecated: use eventlogger.
func (n *NoopAudit) LogResponse(ctx context.Context, in *logical.LogInput) error {
	return nil
}

// LogTestMessage will manually crank the handle on the nodes associated with this backend.
func (n *NoopAudit) LogTestMessage(ctx context.Context, in *logical.LogInput, config map[string]string) error {
	n.l.Lock()
	defer n.l.Unlock()

	// Fake event for test purposes.
	e := &eventlogger.Event{
		Type:      eventlogger.EventType(event.AuditType.String()),
		CreatedAt: time.Now(),
		Formatted: make(map[string][]byte),
		Payload:   in,
	}

	// Try to get the required format from config and default to JSON.
	format, ok := config["format"]
	if !ok {
		format = "json"
	}
	cfg, err := audit.NewFormatterConfig(audit.WithFormat(format))
	if err != nil {
		return fmt.Errorf("cannot create config for formatter node: %w", err)
	}
	// Create a temporary formatter node for reuse.
	f, err := audit.NewEntryFormatter(cfg, n, audit.WithPrefix(config["prefix"]))

	// Go over each node in order from our list.
	for _, id := range n.nodeIDList {
		node, ok := n.nodeMap[id]
		if !ok {
			return fmt.Errorf("node not found: %v", id)
		}

		switch node.Type() {
		case eventlogger.NodeTypeFormatter:
			// Use a temporary formatter node which doesn't persist its salt anywhere.
			if formatNode, ok := node.(*audit.EntryFormatter); ok && formatNode != nil {
				e, err = f.Process(ctx, e)

				// Housekeeping, we should update that we processed some bytes.
				if e != nil {
					b, ok := e.Format(format)
					if ok {
						n.records = append(n.records, b)
					}
				}
			}
		default:
			e, err = node.Process(ctx, e)
		}
	}

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
func (n *NoopAudit) RegisterNodesAndPipeline(broker *eventlogger.Broker, name string) error {
	for id, node := range n.nodeMap {
		if err := broker.RegisterNode(id, node); err != nil {
			return err
		}
	}

	pipeline := eventlogger.Pipeline{
		PipelineID: eventlogger.PipelineID(name),
		EventType:  eventlogger.EventType(event.AuditType.String()),
		NodeIDs:    n.nodeIDList,
	}

	return broker.RegisterPipeline(pipeline)
}

type TestLogger struct {
	hclog.InterceptLogger
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

func (n *NoopAudit) EventType() eventlogger.EventType {
	return eventlogger.EventType(event.AuditType.String())
}

func (n *NoopAudit) HasFiltering() bool {
	return false
}

func (n *NoopAudit) Name() string {
	return n.Config.MountPath
}

func (n *NoopAudit) Nodes() map[eventlogger.NodeID]eventlogger.Node {
	return n.nodeMap
}

func (n *NoopAudit) NodeIDs() []eventlogger.NodeID {
	return n.nodeIDList
}

func (n *NoopAudit) IsFallback() bool {
	return false
}
