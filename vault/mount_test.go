package vault

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/logical"
)

func TestCore_DefaultMountTable(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	verifyDefaultTable(t, c.mounts)

	// Start a second core with same physical
	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := TestCoreUnseal(c2, key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Mount(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	me := &MountEntry{
		Table: mountTableType,
		Path:  "foo",
		Type:  "generic",
	}
	err := c.mount(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("foo/bar")
	if match != "foo/" {
		t.Fatalf("missing mount")
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := TestCoreUnseal(c2, key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Unmount(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	existed, err := c.unmount("secret")
	if !existed || err != nil {
		t.Fatalf("existed: %v; err: %v", existed, err)
	}

	match := c.router.MatchingMount("secret/foo")
	if match != "" {
		t.Fatalf("backend present")
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := TestCoreUnseal(c2, key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Unmount_Cleanup(t *testing.T) {
	noop := &NoopBackend{}
	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(*logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Mount the noop backend
	me := &MountEntry{
		Table: mountTableType,
		Path:  "test/",
		Type:  "noop",
	}
	if err := c.mount(me); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Store the view
	view := c.router.MatchingStorageView("test/")

	// Inject data
	se := &logical.StorageEntry{
		Key:   "plstodelete",
		Value: []byte("test"),
	}
	if err := view.Put(se); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Setup response
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	noop.Response = resp

	// Generate leased secret
	r := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "test/foo",
		ClientToken: root,
	}
	resp, err := c.HandleRequest(r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Unmount, this should cleanup
	if existed, err := c.unmount("test/"); !existed || err != nil {
		t.Fatalf("existed: %v; err: %v", existed, err)
	}

	// Rollback should be invoked
	if noop.Requests[1].Operation != logical.RollbackOperation {
		t.Fatalf("bad: %#v", noop.Requests)
	}

	// Revoke should be invoked
	if noop.Requests[2].Operation != logical.RevokeOperation {
		t.Fatalf("bad: %#v", noop.Requests)
	}
	if noop.Requests[2].Path != "foo" {
		t.Fatalf("bad: %#v", noop.Requests)
	}

	// View should be empty
	out, err := CollectKeys(view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %#v", out)
	}
}

func TestCore_Remount(t *testing.T) {
	c, key, _ := TestCoreUnsealed(t)
	err := c.remount("secret", "foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	match := c.router.MatchingMount("foo/bar")
	if match != "foo/" {
		t.Fatalf("failed remount")
	}

	conf := &CoreConfig{
		Physical:     c.physical,
		DisableMlock: true,
	}
	c2, err := NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	unseal, err := TestCoreUnseal(c2, key)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !unseal {
		t.Fatalf("should be unsealed")
	}

	// Verify matching mount tables
	if !reflect.DeepEqual(c.mounts, c2.mounts) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Remount_Cleanup(t *testing.T) {
	noop := &NoopBackend{}
	c, _, root := TestCoreUnsealed(t)
	c.logicalBackends["noop"] = func(*logical.BackendConfig) (logical.Backend, error) {
		return noop, nil
	}

	// Mount the noop backend
	me := &MountEntry{
		Table: mountTableType,
		Path:  "test/",
		Type:  "noop",
	}
	if err := c.mount(me); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Store the view
	view := c.router.MatchingStorageView("test/")

	// Inject data
	se := &logical.StorageEntry{
		Key:   "plstokeep",
		Value: []byte("test"),
	}
	if err := view.Put(se); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Setup response
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
		},
		Data: map[string]interface{}{
			"foo": "bar",
		},
	}
	noop.Response = resp

	// Generate leased secret
	r := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "test/foo",
		ClientToken: root,
	}
	resp, err := c.HandleRequest(r)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp.Secret.LeaseID == "" {
		t.Fatalf("bad: %#v", resp)
	}

	// Remount, this should cleanup
	if err := c.remount("test/", "new/"); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Rollback should be invoked
	if noop.Requests[1].Operation != logical.RollbackOperation {
		t.Fatalf("bad: %#v", noop.Requests)
	}

	// Revoke should be invoked
	if noop.Requests[2].Operation != logical.RevokeOperation {
		t.Fatalf("bad: %#v", noop.Requests)
	}
	if noop.Requests[2].Path != "foo" {
		t.Fatalf("bad: %#v", noop.Requests)
	}

	// View should not be empty
	out, err := CollectKeys(view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 1 && out[0] != "plstokeep" {
		t.Fatalf("bad: %#v", out)
	}
}

func TestCore_Remount_Protected(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)
	err := c.remount("sys", "foo")
	if err.Error() != "cannot remount 'sys/'" {
		t.Fatalf("err: %v", err)
	}
}

func TestDefaultMountTable(t *testing.T) {
	table := defaultMountTable()
	verifyDefaultTable(t, table)
}

func TestCore_MountTable_UpgradeToTyped(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	c.auditBackends["noop"] = func(config *audit.BackendConfig) (audit.Backend, error) {
		return &NoopAudit{
			Config: config,
		}, nil
	}

	me := &MountEntry{
		Table: auditTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err := c.enableAudit(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	c.credentialBackends["noop"] = func(*logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{}, nil
	}

	me = &MountEntry{
		Table: credentialTableType,
		Path:  "foo",
		Type:  "noop",
	}
	err = c.enableCredential(me)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	testCore_MountTable_UpgradeToTyped_Common(t, c, "mounts")
	testCore_MountTable_UpgradeToTyped_Common(t, c, "audits")
	testCore_MountTable_UpgradeToTyped_Common(t, c, "credentials")
}

func testCore_MountTable_UpgradeToTyped_Common(
	t *testing.T,
	c *Core,
	testType string) {

	var path string
	var mt *MountTable
	switch testType {
	case "mounts":
		path = coreMountConfigPath
		mt = c.mounts
	case "audits":
		path = coreAuditConfigPath
		mt = c.audit
	case "credentials":
		path = coreAuthConfigPath
		mt = c.auth
	}

	// Save the expected table
	goodJson, err := json.Marshal(mt)
	if err != nil {
		t.Fatal(err)
	}

	// Create a pre-typed version
	mt.Type = ""
	for _, entry := range mt.Entries {
		entry.Table = ""
	}

	raw, err := json.Marshal(mt)
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(raw, goodJson) {
		t.Fatalf("bad: values here should be different")
	}

	entry := &Entry{
		Key:   path,
		Value: raw,
	}
	if err := c.barrier.Put(entry); err != nil {
		t.Fatal(err)
	}

	var persistFunc func(*MountTable) error

	// It should load successfully and be upgraded and persisted
	switch testType {
	case "mounts":
		err = c.loadMounts()
		persistFunc = c.persistMounts
		mt = c.mounts
	case "credentials":
		err = c.loadCredentials()
		persistFunc = c.persistAuth
		mt = c.auth
	case "audits":
		err = c.loadAudits()
		persistFunc = c.persistAudit
		mt = c.audit
	}
	if err != nil {
		t.Fatal(err)
	}

	entry, err = c.barrier.Get(path)
	if err != nil {
		t.Fatal(err)
	}

	decompressedBytes, uncompressed, err := compressutil.Decompress(entry.Value)
	if err != nil {
		t.Fatal(err)
	}

	actual := decompressedBytes
	if uncompressed {
		actual = entry.Value
	}

	if strings.TrimSpace(string(actual)) != strings.TrimSpace(string(goodJson)) {
		t.Fatalf("bad: expected\n%s\nactual\n%s\n", string(goodJson), string(actual))
	}

	// Now try saving invalid versions
	origTableType := mt.Type
	mt.Type = "foo"
	if err := persistFunc(mt); err == nil {
		t.Fatal("expected error")
	}

	if len(mt.Entries) > 0 {
		mt.Type = origTableType
		mt.Entries[0].Table = "bar"
		if err := persistFunc(mt); err == nil {
			t.Fatal("expected error")
		}

		mt.Entries[0].Table = mt.Type
		if err := persistFunc(mt); err != nil {
			t.Fatal(err)
		}
	}
}

func verifyDefaultTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 3 {
		t.Fatalf("bad: %v", table.Entries)
	}
	for idx, entry := range table.Entries {
		switch idx {
		case 0:
			if entry.Path != "secret/" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "generic" {
				t.Fatalf("bad: %v", entry)
			}
		case 1:
			if entry.Path != "cubbyhole/" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "cubbyhole" {
				t.Fatalf("bad: %v", entry)
			}
		case 2:
			if entry.Path != "sys/" {
				t.Fatalf("bad: %v", entry)
			}
			if entry.Type != "system" {
				t.Fatalf("bad: %v", entry)
			}
		}
		if entry.Table != mountTableType {
			t.Fatalf("bad: %v", entry)
		}
		if entry.Description == "" {
			t.Fatalf("bad: %v", entry)
		}
		if entry.UUID == "" {
			t.Fatalf("bad: %v", entry)
		}
	}

}
