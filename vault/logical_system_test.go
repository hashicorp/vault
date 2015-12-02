package vault

import (
	"crypto/sha256"
	"reflect"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
)

func TestSystemBackend_RootPaths(t *testing.T) {
	expected := []string{
		"auth/*",
		"remount",
		"revoke-prefix/*",
		"audit",
		"audit/*",
		"seal",
		"raw/*",
		"rotate",
	}

	b := testSystemBackend(t)
	actual := b.SpecialPaths().Root
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSystemBackend_mounts(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.ReadOperation, "mounts")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// We can't know the pointer address ahead of time so simply
	// copy what's given
	exp := map[string]interface{}{
		"secret/": map[string]interface{}{
			"type":        "generic",
			"description": "generic secret storage",
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int),
				"max_lease_ttl":     resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int),
			},
		},
		"sys/": map[string]interface{}{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int),
				"max_lease_ttl":     resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int),
			},
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int),
				"max_lease_ttl":     resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int),
			},
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("Got:\n%#v\nExpected:\n%#v", resp.Data, exp)
	}
}

func TestSystemBackend_mount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "mounts/prod/secret/")
	req.Data["type"] = "generic"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_mount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "mounts/prod/secret/")
	req.Data["type"] = "nope"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "unknown backend type: nope" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_unmount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "mounts/secret/")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_unmount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "mounts/foo/")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching mount" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_remount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "remount")
	req.Data["from"] = "secret"
	req.Data["to"] = "foo"
	req.Data["config"] = structs.Map(MountConfig{})
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_remount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "remount")
	req.Data["from"] = "unknown"
	req.Data["to"] = "foo"
	req.Data["config"] = structs.Map(MountConfig{})
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching mount at 'unknown/'" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_remount_system(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "remount")
	req.Data["from"] = "sys"
	req.Data["to"] = "foo"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "cannot remount 'sys/'" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_renew(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.WriteOperation, "secret/foo")
	req.Data["foo"] = "bar"
	req.ClientToken = root
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read a key with a LeaseID
	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt renew
	req2 := logical.TestRequest(t, logical.WriteOperation, "renew/"+resp.Secret.LeaseID)
	req2.Data["increment"] = "100s"
	resp2, err := b.HandleRequest(req2)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}

	// Should get error about non-renewability
	if resp2.Data["error"] != "lease is not renewable" {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestSystemBackend_renew_invalidID(t *testing.T) {
	b := testSystemBackend(t)

	// Attempt renew
	req := logical.TestRequest(t, logical.WriteOperation, "renew/foobarbaz")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_revoke(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.WriteOperation, "secret/foo")
	req.Data["foo"] = "bar"
	req.Data["lease"] = "1h"
	req.ClientToken = root
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read a key with a LeaseID
	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt revoke
	req2 := logical.TestRequest(t, logical.WriteOperation, "revoke/"+resp.Secret.LeaseID)
	resp2, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt renew
	req3 := logical.TestRequest(t, logical.WriteOperation, "renew/"+resp.Secret.LeaseID)
	resp3, err := b.HandleRequest(req3)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp3.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_revoke_invalidID(t *testing.T) {
	b := testSystemBackend(t)

	// Attempt renew
	req := logical.TestRequest(t, logical.WriteOperation, "revoke/foobarbaz")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_revokePrefix(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.WriteOperation, "secret/foo")
	req.Data["foo"] = "bar"
	req.Data["lease"] = "1h"
	req.ClientToken = root
	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read a key with a LeaseID
	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt revoke
	req2 := logical.TestRequest(t, logical.WriteOperation, "revoke-prefix/secret/")
	resp2, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt renew
	req3 := logical.TestRequest(t, logical.WriteOperation, "renew/"+resp.Secret.LeaseID)
	resp3, err := b.HandleRequest(req3)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp3.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_authTable(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.ReadOperation, "auth")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"token/": map[string]string{
			"type":        "token",
			"description": "token based credentials",
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_enableAuth(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	c.credentialBackends["noop"] = func(*logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{}, nil
	}

	req := logical.TestRequest(t, logical.WriteOperation, "auth/foo")
	req.Data["type"] = "noop"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_enableAuth_invalid(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.WriteOperation, "auth/foo")
	req.Data["type"] = "nope"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "unknown backend type: nope" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_disableAuth(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	c.credentialBackends["noop"] = func(*logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{}, nil
	}

	// Register the backend
	req := logical.TestRequest(t, logical.WriteOperation, "auth/foo")
	req.Data["type"] = "noop"
	b.HandleRequest(req)

	// Deregister it
	req = logical.TestRequest(t, logical.DeleteOperation, "auth/foo")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_disableAuth_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "auth/foo")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching backend" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_policyList(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.ReadOperation, "policy")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"keys": []string{"default", "root"},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_policyCRUD(t *testing.T) {
	b := testSystemBackend(t)

	// Create the policy
	rules := `path "foo/" { policy = "read" }`
	req := logical.TestRequest(t, logical.WriteOperation, "policy/Foo")
	req.Data["rules"] = rules
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the policy
	req = logical.TestRequest(t, logical.ReadOperation, "policy/foo")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"name":  "foo",
		"rules": rules,
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}

	// Read, and make sure that case has been normalized
	req = logical.TestRequest(t, logical.ReadOperation, "policy/Foo")
	resp, err = b.HandleRequest(req)
	if resp != nil {
		t.Fatalf("err: expected nil response, got %#v", *resp)
	}

	// List the policies
	req = logical.TestRequest(t, logical.ReadOperation, "policy")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp = map[string]interface{}{
		"keys": []string{"default", "foo", "root"},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}

	// Delete the policy
	req = logical.TestRequest(t, logical.DeleteOperation, "policy/foo")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read the policy (deleted)
	req = logical.TestRequest(t, logical.ReadOperation, "policy/foo")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// List the policies (deleted)
	req = logical.TestRequest(t, logical.ReadOperation, "policy")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp = map[string]interface{}{
		"keys": []string{"default", "root"},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_enableAudit(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	req := logical.TestRequest(t, logical.WriteOperation, "audit/foo")
	req.Data["type"] = "noop"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_auditHash(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		view := &logical.InmemStorage{}
		view.Put(&logical.StorageEntry{
			Key:   "salt",
			Value: []byte("foo"),
		})
		var err error
		config.Salt, err = salt.NewSalt(view, &salt.Config{
			HMAC:     sha256.New,
			HMACType: "hmac-sha256",
		})
		if err != nil {
			t.Fatal("error getting new salt: %v", err)
		}
		return &NoopAudit{
			Config: config,
		}, nil
	}

	req := logical.TestRequest(t, logical.WriteOperation, "audit/foo")
	req.Data["type"] = "noop"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.WriteOperation, "audit-hash/foo")
	req.Data["input"] = "bar"

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("response or its data was nil")
	}
	hash, ok := resp.Data["hash"]
	if !ok {
		t.Fatalf("did not get hash back in response, response was %#v", resp.Data)
	}
	if hash.(string) != "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317" {
		t.Fatalf("bad hash back: %s", hash.(string))
	}
}

func TestSystemBackend_enableAudit_invalid(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.WriteOperation, "audit/foo")
	req.Data["type"] = "nope"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "unknown backend type: nope" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_auditTable(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	req := logical.TestRequest(t, logical.WriteOperation, "audit/foo")
	req.Data["type"] = "noop"
	req.Data["description"] = "testing"
	req.Data["options"] = map[string]interface{}{
		"foo": "bar",
	}
	b.HandleRequest(req)

	req = logical.TestRequest(t, logical.ReadOperation, "audit")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"foo/": map[string]interface{}{
			"type":        "noop",
			"description": "testing",
			"options": map[string]string{
				"foo": "bar",
			},
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_disableAudit(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	req := logical.TestRequest(t, logical.WriteOperation, "audit/foo")
	req.Data["type"] = "noop"
	req.Data["description"] = "testing"
	req.Data["options"] = map[string]interface{}{
		"foo": "bar",
	}
	b.HandleRequest(req)

	// Deregister it
	req = logical.TestRequest(t, logical.DeleteOperation, "audit/foo")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_disableAudit_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "audit/foo")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching backend" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_rawRead_Protected(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.ReadOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawRead(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.ReadOperation, "raw/"+coreMountConfigPath)
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["value"].(string)[0] != '{' {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_rawWrite_Protected(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawWrite(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "raw/sys/policy/test")
	req.Data["value"] = `path "secret/" { policy = "read" }`
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Read the policy!
	p, err := c.policyStore.GetPolicy("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p == nil || len(p.Paths) == 0 {
		t.Fatalf("missing policy %#v", p)
	}
	if p.Paths[0].Prefix != "secret/" || p.Paths[0].Policy != PathPolicyRead {
		t.Fatalf("Bad: %#v", p)
	}
}

func TestSystemBackend_rawDelete_Protected(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawDelete(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)

	// set the policy!
	p := &Policy{Name: "test"}
	err := c.policyStore.SetPolicy(p)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Delete the policy
	req := logical.TestRequest(t, logical.DeleteOperation, "raw/sys/policy/test")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Policy should be gone
	c.policyStore.lru.Purge()
	out, err := c.policyStore.GetPolicy("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("policy should be gone")
	}
}

func TestSystemBackend_keyStatus(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.ReadOperation, "key-status")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"term": 1,
	}
	delete(resp.Data, "install_time")
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_rotate(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "rotate")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "key-status")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"term": 2,
	}
	delete(resp.Data, "install_time")
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func testSystemBackend(t *testing.T) logical.Backend {
	c, _, _ := TestCoreUnsealed(t)
	bc := &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 30,
		},
	}
	return NewSystemBackend(c, bc)
}

func testCoreSystemBackend(t *testing.T) (*Core, logical.Backend, string) {
	c, _, root := TestCoreUnsealed(t)
	bc := &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 30,
		},
	}
	return c, NewSystemBackend(c, bc), root
}
