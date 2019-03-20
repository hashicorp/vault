package vault

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
)

var (
	// invalidKey is used to test Unseal
	invalidKey = []byte("abcdefghijklmnopqrstuvwxyz")[:17]
)

func TestNewCore_badRedirectAddr(t *testing.T) {
	logger = logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	conf := &CoreConfig{
		RedirectAddr: "127.0.0.1:8200",
		Physical:     inm,
		DisableMlock: true,
	}
	_, err = NewCore(conf)
	if err == nil {
		t.Fatal("should error")
	}
}

func TestSealConfig_Invalid(t *testing.T) {
	s := &SealConfig{
		SecretShares:    2,
		SecretThreshold: 1,
	}
	err := s.Validate()
	if err == nil {
		t.Fatalf("expected err")
	}
}

func TestCore_Unseal_MultiShare(t *testing.T) {
	c := TestCore(t)

	_, err := TestCoreUnseal(c, invalidKey)
	if err != ErrNotInit {
		t.Fatalf("err: %v", err)
	}

	sealConf := &SealConfig{
		SecretShares:    5,
		SecretThreshold: 3,
	}
	res, err := c.Initialize(namespace.RootContext(nil), &InitParams{
		BarrierConfig:  sealConf,
		RecoveryConfig: nil,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !c.Sealed() {
		t.Fatalf("should be sealed")
	}

	if prog, _ := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	for i := 0; i < 5; i++ {
		unseal, err := TestCoreUnseal(c, res.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}

		// Ignore redundant
		_, err = TestCoreUnseal(c, res.SecretShares[i])
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i >= 2 {
			if !unseal {
				t.Fatalf("should be unsealed")
			}
			if prog, _ := c.SecretProgress(); prog != 0 {
				t.Fatalf("bad progress: %d", prog)
			}
		} else {
			if unseal {
				t.Fatalf("should not be unsealed")
			}
			if prog, _ := c.SecretProgress(); prog != i+1 {
				t.Fatalf("bad progress: %d", prog)
			}
		}
	}

	if c.Sealed() {
		t.Fatalf("should not be sealed")
	}

	err = c.Seal(res.RootToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ignore redundant
	err = c.Seal(res.RootToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !c.Sealed() {
		t.Fatalf("should be sealed")
	}
}

func TestCore_Unseal_Single(t *testing.T) {
	c := TestCore(t)

	_, err := TestCoreUnseal(c, invalidKey)
	if err != ErrNotInit {
		t.Fatalf("err: %v", err)
	}

	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}
	res, err := c.Initialize(namespace.RootContext(nil), &InitParams{
		BarrierConfig:  sealConf,
		RecoveryConfig: nil,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !c.Sealed() {
		t.Fatalf("should be sealed")
	}

	if prog, _ := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	unseal, err := TestCoreUnseal(c, res.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !unseal {
		t.Fatalf("should be unsealed")
	}
	if prog, _ := c.SecretProgress(); prog != 0 {
		t.Fatalf("bad progress: %d", prog)
	}

	if c.Sealed() {
		t.Fatalf("should not be sealed")
	}
}

func TestCore_Route_Sealed(t *testing.T) {
	c := TestCore(t)
	sealConf := &SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	}

	ctx := namespace.RootContext(nil)

	// Should not route anything
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "sys/mounts",
	}
	_, err := c.HandleRequest(ctx, req)
	if err != consts.ErrSealed {
		t.Fatalf("err: %v", err)
	}

	res, err := c.Initialize(ctx, &InitParams{
		BarrierConfig:  sealConf,
		RecoveryConfig: nil,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	unseal, err := TestCoreUnseal(c, res.SecretShares[0])
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Should not error after unseal
	req.ClientToken = res.RootToken
	_, err = c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

// Attempt to unseal after doing a first seal
func TestCore_SealUnseal(t *testing.T) {
	c, keys, root := TestCoreUnsealed(t)
	if err := c.Seal(root); err != nil {
		t.Fatalf("err: %v", err)
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("err: should be unsealed")
		}
	}
}

// Attempt to shutdown after unseal
func TestCore_Shutdown(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	if err := c.Shutdown(); err != nil {
		t.Fatalf("err: %v", err)
	}
	if !c.Sealed() {
		t.Fatal("wasn't sealed")
	}
}

// Attempt to seal bad token
func TestCore_Seal_BadToken(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	if err := c.Seal("foo"); err == nil {
		t.Fatalf("err: %v", err)
	}
	if c.Sealed() {
		t.Fatal("was sealed")
	}
}

// GH-3497
func TestCore_Seal_SingleUse(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	c.tokenStore.create(namespace.RootContext(nil), &logical.TokenEntry{
		ID:          "foo",
		NumUses:     1,
		Policies:    []string{"root"},
		NamespaceID: namespace.RootNamespaceID,
	})
	if err := c.Seal("foo"); err != nil {
		t.Fatalf("err: %v", err)
	}
	if !c.Sealed() {
		t.Fatal("not sealed")
	}
	for i, key := range keys {
		unseal, err := TestCoreUnseal(c, key)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if i+1 == len(keys) && !unseal {
			t.Fatalf("err: should be unsealed")
		}
	}
	if err := c.Seal("foo"); err == nil {
		t.Fatal("expected error from revoked token")
	}
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), "foo")
	if err != nil {
		t.Fatal(err)
	}
	if te != nil {
		t.Fatalf("expected nil token entry, got %#v", *te)
	}
}

// Ensure we get a LeaseID
func TestCore_HandleRequest_Lease(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: root,
	}
	ctx := namespace.RootContext(nil)
	resp, err := c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the key
	req.Operation = logical.ReadOperation
	req.Data = nil
	req.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	resp, err = c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Secret.TTL != time.Hour {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	if resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp.Data)
	}
}

func TestCore_HandleRequest_Lease_MaxLength(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1000h",
		},
		ClientToken: root,
	}
	ctx := namespace.RootContext(nil)
	resp, err := c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the key
	req.Operation = logical.ReadOperation
	req.Data = nil
	req.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	resp, err = c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Secret.TTL != c.maxLeaseTTL {
		t.Fatalf("bad: %#v, %d", resp.Secret, c.maxLeaseTTL)
	}
	if resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp.Data)
	}
}

func TestCore_HandleRequest_Lease_DefaultLength(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "0h",
		},
		ClientToken: root,
	}
	ctx := namespace.RootContext(nil)
	resp, err := c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the key
	req.Operation = logical.ReadOperation
	req.Data = nil
	req.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	resp, err = c.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	if resp.Secret.TTL != c.defaultLeaseTTL {
		t.Fatalf("bad: %d, %d", resp.Secret.TTL/time.Second, c.defaultLeaseTTL/time.Second)
	}
	if resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp.Data)
	}
}

func TestCore_HandleRequest_MissingToken(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || !errwrap.Contains(err, logical.ErrInvalidRequest.Error()) {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "missing client token" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestCore_HandleRequest_InvalidToken(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: "foobarbaz",
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || !errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "permission denied" {
		t.Fatalf("bad: %#v", resp)
	}
}

// Check that standard permissions work
func TestCore_HandleRequest_NoSlash(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	req := &logical.Request{
		Operation:   logical.HelpOperation,
		Path:        "secret",
		ClientToken: root,
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
	if _, ok := resp.Data["help"]; !ok {
		t.Fatalf("resp: %v", resp)
	}
}

// Test a root path is denied if non-root
func TestCore_HandleRequest_RootPath(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	testMakeServiceTokenViaCore(t, c, root, "child", "", []string{"test"})

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "sys/policy", // root protected!
		ClientToken: "child",
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || !errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
}

// Test a root path is allowed if non-root but with sudo
func TestCore_HandleRequest_RootPath_WithSudo(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	// Set the 'test' policy object to permit access to sys/policy
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test", // root protected!
		Data: map[string]interface{}{
			"rules": `path "sys/policy" { policy = "sudo" }`,
		},
		ClientToken: root,
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil && (resp.IsError() || len(resp.Data) > 0) {
		t.Fatalf("bad: %#v", resp)
	}

	// Child token (non-root) but with 'test' policy should have access
	testMakeServiceTokenViaCore(t, c, root, "child", "", []string{"test"})
	req = &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "sys/policy", // root protected!
		ClientToken: "child",
	}
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %#v", resp)
	}
}

// Check that standard permissions work
func TestCore_HandleRequest_PermissionDenied(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	testMakeServiceTokenViaCore(t, c, root, "child", "", []string{"test"})

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: "child",
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || !errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
		t.Fatalf("err: %v, resp: %v", err, resp)
	}
}

// Check that standard permissions work
func TestCore_HandleRequest_PermissionAllowed(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)
	testMakeServiceTokenViaCore(t, c, root, "child", "", []string{"test"})

	// Set the 'test' policy object to permit access to secret/
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "secret/*" { policy = "write" }`,
		},
		ClientToken: root,
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil && (resp.IsError() || len(resp.Data) > 0) {
		t.Fatalf("bad: %#v", resp)
	}

	// Write should work now
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: "child",
	}
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestCore_HandleRequest_NoClientToken(t *testing.T) {
	noop := &NoopBackend{
		Response: &logical.Response{},
	}
	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the logical backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo")
	req.Data["type"] = "noop"
	req.Data["description"] = "foo"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to request with connection data
	req = &logical.Request{
		Path: "foo/login",
	}
	req.ClientToken = root
	if _, err := c.HandleRequest(namespace.RootContext(nil), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	ct := noop.Requests[0].ClientToken
	if ct == "" || ct == root {
		t.Fatalf("bad: %#v", noop.Requests)
	}
}

func TestCore_HandleRequest_ConnOnLogin(t *testing.T) {
	noop := &NoopBackend{
		Login:       []string{"login"},
		Response:    &logical.Response{},
		BackendType: logical.TypeCredential,
	}
	c, _, root := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to request with connection data
	req = &logical.Request{
		Path:       "auth/foo/login",
		Connection: &logical.Connection{},
	}
	if _, err := c.HandleRequest(namespace.RootContext(nil), req); err != nil {
		t.Fatalf("err: %v", err)
	}
	if noop.Requests[0].Connection == nil {
		t.Fatalf("bad: %#v", noop.Requests)
	}
}

// Ensure we get a client token
func TestCore_HandleLogin_Token(t *testing.T) {
	noop := &NoopBackend{
		Login: []string{"login"},
		Response: &logical.Response{
			Auth: &logical.Auth{
				Policies: []string{"foo", "bar"},
				Metadata: map[string]string{
					"user": "armon",
				},
				DisplayName: "armon",
			},
		},
		BackendType: logical.TypeCredential,
	}
	c, _, root := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to login
	lreq := &logical.Request{
		Path: "auth/foo/login",
	}
	lresp, err := c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we got a client token back
	clientToken := lresp.Auth.ClientToken
	if clientToken == "" {
		t.Fatalf("bad: %#v", lresp)
	}

	// Check the policy and metadata
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), clientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expect := &logical.TokenEntry{
		ID:       clientToken,
		Accessor: te.Accessor,
		Parent:   "",
		Policies: []string{"bar", "default", "foo"},
		Path:     "auth/foo/login",
		Meta: map[string]string{
			"user": "armon",
		},
		DisplayName:  "foo-armon",
		TTL:          time.Hour * 24,
		CreationTime: te.CreationTime,
		NamespaceID:  namespace.RootNamespaceID,
		CubbyholeID:  te.CubbyholeID,
		Type:         logical.TokenTypeService,
	}

	if diff := deep.Equal(te, expect); diff != nil {
		t.Fatal(diff)
	}

	// Check that we have a lease with default duration
	if lresp.Auth.TTL != noop.System().DefaultLeaseTTL() {
		t.Fatalf("bad: %#v, defaultLeaseTTL: %#v", lresp.Auth, c.defaultLeaseTTL)
	}
}

func TestCore_HandleRequest_AuditTrail(t *testing.T) {
	// Create a noop audit backend
	noop := &NoopAudit{}
	c, _, root := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		noop = &NoopAudit{
			Config: config,
		}
		return noop, nil
	}

	// Enable the audit backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/audit/noop")
	req.Data["type"] = "noop"
	req.ClientToken = root
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make a request
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: root,
	}
	req.ClientToken = root
	if _, err := c.HandleRequest(namespace.RootContext(nil), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the audit trail on request and response
	if len(noop.ReqAuth) != 1 {
		t.Fatalf("bad: %#v", noop)
	}
	auth := noop.ReqAuth[0]
	if auth.ClientToken != root {
		t.Fatalf("bad client token: %#v", auth)
	}
	if len(auth.Policies) != 1 || auth.Policies[0] != "root" {
		t.Fatalf("bad: %#v", auth)
	}
	if len(noop.Req) != 1 || !reflect.DeepEqual(noop.Req[0], req) {
		t.Fatalf("Bad: %#v", noop.Req[0])
	}

	if len(noop.RespAuth) != 2 {
		t.Fatalf("bad: %#v", noop)
	}
	if !reflect.DeepEqual(noop.RespAuth[1], auth) {
		t.Fatalf("bad: %#v", auth)
	}
	if len(noop.RespReq) != 2 || !reflect.DeepEqual(noop.RespReq[1], req) {
		t.Fatalf("Bad: %#v", noop.RespReq[1])
	}
	if len(noop.Resp) != 2 || !reflect.DeepEqual(noop.Resp[1], resp) {
		t.Fatalf("Bad: %#v", noop.Resp[1])
	}
}

func TestCore_HandleRequest_AuditTrail_noHMACKeys(t *testing.T) {
	// Create a noop audit backend
	var noop *NoopAudit
	c, _, root := TestCoreUnsealed(t)
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		noop = &NoopAudit{
			Config: config,
		}
		return noop, nil
	}

	// Specify some keys to not HMAC
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/secret/tune")
	req.Data["audit_non_hmac_request_keys"] = "foo"
	req.ClientToken = root
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/secret/tune")
	req.Data["audit_non_hmac_response_keys"] = "baz"
	req.ClientToken = root
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Enable the audit backend
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/audit/noop")
	req.Data["type"] = "noop"
	req.ClientToken = root
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Make a request
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo": "bar",
		},
		ClientToken: root,
	}
	req.ClientToken = root
	if _, err := c.HandleRequest(namespace.RootContext(nil), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the audit trail on request and response
	if len(noop.ReqAuth) != 1 {
		t.Fatalf("bad: %#v", noop)
	}
	auth := noop.ReqAuth[0]
	if auth.ClientToken != root {
		t.Fatalf("bad client token: %#v", auth)
	}
	if len(auth.Policies) != 1 || auth.Policies[0] != "root" {
		t.Fatalf("bad: %#v", auth)
	}
	if len(noop.Req) != 1 || !reflect.DeepEqual(noop.Req[0], req) {
		t.Fatalf("Bad: %#v", noop.Req[0])
	}
	if len(noop.ReqNonHMACKeys) != 1 || noop.ReqNonHMACKeys[0] != "foo" {
		t.Fatalf("Bad: %#v", noop.ReqNonHMACKeys)
	}
	if len(noop.RespAuth) != 2 {
		t.Fatalf("bad: %#v", noop)
	}
	if !reflect.DeepEqual(noop.RespAuth[1], auth) {
		t.Fatalf("bad: %#v", auth)
	}
	if len(noop.RespReq) != 2 || !reflect.DeepEqual(noop.RespReq[1], req) {
		t.Fatalf("Bad: %#v", noop.RespReq[1])
	}
	if len(noop.Resp) != 2 || !reflect.DeepEqual(noop.Resp[1], resp) {
		t.Fatalf("Bad: %#v", noop.Resp[1])
	}

	// Test for response keys
	// Make a request
	req = &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "secret/test",
		ClientToken: root,
	}
	req.ClientToken = root
	req.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	if _, err := c.HandleRequest(namespace.RootContext(nil), req); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(noop.RespNonHMACKeys) != 1 || noop.RespNonHMACKeys[0] != "baz" {
		t.Fatalf("Bad: %#v", noop.RespNonHMACKeys)
	}
	if len(noop.RespReqNonHMACKeys) != 1 || noop.RespReqNonHMACKeys[0] != "foo" {
		t.Fatalf("Bad: %#v", noop.RespReqNonHMACKeys)
	}
}

// Ensure we get a client token
func TestCore_HandleLogin_AuditTrail(t *testing.T) {
	// Create a badass credential backend that always logs in as armon
	noop := &NoopAudit{}
	noopBack := &NoopBackend{
		Login: []string{"login"},
		Response: &logical.Response{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL: time.Hour,
				},
				Policies: []string{"foo", "bar"},
				Metadata: map[string]string{
					"user": "armon",
				},
			},
		},
		BackendType: logical.TypeCredential,
	}
	c, _, root := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noopBack, nil
	}
	c.auditBackends["noop"] = func(ctx context.Context, config *audit.BackendConfig) (audit.Backend, error) {
		noop = &NoopAudit{
			Config: config,
		}
		return noop, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Enable the audit backend
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/audit/noop")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to login
	lreq := &logical.Request{
		Path: "auth/foo/login",
	}
	lresp, err := c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we got a client token back
	clientToken := lresp.Auth.ClientToken
	if clientToken == "" {
		t.Fatalf("bad: %#v", lresp)
	}

	// Check the audit trail on request and response
	if len(noop.ReqAuth) != 1 {
		t.Fatalf("bad: %#v", noop)
	}
	if len(noop.Req) != 1 || !reflect.DeepEqual(noop.Req[0], lreq) {
		t.Fatalf("Bad: %#v %#v", noop.Req[0], lreq)
	}

	if len(noop.RespAuth) != 2 {
		t.Fatalf("bad: %#v", noop)
	}
	auth := noop.RespAuth[1]
	if auth.ClientToken != clientToken {
		t.Fatalf("bad client token: %#v", auth)
	}
	if len(auth.Policies) != 3 || auth.Policies[0] != "bar" || auth.Policies[1] != "default" || auth.Policies[2] != "foo" {
		t.Fatalf("bad: %#v", auth)
	}
	if len(noop.RespReq) != 2 || !reflect.DeepEqual(noop.RespReq[1], lreq) {
		t.Fatalf("Bad: %#v", noop.RespReq[1])
	}
	if len(noop.Resp) != 2 || !reflect.DeepEqual(noop.Resp[1], lresp) {
		t.Fatalf("Bad: %#v %#v", noop.Resp[1], lresp)
	}
}

// Check that we register a lease for new tokens
func TestCore_HandleRequest_CreateToken_Lease(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	// Create a new credential
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we got a new client token back
	if resp.IsError() {
		t.Fatalf("err: %v %v", err, *resp)
	}
	clientToken := resp.Auth.ClientToken
	if clientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Check the policy and metadata
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), clientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expect := &logical.TokenEntry{
		ID:           clientToken,
		Accessor:     te.Accessor,
		Parent:       root,
		Policies:     []string{"default", "foo"},
		Path:         "auth/token/create",
		DisplayName:  "token",
		CreationTime: te.CreationTime,
		TTL:          time.Hour * 24 * 32,
		NamespaceID:  namespace.RootNamespaceID,
		CubbyholeID:  te.CubbyholeID,
		Type:         logical.TokenTypeService,
	}
	if diff := deep.Equal(te, expect); diff != nil {
		t.Fatal(diff)
	}

	// Check that we have a lease with default duration
	if resp.Auth.TTL != c.defaultLeaseTTL {
		t.Fatalf("bad: %#v", resp.Auth)
	}
}

// Check that we handle excluding the default policy
func TestCore_HandleRequest_CreateToken_NoDefaultPolicy(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	// Create a new credential
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/create")
	req.ClientToken = root
	req.Data["policies"] = []string{"foo"}
	req.Data["no_default_policy"] = true
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we got a new client token back
	clientToken := resp.Auth.ClientToken
	if clientToken == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Check the policy and metadata
	te, err := c.tokenStore.Lookup(namespace.RootContext(nil), clientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	expect := &logical.TokenEntry{
		ID:           clientToken,
		Accessor:     te.Accessor,
		Parent:       root,
		Policies:     []string{"foo"},
		Path:         "auth/token/create",
		DisplayName:  "token",
		CreationTime: te.CreationTime,
		TTL:          time.Hour * 24 * 32,
		NamespaceID:  namespace.RootNamespaceID,
		CubbyholeID:  te.CubbyholeID,
		Type:         logical.TokenTypeService,
	}
	if diff := deep.Equal(te, expect); diff != nil {
		t.Fatal(diff)
	}
}

func TestCore_LimitedUseToken(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	// Create a new credential
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/token/create")
	req.ClientToken = root
	req.Data["num_uses"] = "1"
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Put a secret
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/foo",
		Data: map[string]interface{}{
			"foo": "bar",
		},
		ClientToken: resp.Auth.ClientToken,
	}
	_, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Second operation should fail
	_, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err == nil || !errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
		t.Fatalf("err: %v", err)
	}
}

func TestCore_Standby_Seal(t *testing.T) {
	// Create the first core and initialize it
	logger = logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for core to become active
	TestWaitActive(t, core)

	// Check the leader is local
	isLeader, advertise, _, err := core.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Create the second core and initialize it
	redirectOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, key := range keys {
		if _, err := TestCoreUnseal(core2, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core2.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Core2 should be in standby
	standby, err := core2.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Check the leader is not local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if isLeader {
		t.Fatalf("should not be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Seal the standby core with the correct token. Shouldn't go down
	err = core2.Seal(root)
	if err == nil {
		t.Fatal("should not be sealed")
	}

	keyUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	// Seal the standby core with an invalid token. Shouldn't go down
	err = core2.Seal(keyUUID)
	if err == nil {
		t.Fatal("should not be sealed")
	}
}

func TestCore_StepDown(t *testing.T) {
	// Create the first core and initialize it
	logger = logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for core to become active
	TestWaitActive(t, core)

	// Check the leader is local
	isLeader, advertise, _, err := core.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Create the second core and initialize it
	redirectOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, key := range keys {
		if _, err := TestCoreUnseal(core2, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core2.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Core2 should be in standby
	standby, err := core2.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Check the leader is not local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if isLeader {
		t.Fatalf("should not be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	req := &logical.Request{
		ClientToken: root,
		Path:        "sys/step-down",
	}

	// Create an identifier for the request
	req.ID, err = uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("failed to generate identifier for the request: path: %s err: %v", req.Path, err)
	}

	// Step down core
	err = core.StepDown(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal("error stepping down core 1")
	}

	// Give time to switch leaders
	time.Sleep(5 * time.Second)

	// Core1 should be in standby
	standby, err = core.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Check the leader is core2
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal2 {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal2)
	}

	// Check the leader is not local
	isLeader, advertise, _, err = core.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if isLeader {
		t.Fatalf("should not be leader")
	}
	if advertise != redirectOriginal2 {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal2)
	}

	// Step down core2
	err = core2.StepDown(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatal("error stepping down core 1")
	}

	// Give time to switch leaders -- core 1 will still be waiting on its
	// cooling off period so give it a full 10 seconds to recover
	time.Sleep(10 * time.Second)

	// Core2 should be in standby
	standby, err = core2.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Check the leader is core1
	isLeader, advertise, _, err = core.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Check the leader is not local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if isLeader {
		t.Fatalf("should not be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}
}

func TestCore_CleanLeaderPrefix(t *testing.T) {
	// Create the first core and initialize it
	logger = logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for core to become active
	TestWaitActive(t, core)

	// Ensure that the original clean function has stopped running
	time.Sleep(2 * time.Second)

	// Put several random entries
	for i := 0; i < 5; i++ {
		keyUUID, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}
		valueUUID, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}
		core.barrier.Put(namespace.RootContext(nil), &logical.StorageEntry{
			Key:   coreLeaderPrefix + keyUUID,
			Value: []byte(valueUUID),
		})
	}

	entries, err := core.barrier.List(namespace.RootContext(nil), coreLeaderPrefix)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(entries) != 6 {
		t.Fatalf("wrong number of core leader prefix entries, got %d", len(entries))
	}

	// Check the leader is local
	isLeader, advertise, _, err := core.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Create a second core, attached to same in-memory store
	redirectOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, key := range keys {
		if _, err := TestCoreUnseal(core2, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core2.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Core2 should be in standby
	standby, err := core2.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Check the leader is not local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if isLeader {
		t.Fatalf("should not be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Seal the first core, should step down
	err = core.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Core should be in standby
	standby, err = core.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Wait for core2 to become active
	TestWaitActive(t, core2)

	// Check the leader is local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal2 {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal2)
	}

	// Give time for the entries to clear out; it is conservative at 1/second
	time.Sleep(10 * leaderPrefixCleanDelay)

	entries, err = core2.barrier.List(namespace.RootContext(nil), coreLeaderPrefix)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("wrong number of core leader prefix entries, got %d", len(entries))
	}
}

func TestCore_Standby(t *testing.T) {
	logger = logging.NewVaultLogger(log.Trace)

	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	testCore_Standby_Common(t, inmha, inmha.(physical.HABackend))
}

func TestCore_Standby_SeparateHA(t *testing.T) {
	logger = logging.NewVaultLogger(log.Trace)

	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha2, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	testCore_Standby_Common(t, inmha, inmha2.(physical.HABackend))
}

func testCore_Standby_Common(t *testing.T, inm physical.Backend, inmha physical.HABackend) {
	// Create the first core and initialize it
	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha,
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Wait for core to become active
	TestWaitActive(t, core)

	testCoreAddSecretMount(t, core, root)

	// Put a secret
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/foo",
		Data: map[string]interface{}{
			"foo": "bar",
		},
		ClientToken: root,
	}
	_, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the leader is local
	isLeader, advertise, _, err := core.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Create a second core, attached to same in-memory store
	redirectOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha,
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, key := range keys {
		if _, err := TestCoreUnseal(core2, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Verify unsealed
	if core2.Sealed() {
		t.Fatal("should not be sealed")
	}

	// Core2 should be in standby
	standby, err := core2.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Request should fail in standby mode
	_, err = core2.HandleRequest(namespace.RootContext(nil), req)
	if err != consts.ErrStandby {
		t.Fatalf("err: %v", err)
	}

	// Check the leader is not local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if isLeader {
		t.Fatalf("should not be leader")
	}
	if advertise != redirectOriginal {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal)
	}

	// Seal the first core, should step down
	err = core.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Core should be in standby
	standby, err = core.Standby()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !standby {
		t.Fatalf("should be standby")
	}

	// Wait for core2 to become active
	TestWaitActive(t, core2)

	// Read the secret
	req = &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "secret/foo",
		ClientToken: root,
	}
	resp, err := core2.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if resp.Data["foo"] != "bar" {
		t.Fatalf("bad: %#v", resp)
	}

	// Check the leader is local
	isLeader, advertise, _, err = core2.Leader()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !isLeader {
		t.Fatalf("should be leader")
	}
	if advertise != redirectOriginal2 {
		t.Fatalf("Bad advertise: %v, orig is %v", advertise, redirectOriginal2)
	}

	if inm.(*inmem.InmemHABackend) == inmha.(*inmem.InmemHABackend) {
		lockSize := inm.(*inmem.InmemHABackend).LockMapSize()
		if lockSize == 0 {
			t.Fatalf("locks not used with only one HA backend")
		}
	} else {
		lockSize := inmha.(*inmem.InmemHABackend).LockMapSize()
		if lockSize == 0 {
			t.Fatalf("locks not used with expected HA backend")
		}

		lockSize = inm.(*inmem.InmemHABackend).LockMapSize()
		if lockSize != 0 {
			t.Fatalf("locks used with unexpected HA backend")
		}
	}
}

// Ensure that InternalData is never returned
func TestCore_HandleRequest_Login_InternalData(t *testing.T) {
	noop := &NoopBackend{
		Login: []string{"login"},
		Response: &logical.Response{
			Auth: &logical.Auth{
				Policies: []string{"foo", "bar"},
				InternalData: map[string]interface{}{
					"foo": "bar",
				},
			},
		},
		BackendType: logical.TypeCredential,
	}

	c, _, root := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to login
	lreq := &logical.Request{
		Path: "auth/foo/login",
	}
	lresp, err := c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we do not get the internal data
	if lresp.Auth.InternalData != nil {
		t.Fatalf("bad: %#v", lresp)
	}
}

// Ensure that InternalData is never returned
func TestCore_HandleRequest_InternalData(t *testing.T) {
	noop := &NoopBackend{
		Response: &logical.Response{
			Secret: &logical.Secret{
				InternalData: map[string]interface{}{
					"foo": "bar",
				},
			},
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	}

	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to read
	lreq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "foo/test",
		ClientToken: root,
	}
	lreq.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	lresp, err := c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure we do not get the internal data
	if lresp.Secret.InternalData != nil {
		t.Fatalf("bad: %#v", lresp)
	}
}

// Ensure login does not return a secret
func TestCore_HandleLogin_ReturnSecret(t *testing.T) {
	// Create a badass credential backend that always logs in as armon
	noopBack := &NoopBackend{
		Login: []string{"login"},
		Response: &logical.Response{
			Secret: &logical.Secret{},
			Auth: &logical.Auth{
				Policies: []string{"foo", "bar"},
			},
		},
		BackendType: logical.TypeCredential,
	}
	c, _, root := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noopBack, nil
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to login
	lreq := &logical.Request{
		Path: "auth/foo/login",
	}
	_, err = c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != ErrInternalError {
		t.Fatalf("err: %v", err)
	}
}

// Renew should return the same lease back
func TestCore_RenewSameLease(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	// Create a leasable secret
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: root,
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the key
	req.Operation = logical.ReadOperation
	req.Data = nil
	req.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp.Secret)
	}
	original := resp.Secret.LeaseID

	// Renew the lease
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/renew/"+resp.Secret.LeaseID)
	req.ClientToken = root
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the lease did not change
	if resp.Secret.LeaseID != original {
		t.Fatalf("lease id changed: %s %s", original, resp.Secret.LeaseID)
	}

	// Renew the lease (alternate path)
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/leases/renew/"+resp.Secret.LeaseID)
	req.ClientToken = root
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the lease did not change
	if resp.Secret.LeaseID != original {
		t.Fatalf("lease id changed: %s %s", original, resp.Secret.LeaseID)
	}
}

// Renew of a token should not create a new lease
func TestCore_RenewToken_SingleRegister(t *testing.T) {
	c, _, root := TestCoreUnsealed(t)

	// Create a new token
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "auth/token/create",
		Data: map[string]interface{}{
			"lease": "1h",
		},
		ClientToken: root,
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	newClient := resp.Auth.ClientToken

	// Renew the token
	req = logical.TestRequest(t, logical.UpdateOperation, "auth/token/renew")
	req.ClientToken = newClient
	req.Data = map[string]interface{}{
		"token": newClient,
	}
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Revoke using the renew prefix
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/revoke-prefix/auth/token/renew/")
	req.ClientToken = root
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify our token is still valid (e.g. we did not get invalidated by the revoke)
	req = logical.TestRequest(t, logical.UpdateOperation, "auth/token/lookup")
	req.Data = map[string]interface{}{
		"token": newClient,
	}
	req.ClientToken = newClient
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the token exists
	if resp.Data["id"] != newClient {
		t.Fatalf("bad: %#v", resp.Data)
	}
}

// Based on bug GH-203, attempt to disable a credential backend with leased secrets
func TestCore_EnableDisableCred_WithLease(t *testing.T) {
	noopBack := &NoopBackend{
		Login: []string{"login"},
		Response: &logical.Response{
			Auth: &logical.Auth{
				Policies: []string{"root"},
			},
		},
		BackendType: logical.TypeCredential,
	}

	c, _, root := TestCoreUnsealed(t)
	c.credentialBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noopBack, nil
	}

	var secretWritingPolicy = `
name = "admins"
path "secret/*" {
	capabilities = ["update", "create", "read"]
}
`

	ps := c.policyStore
	policy, _ := ParseACLPolicy(namespace.RootNamespace, secretWritingPolicy)
	if err := ps.SetPolicy(namespace.RootContext(nil), policy); err != nil {
		t.Fatal(err)
	}

	// Enable the credential backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/auth/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to login -- should fail because we don't allow root to be returned
	lreq := &logical.Request{
		Path: "auth/foo/login",
	}
	lresp, err := c.HandleRequest(namespace.RootContext(nil), lreq)
	if err == nil || lresp == nil || !lresp.IsError() {
		t.Fatalf("expected error trying to auth and receive root policy")
	}

	// Fix and try again
	noopBack.Response.Auth.Policies = []string{"admins"}
	lreq = &logical.Request{
		Path: "auth/foo/login",
	}
	lresp, err = c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Create a leasable secret
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "secret/test",
		Data: map[string]interface{}{
			"foo":   "bar",
			"lease": "1h",
		},
		ClientToken: lresp.Auth.ClientToken,
	}
	resp, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the key
	req.Operation = logical.ReadOperation
	req.Data = nil
	req.SetTokenEntry(&logical.TokenEntry{ID: root, NamespaceID: "root", Policies: []string{"root"}})
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp.Secret)
	}

	// Renew the lease
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/leases/renew")
	req.Data = map[string]interface{}{
		"lease_id": resp.Secret.LeaseID,
	}
	req.ClientToken = lresp.Auth.ClientToken
	_, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Disable the credential backend
	req = logical.TestRequest(t, logical.DeleteOperation, "sys/auth/foo")
	req.ClientToken = root
	resp, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp)
	}
}

func TestCore_HandleRequest_MountPointType(t *testing.T) {
	noop := &NoopBackend{
		Response: &logical.Response{},
	}
	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the logical backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo")
	req.Data["type"] = "noop"
	req.Data["description"] = "foo"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to request
	req = &logical.Request{
		Operation:  logical.ReadOperation,
		Path:       "foo/test",
		Connection: &logical.Connection{},
	}
	req.ClientToken = root
	if _, err := c.HandleRequest(namespace.RootContext(nil), req); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify Path, MountPoint, and MountType
	if noop.Requests[0].Path != "test" {
		t.Fatalf("bad: %#v", noop.Requests)
	}
	if noop.Requests[0].MountPoint != "foo/" {
		t.Fatalf("bad: %#v", noop.Requests)
	}
	if noop.Requests[0].MountType != "noop" {
		t.Fatalf("bad: %#v", noop.Requests)
	}
}

func TestCore_Standby_Rotate(t *testing.T) {
	// Create the first core and initialize it
	logger = logging.NewVaultLogger(log.Trace)

	inm, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	redirectOriginal := "http://127.0.0.1:8200"
	core, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	keys, root := TestCoreInit(t, core)
	for _, key := range keys {
		if _, err := TestCoreUnseal(core, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Wait for core to become active
	TestWaitActive(t, core)

	// Create a second core, attached to same in-memory store
	redirectOriginal2 := "http://127.0.0.1:8500"
	core2, err := NewCore(&CoreConfig{
		Physical:     inm,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: redirectOriginal2,
		DisableMlock: true,
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	for _, key := range keys {
		if _, err := TestCoreUnseal(core2, TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Rotate the encryption key
	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/rotate",
		ClientToken: root,
	}
	_, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Seal the first core, should step down
	err = core.Seal(root)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait for core2 to become active
	TestWaitActive(t, core2)

	// Read the key status
	req = &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "sys/key-status",
		ClientToken: root,
	}
	resp, err := core2.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Verify the response
	if resp.Data["term"] != 2 {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestCore_HandleRequest_Headers(t *testing.T) {
	noop := &NoopBackend{
		Response: &logical.Response{
			Data: map[string]interface{}{},
		},
	}

	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Mount tune
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo/tune")
	req.Data["passthrough_request_headers"] = []string{"Should-Passthrough", "should-passthrough-case-insensitive"}
	req.ClientToken = root
	_, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to read
	lreq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "foo/test",
		ClientToken: root,
		Headers: map[string][]string{
			"Should-Passthrough":                  []string{"foo"},
			"Should-Passthrough-Case-Insensitive": []string{"baz"},
			"Should-Not-Passthrough":              []string{"bar"},
			consts.AuthHeaderName:                 []string{"nope"},
		},
	}
	_, err = c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the headers
	headers := noop.Requests[0].Headers

	// Test passthrough values
	if val, ok := headers["Should-Passthrough"]; ok {
		expected := []string{"foo"}
		if !reflect.DeepEqual(val, expected) {
			t.Fatalf("expected: %v, got: %v", expected, val)
		}
	} else {
		t.Fatalf("expected 'Should-Passthrough' to be present in the headers map")
	}

	if val, ok := headers["Should-Passthrough-Case-Insensitive"]; ok {
		expected := []string{"baz"}
		if !reflect.DeepEqual(val, expected) {
			t.Fatalf("expected: %v, got: %v", expected, val)
		}
	} else {
		t.Fatal("expected 'Should-Passthrough-Case-Insensitive' to be present in the headers map")
	}

	if _, ok := headers["Should-Not-Passthrough"]; ok {
		t.Fatal("did not expect 'Should-Not-Passthrough' to be in the headers map")
	}

	if _, ok := headers[consts.AuthHeaderName]; ok {
		t.Fatalf("did not expect %q to be in the headers map", consts.AuthHeaderName)
	}
}

func TestCore_HandleRequest_Headers_denyList(t *testing.T) {
	noop := &NoopBackend{
		Response: &logical.Response{
			Data: map[string]interface{}{},
		},
	}

	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(context.Context, *logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Enable the backend
	req := logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo")
	req.Data["type"] = "noop"
	req.ClientToken = root
	_, err := c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Mount tune
	req = logical.TestRequest(t, logical.UpdateOperation, "sys/mounts/foo/tune")
	req.Data["passthrough_request_headers"] = []string{"Authorization", consts.AuthHeaderName}
	req.ClientToken = root
	_, err = c.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to read
	lreq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "foo/test",
		ClientToken: root,
		Headers: map[string][]string{
			consts.AuthHeaderName: []string{"foo"},
		},
	}
	_, err = c.HandleRequest(namespace.RootContext(nil), lreq)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Check the headers
	headers := noop.Requests[0].Headers

	// Test passthrough values, they should not be present in the backend
	if _, ok := headers[consts.AuthHeaderName]; ok {
		t.Fatalf("did not expect %q to be in the headers map", consts.AuthHeaderName)
	}
}
