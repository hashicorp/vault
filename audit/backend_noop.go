// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/eventlogger"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	_ Backend          = (*NoopAudit)(nil)
	_ eventlogger.Node = (*noopWrapper)(nil)
)

// noopWrapper is designed to wrap a formatter node in order to allow access to
// bytes formatted, headers formatted and parts of the logical.LogInput.
// Some older tests relied on being able to query this information so while those
// tests stick around we should look after them.
type noopWrapper struct {
	format  string
	node    eventlogger.Node
	backend *NoopAudit
}

// SetListener provides a callback func to the NoopAudit which can be invoked
// during processing of the Event.
//
// Deprecated: SetListener should not be used in new tests.
func (n *NoopAudit) SetListener(listener func(event *Event)) {
	n.listener = listener
}

// NoopAudit only exists to allow legacy tests to continue working.
//
// Deprecated: NoopAudit should not be used in new tests.
type NoopAudit struct {
	Config *BackendConfig

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

	listener func(event *Event)
}

// noopHeaderFormatter can be used within no-op audit devices to do nothing when
// it comes to only allow configured headers to appear in the result.
// Whatever is passed in will be returned (nil becomes an empty map) in lowercase.
type noopHeaderFormatter struct{}

// ApplyConfig implements the relevant interface to make noopHeaderFormatter an HeaderFormatter.
func (f *noopHeaderFormatter) ApplyConfig(_ context.Context, headers map[string][]string, _ Salter) (result map[string][]string, retErr error) {
	if len(headers) < 1 {
		return map[string][]string{}, nil
	}

	// Make a copy of the incoming headers with everything lower so we can
	// case-insensitively compare
	lowerHeaders := make(map[string][]string, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	return lowerHeaders, nil
}

// NewNoopAudit should be used to create a NoopAudit as it handles creation of a
// predictable salt and wraps eventlogger nodes so information can be retrieved on
// what they've seen or formatted.
//
// Deprecated: NewNoopAudit only exists to allow legacy tests to continue working.
func NewNoopAudit(config *BackendConfig) (*NoopAudit, error) {
	view := &logical.InmemStorage{}

	// Create the salt with a known key for predictable hmac values.
	se := &logical.StorageEntry{Key: "salt", Value: []byte("foo")}
	err := view.Put(context.Background(), se)
	if err != nil {
		return nil, err
	}

	// Override the salt related config settings.
	backendConfig := &BackendConfig{
		SaltView: view,
		SaltConfig: &salt.Config{
			HMAC:     sha256.New,
			HMACType: "hmac-sha256",
		},
		Config:    config.Config,
		MountPath: config.MountPath,
	}

	noopBackend := &NoopAudit{
		Config:     backendConfig,
		nodeIDList: make([]eventlogger.NodeID, 2),
		nodeMap:    make(map[eventlogger.NodeID]eventlogger.Node, 2),
	}

	cfg, err := newFormatterConfig(&noopHeaderFormatter{}, nil)
	if err != nil {
		return nil, err
	}

	formatterNodeID, err := event.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating random NodeID for formatter node: %w", err)
	}

	formatterNode, err := newEntryFormatter(config.MountPath, cfg, noopBackend, config.Logger)
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}

	// Wrap the formatting node, so we can get any bytes that were formatted etc.
	wrappedFormatter := &noopWrapper{format: "json", node: formatterNode, backend: noopBackend}

	noopBackend.nodeIDList[0] = formatterNodeID
	noopBackend.nodeMap[formatterNodeID] = wrappedFormatter

	sinkNode := event.NewNoopSink()
	sinkNodeID, err := event.GenerateNodeID()
	if err != nil {
		return nil, fmt.Errorf("error generating random NodeID for sink node: %w", err)
	}

	noopBackend.nodeIDList[1] = sinkNodeID
	noopBackend.nodeMap[sinkNodeID] = sinkNode

	return noopBackend, nil
}

// NoopAuditFactory should be used when the test needs a way to access bytes that
// have been formatted by the pipeline during audit requests.
// The records parameter will be repointed to the one used within the pipeline.
//
// Deprecated: NoopAuditFactory only exists to allow legacy tests to continue working.
func NoopAuditFactory(records **[][]byte) Factory {
	return func(config *BackendConfig, _ HeaderFormatter) (Backend, error) {
		n, err := NewNoopAudit(config)
		if err != nil {
			return nil, err
		}
		if records != nil {
			*records = &n.records
		}

		return n, nil
	}
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
	a, ok := e.Payload.(*Event)
	if !ok {
		return nil, errors.New("cannot parse payload as an audit event")
	}

	if n.backend.listener != nil {
		n.backend.listener(a)
	}

	in := a.Data

	// Depending on the type of the audit event (request or response) we need to
	// track different things.
	switch a.Subtype {
	case RequestType:
		n.backend.ReqAuth = append(n.backend.ReqAuth, in.Auth)
		n.backend.Req = append(n.backend.Req, in.Request)
		n.backend.ReqNonHMACKeys = in.NonHMACReqDataKeys
		n.backend.ReqErrs = append(n.backend.ReqErrs, in.OuterErr)

		if n.backend.ReqErr != nil {
			return nil, n.backend.ReqErr
		}
	case ResponseType:
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
	if a.Subtype == RequestType {
		reqEntry := &entry{}
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

// LogTestMessage will manually crank the handle on the nodes associated with this backend.
func (n *NoopAudit) LogTestMessage(ctx context.Context, in *logical.LogInput) error {
	if len(n.nodeIDList) > 0 {
		return processManual(ctx, in, n.nodeIDList, n.nodeMap)
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

func (n *NoopAudit) Reload() error {
	return nil
}

func (n *NoopAudit) Invalidate(_ context.Context) {
	n.saltMutex.Lock()
	defer n.saltMutex.Unlock()
	n.salt = nil
}

func (n *NoopAudit) EventType() eventlogger.EventType {
	return event.AuditType.AsEventType()
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

// Deprecated: TestNoopAudit only exists to allow legacy tests to continue working.
func TestNoopAudit(t *testing.T, path string, config map[string]string) *NoopAudit {
	cfg := &BackendConfig{
		Config:    config,
		MountPath: path,
		Logger:    corehelpers.NewTestLogger(t),
	}
	n, err := NewNoopAudit(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return n
}
