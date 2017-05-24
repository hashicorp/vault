package vault

import (
	"crypto/sha256"
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
		"raw/*",
		"replication/primary/secondary-token",
		"replication/reindex",
		"rotate",
		"config/auditing/*",
		"plugins/catalog/*",
		"revoke-prefix/*",
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
				"default_lease_ttl": resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["secret/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"force_no_cache":    false,
			},
			"local": false,
		},
		"sys/": map[string]interface{}{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["sys/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"force_no_cache":    false,
			},
			"local": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["default_lease_ttl"].(int64),
				"max_lease_ttl":     resp.Data["cubbyhole/"].(map[string]interface{})["config"].(map[string]interface{})["max_lease_ttl"].(int64),
				"force_no_cache":    false,
			},
			"local": true,
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("Got:\n%#v\nExpected:\n%#v", resp.Data, exp)
	}
}

func TestSystemBackend_mount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/prod/secret/")
	req.Data["type"] = "generic"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_mount_force_no_cache(t *testing.T) {
	core, b, _ := testCoreSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "mounts/prod/secret/")
	req.Data["type"] = "generic"
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

	policy, _ := Parse(capabilitiesPolicy)
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

	policy, _ := Parse(capabilitiesPolicy)
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
		t.Fatal("generic leases are not renewable")
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
	b, err := NewSystemBackend(core, bc)
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
	b, err := NewSystemBackend(core, bc)
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
			"config": map[string]interface{}{
				"default_lease_ttl": int64(0),
				"max_lease_ttl":     int64(0),
			},
			"local": false,
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
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.ReadOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawWrite_Protected(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.UpdateOperation, "raw/"+keyringPath)
	_, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
}

func TestSystemBackend_rawReadWrite(t *testing.T) {
	c, b, _ := testCoreSystemBackend(t)

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

	// Read the policy!
	p, err := c.policyStore.GetPolicy("test")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p == nil || len(p.Paths) == 0 {
		t.Fatalf("missing policy %#v", p)
	}
	if p.Paths[0].Prefix != "secret/" || p.Paths[0].Policy != ReadCapability {
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
	bc := &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}

	b, err := NewSystemBackend(c, bc)
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func testCoreSystemBackend(t *testing.T) (*Core, logical.Backend, string) {
	c, _, root := TestCoreUnsealed(t)
	bc := &logical.BackendConfig{
		Logger: c.logger,
		System: logical.StaticSystemView{
			DefaultLeaseTTLVal: time.Hour * 24,
			MaxLeaseTTLVal:     time.Hour * 24 * 32,
		},
	}

	b, err := NewSystemBackend(c, bc)
	if err != nil {
		t.Fatal(err)
	}
	return c, b, root
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

	expectedBuiltin := &pluginutil.PluginRunner{
		Name:    "mysql-database-plugin",
		Builtin: true,
	}
	expectedBuiltin.BuiltinFactory, _ = builtinplugins.Get("mysql-database-plugin")

	p := resp.Data["plugin"].(*pluginutil.PluginRunner)
	if &(p.BuiltinFactory) == &(expectedBuiltin.BuiltinFactory) {
		t.Fatal("expected BuiltinFactory did not match actual")
	}

	expectedBuiltin.BuiltinFactory = nil
	p.BuiltinFactory = nil
	if !reflect.DeepEqual(p, expectedBuiltin) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", resp.Data["plugin"].(*pluginutil.PluginRunner), expectedBuiltin)
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

	expected := &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Command: filepath.Join(sym, filepath.Base(file.Name())),
		Args:    []string{"--test"},
		Sha256:  []byte{'1'},
		Builtin: false,
	}
	if !reflect.DeepEqual(resp.Data["plugin"].(*pluginutil.PluginRunner), expected) {
		t.Fatalf("expected did not match actual, got %#v\n expected %#v\n", resp.Data["plugin"].(*pluginutil.PluginRunner), expected)
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
