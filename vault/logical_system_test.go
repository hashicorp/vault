package vault

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

func TestSystemBackend_RootPaths(t *testing.T) {
	expected := []string{
		"auth/*",
		"remount",
		"audit",
		"audit/*",
		"raw",
		"raw/*",
		"replication/primary/secondary-token",
		"replication/reindex",
		"rotate",
		"config/cors",
		"config/auditing/*",
		"plugins/catalog/*",
		"revoke-prefix/*",
		"revoke-force/*",
		"leases/revoke-prefix/*",
		"leases/revoke-force/*",
		"leases/lookup/*",
	}

	b := testSystemBackend(t)
	actual := b.SpecialPaths().Root
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSystemConfigCORS(t *testing.T) {
	b := testSystemBackend(t)
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "")
	b.(*SystemBackend).Core.systemBarrierView = view

	req := logical.TestRequest(t, logical.UpdateOperation, "config/cors")
	req.Data["allowed_origins"] = "http://www.example.com"
	req.Data["allowed_headers"] = "X-Custom-Header"
	_, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	expected := &logical.Response{
		Data: map[string]interface{}{
			"enabled":         true,
			"allowed_origins": []string{"http://www.example.com"},
			"allowed_headers": append(StdAllowedHeaders, "X-Custom-Header"),
		},
	}

	req = logical.TestRequest(t, logical.ReadOperation, "config/cors")
	actual, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}

	req = logical.TestRequest(t, logical.DeleteOperation, "config/cors")
	_, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "config/cors")
	actual, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected = &logical.Response{
		Data: map[string]interface{}{
			"enabled": false,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("DELETE FAILED -- bad: %#v", actual)
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
			"type":        "kv",
			"description": "key/value secret storage",
			"accessor":    resp.Data["secret/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
			"accessor":    resp.Data["sys/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"accessor":    resp.Data["cubbyhole/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"accessor":    resp.Data["identity/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["identity/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["identity/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", resp.Data, exp)
	}
}

func TestSystemBackend_mount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/prod/secret/")
	req.Data["type"] = "kv"
	req.Data["local"] = true
	req.Data["seal_wrap"] = true

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "mounts")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// We can't know the pointer address ahead of time so simply
	// copy what's given
	exp := map[string]interface{}{
		"secret/": map[string]interface{}{
			"type":        "kv",
			"description": "key/value secret storage",
			"accessor":    resp.Data["secret/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
			"accessor":    resp.Data["sys/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"accessor":    resp.Data["cubbyhole/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"accessor":    resp.Data["identity/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["identity/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["identity/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     false,
			"seal_wrap": false,
		},
		"prod/secret/": map[string]interface{}{
			"description": "",
			"type":        "kv",
			"accessor":    resp.Data["prod/secret/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["identity/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["identity/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"plugin_name":       "",
				"force_no_cache":    false,
			},
			"local":     true,
			"seal_wrap": true,
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", resp.Data, exp)
	}

}

func TestSystemBackend_mount_force_no_cache(t *testing.T) {
	core, b, _ := testCoreSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/prod/secret/")
	req.Data["type"] = "kv"
	req.Data["config"] = map[string]interface{}{
		"force_no_cache": true,
	}

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	mountEntry := core.router.MatchingMountEntry("prod/secret/")
	if mountEntry == nil {
		t.Fatalf("missing mount entry")
	}
	if !mountEntry.Config.ForceNoCache {
		t.Fatalf("bad config %#v", mountEntry)
	}
}

func TestSystemBackend_mount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/prod/secret/")
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

var capabilitiesPolicy = `
name = "test"
path "foo/bar*" {
	capabilities = ["create", "sudo", "update"]
}
path "sys/capabilities*" {
	capabilities = ["update"]
}
`

func TestSystemBackend_Capabilities(t *testing.T) {
	testCapabilities(t, "capabilities")
	testCapabilities(t, "capabilities-self")
}

func testCapabilities(t *testing.T, endpoint string) {
	core, b, rootToken := testCoreSystemBackend(t)
	req := logical.TestRequest(t, logical.UpdateOperation, endpoint)
	req.Data["token"] = rootToken
	req.Data["path"] = "any_path"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}

	actual := resp.Data["capabilities"]
	expected := []string{"root"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	policy, _ := ParseACLPolicy(capabilitiesPolicy)
	err = core.policyStore.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	testMakeToken(t, core.tokenStore, rootToken, "tokenid", "", []string{"test"})
	req = logical.TestRequest(t, logical.UpdateOperation, endpoint)
	req.Data["token"] = "tokenid"
	req.Data["path"] = "foo/bar"

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}

	actual = resp.Data["capabilities"]
	expected = []string{"create", "sudo", "update"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}

func TestSystemBackend_CapabilitiesAccessor(t *testing.T) {
	core, b, rootToken := testCoreSystemBackend(t)
	te, err := core.tokenStore.Lookup(rootToken)
	if err != nil {
		t.Fatal(err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "capabilities-accessor")
	// Accessor of root token
	req.Data["accessor"] = te.Accessor
	req.Data["path"] = "any_path"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}

	actual := resp.Data["capabilities"]
	expected := []string{"root"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	policy, _ := ParseACLPolicy(capabilitiesPolicy)
	err = core.policyStore.SetPolicy(policy)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	testMakeToken(t, core.tokenStore, rootToken, "tokenid", "", []string{"test"})

	te, err = core.tokenStore.Lookup("tokenid")
	if err != nil {
		t.Fatal(err)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "capabilities-accessor")
	req.Data["accessor"] = te.Accessor
	req.Data["path"] = "foo/bar"

	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatalf("bad: %v", resp)
	}

	actual = resp.Data["capabilities"]
	expected = []string{"create", "sudo", "update"}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}

func TestSystemBackend_remount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "remount")
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

	req := logical.TestRequest(t, logical.UpdateOperation, "remount")
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

	req := logical.TestRequest(t, logical.UpdateOperation, "remount")
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

func TestSystemBackend_leases(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
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

	// Read lease
	req = logical.TestRequest(t, logical.UpdateOperation, "leases/lookup")
	req.Data["lease_id"] = resp.Secret.LeaseID
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["renewable"] == nil || resp.Data["renewable"].(bool) {
		t.Fatal("kv leases are not renewable")
	}

	// Invalid lease
	req = logical.TestRequest(t, logical.UpdateOperation, "leases/lookup")
	req.Data["lease_id"] = "invalid"
	resp, err = b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("expected invalid request, got err: %v", err)
	}
}

func TestSystemBackend_leases_list(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
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

	// List top level
	req = logical.TestRequest(t, logical.ListOperation, "leases/lookup/")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	var keys []string
	if err := mapstructure.WeakDecode(resp.Data["keys"], &keys); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("Expected 1 subkey lease, got %d: %#v", len(keys), keys)
	}
	if keys[0] != "secret/" {
		t.Fatal("Expected only secret subkey")
	}

	// List lease
	req = logical.TestRequest(t, logical.ListOperation, "leases/lookup/secret/foo")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	keys = []string{}
	if err := mapstructure.WeakDecode(resp.Data["keys"], &keys); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("Expected 1 secret lease, got %d: %#v", len(keys), keys)
	}

	// Generate multiple leases
	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "secret/foo")
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.ListOperation, "leases/lookup/secret/foo")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	keys = []string{}
	if err := mapstructure.WeakDecode(resp.Data["keys"], &keys); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 3 {
		t.Fatalf("Expected 3 secret lease, got %d: %#v", len(keys), keys)
	}

	// Listing subkeys
	req = logical.TestRequest(t, logical.UpdateOperation, "secret/bar")
	req.Data["foo"] = "bar"
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Read a key with a LeaseID
	req = logical.TestRequest(t, logical.ReadOperation, "secret/bar")
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Secret == nil || resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	req = logical.TestRequest(t, logical.ListOperation, "leases/lookup/secret")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.Data == nil {
		t.Fatalf("bad: %#v", resp)
	}
	keys = []string{}
	if err := mapstructure.WeakDecode(resp.Data["keys"], &keys); err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("Expected 2 secret lease, got %d: %#v", len(keys), keys)
	}
	expected := []string{"bar/", "foo/"}
	if !reflect.DeepEqual(expected, keys) {
		t.Fatalf("exp: %#v, act: %#v", expected, keys)
	}
}

func TestSystemBackend_renew(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
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
	req2 := logical.TestRequest(t, logical.UpdateOperation, "leases/renew/"+resp.Secret.LeaseID)
	resp2, err := b.HandleRequest(req2)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}

	// Should get error about non-renewability
	if resp2.Data["error"] != "lease is not renewable" {
		t.Fatalf("bad: %#v", resp)
	}

	// Add a TTL to the lease
	req = logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
	req.Data["foo"] = "bar"
	req.Data["ttl"] = "180s"
	req.ClientToken = root
	resp, err = core.HandleRequest(req)
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
	req2 = logical.TestRequest(t, logical.UpdateOperation, "leases/renew/"+resp.Secret.LeaseID)
	resp2, err = b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp2.IsError() {
		t.Fatalf("got an error")
	}
	if resp2.Data == nil {
		t.Fatal("nil data")
	}
	if resp.Secret.TTL != 180*time.Second {
		t.Fatalf("bad lease duration: %v", resp.Secret.TTL)
	}

	// Test the other route path
	req2 = logical.TestRequest(t, logical.UpdateOperation, "leases/renew")
	req2.Data["lease_id"] = resp.Secret.LeaseID
	resp2, err = b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp2.IsError() {
		t.Fatalf("got an error")
	}
	if resp2.Data == nil {
		t.Fatal("nil data")
	}
	if resp.Secret.TTL != 180*time.Second {
		t.Fatalf("bad lease duration: %v", resp.Secret.TTL)
	}

	// Test orig path
	req2 = logical.TestRequest(t, logical.UpdateOperation, "renew")
	req2.Data["lease_id"] = resp.Secret.LeaseID
	resp2, err = b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp2.IsError() {
		t.Fatalf("got an error")
	}
	if resp2.Data == nil {
		t.Fatal("nil data")
	}
	if resp.Secret.TTL != 180*time.Second {
		t.Fatalf("bad lease duration: %v", resp.Secret.TTL)
	}
}

func TestSystemBackend_renew_invalidID(t *testing.T) {
	b := testSystemBackend(t)

	// Attempt renew
	req := logical.TestRequest(t, logical.UpdateOperation, "leases/renew/foobarbaz")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}

	// Attempt renew with other method
	req = logical.TestRequest(t, logical.UpdateOperation, "leases/renew")
	req.Data["lease_id"] = "foobarbaz"
	resp, err = b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_renew_invalidID_origUrl(t *testing.T) {
	b := testSystemBackend(t)

	// Attempt renew
	req := logical.TestRequest(t, logical.UpdateOperation, "renew/foobarbaz")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}

	// Attempt renew with other method
	req = logical.TestRequest(t, logical.UpdateOperation, "renew")
	req.Data["lease_id"] = "foobarbaz"
	resp, err = b.HandleRequest(req)
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
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
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
	req2 := logical.TestRequest(t, logical.UpdateOperation, "revoke/"+resp.Secret.LeaseID)
	resp2, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt renew
	req3 := logical.TestRequest(t, logical.UpdateOperation, "renew/"+resp.Secret.LeaseID)
	resp3, err := b.HandleRequest(req3)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp3.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
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

	// Test the other route path
	req2 = logical.TestRequest(t, logical.UpdateOperation, "revoke")
	req2.Data["lease_id"] = resp.Secret.LeaseID
	resp2, err = b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
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

	// Test the other route path
	req2 = logical.TestRequest(t, logical.UpdateOperation, "leases/revoke")
	req2.Data["lease_id"] = resp.Secret.LeaseID
	resp2, err = b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestSystemBackend_revoke_invalidID(t *testing.T) {
	b := testSystemBackend(t)

	// Attempt revoke
	req := logical.TestRequest(t, logical.UpdateOperation, "leases/revoke/foobarbaz")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Attempt revoke with other method
	req = logical.TestRequest(t, logical.UpdateOperation, "leases/revoke")
	req.Data["lease_id"] = "foobarbaz"
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_revoke_invalidID_origUrl(t *testing.T) {
	b := testSystemBackend(t)

	// Attempt revoke
	req := logical.TestRequest(t, logical.UpdateOperation, "revoke/foobarbaz")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Attempt revoke with other method
	req = logical.TestRequest(t, logical.UpdateOperation, "revoke")
	req.Data["lease_id"] = "foobarbaz"
	resp, err = b.HandleRequest(req)
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
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
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
	req2 := logical.TestRequest(t, logical.UpdateOperation, "leases/revoke-prefix/secret/")
	resp2, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt renew
	req3 := logical.TestRequest(t, logical.UpdateOperation, "leases/renew/"+resp.Secret.LeaseID)
	resp3, err := b.HandleRequest(req3)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp3.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_revokePrefix_origUrl(t *testing.T) {
	core, b, root := testCoreSystemBackend(t)

	// Create a key with a lease
	req := logical.TestRequest(t, logical.UpdateOperation, "secret/foo")
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
	req2 := logical.TestRequest(t, logical.UpdateOperation, "revoke-prefix/secret/")
	resp2, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp2)
	}
	if resp2 != nil {
		t.Fatalf("bad: %#v", resp)
	}

	// Attempt renew
	req3 := logical.TestRequest(t, logical.UpdateOperation, "renew/"+resp.Secret.LeaseID)
	resp3, err := b.HandleRequest(req3)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp3.Data["error"] != "lease not found or lease is not renewable" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_revokePrefixAuth(t *testing.T) {
	core, ts, _, _ := TestCoreWithTokenStore(t)
	bc := &logical.BackendConfig{
		Logger: core.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}
	b := NewSystemBackend(core)
	err := b.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	exp := ts.expiration

	te := &TokenEntry{
		ID:   "foo",
		Path: "auth/github/login/bar",
	}
	err = ts.create(te)
	if err != nil {
		t.Fatal(err)
	}

	te, err = ts.Lookup("foo")
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("token entry was nil")
	}

	// Create a new token
	auth := &logical.Auth{
		ClientToken: te.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
	}
	err = exp.RegisterAuth(te.Path, auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "leases/revoke-prefix/auth/github/")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	te, err = ts.Lookup(te.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te != nil {
		t.Fatalf("bad: %v", te)
	}
}

func TestSystemBackend_revokePrefixAuth_origUrl(t *testing.T) {
	core, ts, _, _ := TestCoreWithTokenStore(t)
	bc := &logical.BackendConfig{
		Logger: core.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}
	b := NewSystemBackend(core)
	err := b.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	exp := ts.expiration

	te := &TokenEntry{
		ID:   "foo",
		Path: "auth/github/login/bar",
	}
	err = ts.create(te)
	if err != nil {
		t.Fatal(err)
	}

	te, err = ts.Lookup("foo")
	if err != nil {
		t.Fatal(err)
	}
	if te == nil {
		t.Fatal("token entry was nil")
	}

	// Create a new token
	auth := &logical.Auth{
		ClientToken: te.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
	}
	err = exp.RegisterAuth(te.Path, auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "revoke-prefix/auth/github/")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp != nil {
		t.Fatalf("bad: %#v", resp)
	}

	te, err = ts.Lookup(te.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if te != nil {
		t.Fatalf("bad: %v", te)
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
		"token/": map[string]interface{}{
			"type":        "token",
			"description": "token based credentials",
			"accessor":    resp.Data["token/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": int64(0),
				"max_lease_ttl":     int64(0),
			},
			"local":     false,
			"seal_wrap": false,
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

	req := logical.TestRequest(t, logical.UpdateOperation, "auth/foo")
	req.Data["type"] = "noop"
	req.Data["local"] = true
	req.Data["seal_wrap"] = true

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "auth")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatal("resp is nil")
	}

	exp := map[string]interface{}{
		"foo/": map[string]interface{}{
			"type":        "noop",
			"description": "",
			"accessor":    resp.Data["foo/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": int64(0),
				"max_lease_ttl":     int64(0),
			},
			"local":     true,
			"seal_wrap": true,
		},
		"token/": map[string]interface{}{
			"type":        "token",
			"description": "token based credentials",
			"accessor":    resp.Data["token/"].(map[string]interface{})["accessor"],
			"config": map[string]interface{}{
				"default_lease_ttl": int64(0),
				"max_lease_ttl":     int64(0),
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_enableAuth_invalid(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/foo")
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
	req := logical.TestRequest(t, logical.UpdateOperation, "auth/foo")
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

func TestSystemBackend_policyList(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.ReadOperation, "policy")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"keys":     []string{"default", "root"},
		"policies": []string{"default", "root"},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_policyCRUD(t *testing.T) {
	b := testSystemBackend(t)

	// Create the policy
	rules := `path "foo/" { policy = "read" }`
	req := logical.TestRequest(t, logical.UpdateOperation, "policy/Foo")
	req.Data["rules"] = rules
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v %#v", err, resp)
	}
	if resp != nil && (resp.IsError() || len(resp.Data) > 0) {
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
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp = map[string]interface{}{
		"name":  "foo",
		"rules": rules,
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}

	// List the policies
	req = logical.TestRequest(t, logical.ReadOperation, "policy")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp = map[string]interface{}{
		"keys":     []string{"default", "foo", "root"},
		"policies": []string{"default", "foo", "root"},
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
		"keys":     []string{"default", "root"},
		"policies": []string{"default", "root"},
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

	req := logical.TestRequest(t, logical.UpdateOperation, "audit/foo")
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
		config.SaltView = view
		config.SaltConfig = &salt.Config{
			HMAC:     sha256.New,
			HMACType: "hmac-sha256",
			Location: salt.DefaultLocation,
		}
		return &NoopAudit{
			Config: config,
		}, nil
	}

	req := logical.TestRequest(t, logical.UpdateOperation, "audit/foo")
	req.Data["type"] = "noop"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.UpdateOperation, "audit-hash/foo")
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
	req := logical.TestRequest(t, logical.UpdateOperation, "audit/foo")
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

	req := logical.TestRequest(t, logical.UpdateOperation, "audit/foo")
	req.Data["type"] = "noop"
	req.Data["description"] = "testing"
	req.Data["options"] = map[string]interface{}{
		"foo": "bar",
	}
	req.Data["local"] = true
	b.HandleRequest(req)

	req = logical.TestRequest(t, logical.ReadOperation, "audit")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"foo/": map[string]interface{}{
			"path":        "foo/",
			"type":        "noop",
			"description": "testing",
			"options": map[string]string{
				"foo": "bar",
			},
			"local": true,
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

	req := logical.TestRequest(t, logical.UpdateOperation, "audit/foo")
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

func TestSystemBackend_rawRead_Protected(t *testing.T) {
	b := testSystemBackendRaw(t)

	req := logical.TestRequest(t, logical.ReadOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawWrite_Protected(t *testing.T) {
	b := testSystemBackendRaw(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawReadWrite(t *testing.T) {
	_, b, _ := testCoreSystemBackendRaw(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "raw/sys/policy/test")
	req.Data["value"] = `path "secret/" { policy = "read" }`
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	// Read via raw API
	req = logical.TestRequest(t, logical.ReadOperation, "raw/sys/policy/test")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !strings.HasPrefix(resp.Data["value"].(string), "path") {
		t.Fatalf("bad: %v", resp)
	}

	// Note: since the upgrade code is gone that upgraded from 0.1, we can't
	// simply parse this out directly via GetPolicy, so the test now ends here.
}

func TestSystemBackend_rawDelete_Protected(t *testing.T) {
	b := testSystemBackendRaw(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawDelete(t *testing.T) {
	c, b, _ := testCoreSystemBackendRaw(t)

	// set the policy!
	p := &Policy{
		Name: "test",
		Type: PolicyTypeACL,
	}
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
	c.policyStore.tokenPoliciesLRU.Purge()
	out, err := c.policyStore.GetPolicy("test", PolicyTypeToken)
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

	req := logical.TestRequest(t, logical.UpdateOperation, "rotate")
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
	return testSystemBackendInternal(t, c)
}

func testSystemBackendRaw(t *testing.T) logical.Backend {
	c, _, _ := TestCoreUnsealedRaw(t)
	return testSystemBackendInternal(t, c)
}

func testCoreSystemBackend(t *testing.T) (*Core, logical.Backend, string) {
	c, _, root := TestCoreUnsealed(t)
	return c, testSystemBackendInternal(t, c), root
}

func testCoreSystemBackendRaw(t *testing.T) (*Core, logical.Backend, string) {
	c, _, root := TestCoreUnsealedRaw(t)
	return c, testSystemBackendInternal(t, c), root
}

func testSystemBackendInternal(t *testing.T, c *Core) logical.Backend {
	bc := &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}

	b := NewSystemBackend(c)
	err := b.Backend.Setup(bc)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestSystemBackend_PluginCatalog_CRUD(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)
	// Bootstrap the pluginCatalog
	sym, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	c.pluginCatalog.directory = sym

	req := logical.TestRequest(t, logical.ListOperation, "plugins/catalog/")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(resp.Data["keys"].([]string)) != len(builtinplugins.Keys()) {
		t.Fatalf("Wrong number of plugins, got %d, expected %d", len(resp.Data["keys"].([]string)), len(builtinplugins.Keys()))
	}

	req = logical.TestRequest(t, logical.ReadOperation, "plugins/catalog/mysql-database-plugin")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	actualRespData := resp.Data

	expectedBuiltin := &pluginutil.PluginRunner{
		Name:    "mysql-database-plugin",
		Builtin: true,
	}
	expectedRespData := structs.New(expectedBuiltin).Map()

	if !reflect.DeepEqual(actualRespData, expectedRespData) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", actualRespData, expectedRespData)
	}

	// Set a plugin
	file, err := ioutil.TempFile(os.TempDir(), "temp")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	command := fmt.Sprintf("%s --test", filepath.Base(file.Name()))
	req = logical.TestRequest(t, logical.UpdateOperation, "plugins/catalog/test-plugin")
	req.Data["sha_256"] = hex.EncodeToString([]byte{'1'})
	req.Data["command"] = command
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "plugins/catalog/test-plugin")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	actual := resp.Data

	expectedRunner := &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Command: filepath.Join(sym, filepath.Base(file.Name())),
		Args:    []string{"--test"},
		Sha256:  []byte{'1'},
		Builtin: false,
	}
	expected := structs.New(expectedRunner).Map()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", actual, expected)
	}

	// Delete plugin
	req = logical.TestRequest(t, logical.DeleteOperation, "plugins/catalog/test-plugin")
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "plugins/catalog/test-plugin")
	resp, err = b.HandleRequest(req)
	if resp != nil || err != nil {
		t.Fatalf("expected nil response, plugin not deleted correctly got resp: %v, err: %v", resp, err)
	}
}

func TestSystemBackend_ToolsHash(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.UpdateOperation, "tools/hash")
	req.Data = map[string]interface{}{
		"input": "dGhlIHF1aWNrIGJyb3duIGZveA==",
	}
	_, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	doRequest := func(req *logical.Request, errExpected bool, expected string) {
		t.Helper()
		resp, err := b.HandleRequest(req)
		if err != nil && !errExpected {
			t.Fatal(err)
		}
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if errExpected {
			if !resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
			}
			return
		}
		if resp.IsError() {
			t.Fatalf("bad: got error response: %#v", *resp)
		}
		sum, ok := resp.Data["sum"]
		if !ok {
			t.Fatal("no sum key found in returned data")
		}
		if sum.(string) != expected {
			t.Fatal("mismatched hashes")
		}
	}

	// Test defaults -- sha2-256
	doRequest(req, false, "9ecb36561341d18eb65484e833efea61edc74b84cf5e6ae1b81c63533e25fc8f")

	// Test algorithm selection in the path
	req.Path = "tools/hash/sha2-224"
	doRequest(req, false, "ea074a96cabc5a61f8298a2c470f019074642631a49e1c5e2f560865")

	// Reset and test algorithm selection in the data
	req.Path = "tools/hash"
	req.Data["algorithm"] = "sha2-224"
	doRequest(req, false, "ea074a96cabc5a61f8298a2c470f019074642631a49e1c5e2f560865")

	req.Data["algorithm"] = "sha2-384"
	doRequest(req, false, "15af9ec8be783f25c583626e9491dbf129dd6dd620466fdf05b3a1d0bb8381d30f4d3ec29f923ff1e09a0f6b337365a6")

	req.Data["algorithm"] = "sha2-512"
	doRequest(req, false, "d9d380f29b97ad6a1d92e987d83fa5a02653301e1006dd2bcd51afa59a9147e9caedaf89521abc0f0b682adcd47fb512b8343c834a32f326fe9bef00542ce887")

	// Test returning as base64
	req.Data["format"] = "base64"
	doRequest(req, false, "2dOA8puXrWodkumH2D+loCZTMB4QBt0rzVGvpZqRR+nK7a+JUhq8DwtoKtzUf7USuDQ8g0oy8yb+m+8AVCzohw==")

	// Test bad input/format/algorithm
	req.Data["format"] = "base92"
	doRequest(req, true, "")

	req.Data["format"] = "hex"
	req.Data["algorithm"] = "foobar"
	doRequest(req, true, "")

	req.Data["algorithm"] = "sha2-256"
	req.Data["input"] = "foobar"
	doRequest(req, true, "")
}

func TestSystemBackend_ToolsRandom(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.UpdateOperation, "tools/random")

	_, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	doRequest := func(req *logical.Request, errExpected bool, format string, numBytes int) {
		t.Helper()
		getResponse := func() []byte {
			resp, err := b.HandleRequest(req)
			if err != nil && !errExpected {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if errExpected {
				if !resp.IsError() {
					t.Fatalf("bad: got error response: %#v", *resp)
				}
				return nil
			}
			if resp.IsError() {
				t.Fatalf("bad: got error response: %#v", *resp)
			}
			if _, ok := resp.Data["random_bytes"]; !ok {
				t.Fatal("no random_bytes found in response")
			}

			outputStr := resp.Data["random_bytes"].(string)
			var outputBytes []byte
			switch format {
			case "base64":
				outputBytes, err = base64.StdEncoding.DecodeString(outputStr)
			case "hex":
				outputBytes, err = hex.DecodeString(outputStr)
			default:
				t.Fatal("unknown format")
			}
			if err != nil {
				t.Fatal(err)
			}

			return outputBytes
		}

		rand1 := getResponse()
		// Expected error
		if rand1 == nil {
			return
		}
		rand2 := getResponse()
		if len(rand1) != numBytes || len(rand2) != numBytes {
			t.Fatal("length of output random bytes not what is exepcted")
		}
		if reflect.DeepEqual(rand1, rand2) {
			t.Fatal("found identical ouputs")
		}
	}

	// Test defaults
	doRequest(req, false, "base64", 32)

	// Test size selection in the path
	req.Path = "tools/random/24"
	req.Data["format"] = "hex"
	doRequest(req, false, "hex", 24)

	// Test bad input/format
	req.Path = "tools/random"
	req.Data["format"] = "base92"
	doRequest(req, true, "", 0)

	req.Data["format"] = "hex"
	req.Data["bytes"] = -1
	doRequest(req, true, "", 0)
}
