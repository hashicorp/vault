package vault

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"errors"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
)

type NoopAudit struct {
	Config  *audit.BackendConfig
	ReqErr  error
	ReqAuth []*logical.Auth
	Req     []*logical.Request
	ReqErrs []error

	RespErr  error
	RespAuth []*logical.Auth
	RespReq  []*logical.Request
	Resp     []*logical.Response
	RespErrs []error
}

func (n *NoopAudit) LogRequest(a *logical.Auth, r *logical.Request, err error) error {
	n.ReqAuth = append(n.ReqAuth, a)
	n.Req = append(n.Req, r)
	n.ReqErrs = append(n.ReqErrs, err)
	return n.ReqErr
}

func (n *NoopAudit) LogResponse(a *logical.Auth, r *logical.Request, re *logical.Response, err error) error {
	n.RespAuth = append(n.RespAuth, a)
	n.RespReq = append(n.RespReq, r)
	n.Resp = append(n.Resp, re)
	n.RespErrs = append(n.RespErrs, err)
	return n.RespErr
}

func (n *NoopAudit) GetHash(data string) string {
	return n.Config.Salt.GetIdentifiedHMAC(data)
}

func TestCore_EnableAudit(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	me := &MountEntry{
		Path: "foo",
		Type: "noop",
	}
	err := c.enableAudit(me)
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
	conf.AuditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
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

func TestCore_DisableAudit(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	err := c.disableAudit("foo")
	if err.Error() != "no matching backend" {
		t.Fatalf("err: %v", err)
	}

	me := &MountEntry{
		Path: "foo",
		Type: "noop",
	}
	err = c.enableAudit(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = c.disableAudit("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
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
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.audit, c2.audit) {
		t.Fatalf("mismatch: %v %v", c.audit, c2.audit)
	}
}

func TestCore_DefaultAuditTable(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
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
	unseal, err := c2.Unseal(key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
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
}

func TestAuditBroker_LogRequest(t *testing.T) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	b := NewAuditBroker(l)
	a1 := &NoopAudit{}
	a2 := &NoopAudit{}
	b.Register("foo", a1, nil)
	b.Register("bar", a2, nil)

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
	reqErrs := errors.New("errs")

	err := b.LogRequest(auth, req, reqErrs)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, a := range []*NoopAudit{a1, a2} {
		if !reflect.DeepEqual(a.ReqAuth[0], auth) {
			t.Fatalf("Bad: %#v", a.ReqAuth[0])
		}
		if !reflect.DeepEqual(a.Req[0], req) {
			t.Fatalf("Bad: %#v", a.Req[0])
		}
		if !reflect.DeepEqual(a.ReqErrs[0], reqErrs) {
			t.Fatalf("Bad: %#v", a.ReqErrs[0])
		}
	}

	// Should still work with one failing backend
	a1.ReqErr = fmt.Errorf("failed")
	if err := b.LogRequest(auth, req, nil); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL work with both failing backends
	a2.ReqErr = fmt.Errorf("failed")
	if err := b.LogRequest(auth, req, nil); err.Error() != "no audit backend succeeded in logging the request" {
		t.Fatalf("err: %v", err)
	}
}

func TestAuditBroker_LogResponse(t *testing.T) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	b := NewAuditBroker(l)
	a1 := &NoopAudit{}
	a2 := &NoopAudit{}
	b.Register("foo", a1, nil)
	b.Register("bar", a2, nil)

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

	err := b.LogResponse(auth, req, resp, respErr)
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
			t.Fatalf("Bad: %#v", a.RespErrs[0])
		}
	}

	// Should still work with one failing backend
	a1.RespErr = fmt.Errorf("failed")
	err = b.LogResponse(auth, req, resp, respErr)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should FAIL work with both failing backends
	a2.RespErr = fmt.Errorf("failed")
	err = b.LogResponse(auth, req, resp, respErr)
	if err.Error() != "no audit backend succeeded in logging the response" {
		t.Fatalf("err: %v", err)
	}
}
