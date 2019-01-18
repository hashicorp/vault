package vault

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"errors"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
)

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
	RespNonHMACKeys    []string
	RespReqNonHMACKeys []string
	RespErrs           []error

	salt      *salt.Salt
	saltMutex sync.RWMutex
}

func (n *NoopAudit) LogRequest(ctx context.Context, in *audit.LogInput) error {
	n.ReqAuth = append(n.ReqAuth, in.Auth)
	n.Req = append(n.Req, in.Request)
	n.ReqHeaders = append(n.ReqHeaders, in.Request.Headers)
	n.ReqNonHMACKeys = in.NonHMACReqDataKeys
	n.ReqErrs = append(n.ReqErrs, in.OuterErr)
	return n.ReqErr
}

func (n *NoopAudit) LogResponse(ctx context.Context, in *audit.LogInput) error {
	n.RespAuth = append(n.RespAuth, in.Auth)
	n.RespReq = append(n.RespReq, in.Request)
	n.Resp = append(n.Resp, in.Response)
	n.RespErrs = append(n.RespErrs, in.OuterErr)

	if in.Response != nil {
		n.RespNonHMACKeys = in.NonHMACRespDataKeys
		n.RespReqNonHMACKeys = in.NonHMACReqDataKeys
	}

	return n.RespErr
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

func TestAudit_ReadOnlyViewDuringMount(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		err := config.SaltView.Put(ctx, &logical.StorageEntry{
			Key:   "bar",
			Value: []byte("baz"),
		})
		if err == nil || !strings.Contains(err.Error(), logical.ErrSetupReadOnly.Error()) {
			t.Fatalf("expected a read-only error")
		}
		return &NoopAudit{}, nil
	}

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableAudit(namespace.RootContext(nil), me, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_EnableAudit(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableAudit(namespace.RootContext(nil), me, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !c.auditBroker.IsRegistered("foo/") {
		t.Fatalf("missing audit backend")
	}

	conf := &CoreConfig{
		Physical:      c.physical,
		AuditBackends: make(map[string]audit.Factory),
		DisableMlock:  true,
	}
	conf.AuditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching audit tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}

	// Check for registration
	if !c2.auditBroker.IsRegistered("foo/") {
		t.Fatalf("missing audit backend")
	}
}

func TestCore_EnableAudit_MixedFailures(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	c.auditBackends["fail"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return nil, fmt.Errorf("failing enabling")
	}

	c.audit = &MountTable{
		Type: auditTableType,
		Entries: []*MountEntry{
			&MountEntry{
				Table: auditTableType,
				Path:  "noop/",
				Type:  "noop",
				UUID:  "abcd",
			},
			&MountEntry{
				Table: auditTableType,
				Path:  "noop2/",
				Type:  "noop",
				UUID:  "bcde",
			},
		},
	}

	// Both should set up successfully
	err := c.setupAudits(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// We expect this to work because the other entry is still valid
	c.audit.Entries[0].Type = "fail"
	err = c.setupAudits(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// No audit backend set up successfully, so expect error
	c.audit.Entries[1].Type = "fail"
	err = c.setupAudits(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

// Test that the local table actually gets populated as expected with local
// entries, and that upon reading the entries from both are recombined
// correctly
func TestCore_EnableAudit_Local(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	c.auditBackends["fail"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return nil, fmt.Errorf("failing enabling")
	}

	c.audit = &MountTable{
		Type: auditTableType,
		Entries: []*MountEntry{
			&MountEntry{
				Table:       auditTableType,
				Path:        "noop/",
				Type:        "noop",
				UUID:        "abcd",
				Accessor:    "noop-abcd",
				NamespaceID: namespace.RootNamespaceID,
				namespace:   namespace.RootNamespace,
			},
			&MountEntry{
				Table:       auditTableType,
				Path:        "noop2/",
				Type:        "noop",
				UUID:        "bcde",
				Accessor:    "noop-bcde",
				NamespaceID: namespace.RootNamespaceID,
				namespace:   namespace.RootNamespace,
			},
		},
	}

	// Both should set up successfully
	err := c.setupAudits(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	rawLocal, err := c.barrier.Get(context.Background(), coreLocalAuditConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local audit")
	}
	localAuditTable := &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localAuditTable); err != nil {
		t.Fatal(err)
	}
	if len(localAuditTable.Entries) > 0 {
		t.Fatalf("expected no entries in local audit table, got %#v", localAuditTable)
	}

	c.audit.Entries[1].Local = true
	if err := c.persistAudit(context.Background(), c.audit, false); err != nil {
		t.Fatal(err)
	}

	rawLocal, err = c.barrier.Get(context.Background(), coreLocalAuditConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local audit")
	}
	localAuditTable = &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localAuditTable); err != nil {
		t.Fatal(err)
	}
	if len(localAuditTable.Entries) != 1 {
		t.Fatalf("expected one entry in local audit table, got %#v", localAuditTable)
	}

	oldAudit := c.audit
	if err := c.loadAudits(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(oldAudit, c.audit) {
		t.Fatalf("expected\n%#v\ngot\n%#v\n", oldAudit, c.audit)
	}

	if len(c.audit.Entries) != 2 {
		t.Fatalf("expected two audit entries, got %#v", localAuditTable)
	}
}

func TestCore_DisableAudit(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	existed, err := c.disableAudit(namespace.RootContext(nil), "foo", true)
	if existed && err != nil {
		t.Fatalf("existed: %v; err: %v", existed, err)
	}

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err = c.enableAudit(namespace.RootContext(nil), me, true)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	existed, err = c.disableAudit(namespace.RootContext(nil), "foo", true)
	if !existed || err != nil {
		t.Fatalf("existed: %v; err: %v", existed, err)
	}

	// Check for registration
	if c.auditBroker.IsRegistered("foo") {
		t.Fatalf("audit backend present")
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch:\n%#v\n%#v", c.audit, c2.audit)
	}
}

func TestCore_DefaultAuditTable(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	verifyDefaultAuditTable(t, c.audit)

	// Verify we have an audit broker
	if c.auditBroker == nil {
		t.Fatalf("missing audit broker")
	}

	// Start a second core with same physical
	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c2, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("should be unsealed")
		}
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}
}

func TestDefaultAuditTable(t *testing.T) {
	table := defaultAuditTable()
	verifyDefaultAuditTable(t, table)
}

func verifyDefaultAuditTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 0 {
		t.Fatalf("bad: %v", table.Entries)
	}
	if table.Type != auditTableType {
		t.Fatalf("bad: %v", *table)
	}
}

func TestAuditBroker_LogRequest(t *testing.T) {
	l := logging.NewVaultLogger(log.Trace)
	b := NewAuditBroker(l)
	a1 := &NoopAudit{}
	a2 := &NoopAudit{}
	b.Register("foo", a1, nil, false)
	b.Register("bar", a2, nil, false)

	auth := &logical.Auth{
		ClientToken: "foo",
		Policies:    []string{"dev", "ops"},
		Metadata: map[string]string{
			"user":   "armon",
			"source": "github",
		},
	}
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}

	// Copy so we can verify nothing changed
	authCopyRaw, err := copystructure.Copy(auth)
	if err != nil {
		t.Fatal(err)
	}
	authCopy := authCopyRaw.(*logical.Auth)

	reqCopyRaw, err := copystructure.Copy(req)
	if err != nil {
		t.Fatal(err)
	}
	reqCopy := reqCopyRaw.(*logical.Request)

	// Create an identifier for the request to verify against
	req.ID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate identifier for the request: path%s err: %v", req.Path, err)
	}
	reqCopy.ID = req.ID

	reqErrs := errors.New("errs")

	headersConf := &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
	}

	logInput := &audit.LogInput{
		Auth:     authCopy,
		Request:  reqCopy,
		OuterErr: reqErrs,
	}
	err = b.LogRequest(context.Background(), logInput, headersConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, a := range []*NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.ReqAuth[0], auth) {
			t.Fatalf("Bad: %#v", a.ReqAuth[0])
		}
		if !reflect.DeepEqual(a.Req[0], req) {
			t.Fatalf("Bad: %#v\n wanted %#v", a.Req[0], req)
		}
		if !reflect.DeepEqual(a.ReqErrs[0], reqErrs) {
			t.Fatalf("Bad: %#v", a.ReqErrs[0])
		}
	}

	// Should still work with one failing backend
	a1.ReqErr = fmt.Errorf("failed")
	logInput = &audit.LogInput{
		Auth:    auth,
		Request: req,
	}
	if err := b.LogRequest(context.Background(), logInput, headersConf); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL work with both failing backends
	a2.ReqErr = fmt.Errorf("failed")
	if err := b.LogRequest(context.Background(), logInput, headersConf); !errwrap.Contains(err, "no audit backend succeeded in logging the request") {
		t.Fatalf("err: %v", err)
	}
}

func TestAuditBroker_LogResponse(t *testing.T) {
	l := logging.NewVaultLogger(log.Trace)
	b := NewAuditBroker(l)
	a1 := &NoopAudit{}
	a2 := &NoopAudit{}
	b.Register("foo", a1, nil, false)
	b.Register("bar", a2, nil, false)

	auth := &logical.Auth{
		NumUses:     10,
		ClientToken: "foo",
		Policies:    []string{"dev", "ops"},
		Metadata: map[string]string{
			"user":   "armon",
			"source": "github",
		},
	}
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 1 * time.Hour,
			},
		},
		Data: map[string]interface{}{
			"user":     "root",
			"password": "password",
		},
	}
	respErr := fmt.Errorf("permission denied")

	// Copy so we can verify nothing changed
	authCopyRaw, err := copystructure.Copy(auth)
	if err != nil {
		t.Fatal(err)
	}
	authCopy := authCopyRaw.(*logical.Auth)

	reqCopyRaw, err := copystructure.Copy(req)
	if err != nil {
		t.Fatal(err)
	}
	reqCopy := reqCopyRaw.(*logical.Request)

	respCopyRaw, err := copystructure.Copy(resp)
	if err != nil {
		t.Fatal(err)
	}
	respCopy := respCopyRaw.(*logical.Response)

	headersConf := &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
	}

	logInput := &audit.LogInput{
		Auth:     authCopy,
		Request:  reqCopy,
		Response: respCopy,
		OuterErr: respErr,
	}
	err = b.LogResponse(context.Background(), logInput, headersConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, a := range []*NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.RespAuth[0], auth) {
			t.Fatalf("Bad: %#v", a.ReqAuth[0])
		}
		if !reflect.DeepEqual(a.RespReq[0], req) {
			t.Fatalf("Bad: %#v", a.Req[0])
		}
		if !reflect.DeepEqual(a.Resp[0], resp) {
			t.Fatalf("Bad: %#v", a.Resp[0])
		}
		if !reflect.DeepEqual(a.RespErrs[0], respErr) {
			t.Fatalf("Expected\n%v\nGot\n%#v", respErr, a.RespErrs[0])
		}
	}

	// Should still work with one failing backend
	a1.RespErr = fmt.Errorf("failed")
	logInput = &audit.LogInput{
		Auth:     auth,
		Request:  req,
		Response: resp,
		OuterErr: respErr,
	}
	err = b.LogResponse(context.Background(), logInput, headersConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL work with both failing backends
	a2.RespErr = fmt.Errorf("failed")
	err = b.LogResponse(context.Background(), logInput, headersConf)
	if !strings.Contains(err.Error(), "no audit backend succeeded in logging the response") {
		t.Fatalf("err: %v", err)
	}
}

func TestAuditBroker_AuditHeaders(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	b := NewAuditBroker(logger)
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "headers/")
	a1 := &NoopAudit{}
	a2 := &NoopAudit{}
	b.Register("foo", a1, nil, false)
	b.Register("bar", a2, nil, false)

	auth := &logical.Auth{
		ClientToken: "foo",
		Policies:    []string{"dev", "ops"},
		Metadata: map[string]string{
			"user":   "armon",
			"source": "github",
		},
	}
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
		Headers: map[string][]string{
			"X-Test-Header":  []string{"foo"},
			"X-Vault-Header": []string{"bar"},
			"Content-Type":   []string{"baz"},
		},
	}
	respErr := fmt.Errorf("permission denied")

	// Copy so we can verify nothing changed
	reqCopyRaw, err := copystructure.Copy(req)
	if err != nil {
		t.Fatal(err)
	}
	reqCopy := reqCopyRaw.(*logical.Request)

	headersConf := &AuditedHeadersConfig{
		view: view,
	}
	headersConf.add(context.Background(), "X-Test-Header", false)
	headersConf.add(context.Background(), "X-Vault-Header", false)

	logInput := &audit.LogInput{
		Auth:     auth,
		Request:  reqCopy,
		OuterErr: respErr,
	}
	err = b.LogRequest(context.Background(), logInput, headersConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := map[string][]string{
		"x-test-header":  []string{"foo"},
		"x-vault-header": []string{"bar"},
	}

	for _, a := range []*NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.ReqHeaders[0], expected) {
			t.Fatalf("Bad audited headers: %#v", a.Req[0].Headers)
		}
	}

	// Should still work with one failing backend
	a1.ReqErr = fmt.Errorf("failed")
	logInput = &audit.LogInput{
		Auth:     auth,
		Request:  req,
		OuterErr: respErr,
	}
	err = b.LogRequest(context.Background(), logInput, headersConf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL work with both failing backends
	a2.ReqErr = fmt.Errorf("failed")
	err = b.LogRequest(context.Background(), logInput, headersConf)
	if !errwrap.Contains(err, "no audit backend succeeded in logging the request") {
		t.Fatalf("err: %v", err)
	}
}
