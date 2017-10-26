package vault

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/compressutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
)

func TestCore_DefaultMountTable(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
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
	if !reflect.DeepEqual(c.mounts.sortEntriesByPath(), c2.mounts.sortEntriesByPath()) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

func TestCore_Mount(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	me := &MountEntry{
		Table: mountTableType,
		Path:  "foo",
		Type:  "kv",
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
	if !reflect.DeepEqual(c.mounts.sortEntriesByPath(), c2.mounts.sortEntriesByPath()) {
		t.Fatalf("mismatch: %v %v", c.mounts, c2.mounts)
	}
}

// Test that the local table actually gets populated as expected with local
// entries, and that upon reading the entries from both are recombined
// correctly
func TestCore_Mount_Local(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	c.mounts = &MountTable{
		Type: mountTableType,
		Entries: []*MountEntry{
			&MountEntry{
				Table:    mountTableType,
				Path:     "noop/",
				Type:     "kv",
				UUID:     "abcd",
				Accessor: "kv-abcd",
			},
			&MountEntry{
				Table:    mountTableType,
				Path:     "noop2/",
				Type:     "kv",
				UUID:     "bcde",
				Accessor: "kv-bcde",
			},
		},
	}

	// Both should set up successfully
	err := c.setupMounts()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.mounts.Entries) != 2 {
		t.Fatalf("expected two entries, got %d", len(c.mounts.Entries))
	}

	rawLocal, err := c.barrier.Get(coreLocalMountConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local mounts")
	}
	localMountsTable := &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localMountsTable); err != nil {
		t.Fatal(err)
	}
	if len(localMountsTable.Entries) != 1 || localMountsTable.Entries[0].Type != "cubbyhole" {
		t.Fatalf("expected only cubbyhole entry in local mount table, got %#v", localMountsTable)
	}

	c.mounts.Entries[1].Local = true
	if err := c.persistMounts(c.mounts, false); err != nil {
		t.Fatal(err)
	}

	rawLocal, err = c.barrier.Get(coreLocalMountConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	if rawLocal == nil {
		t.Fatal("expected non-nil local mount")
	}
	localMountsTable = &MountTable{}
	if err := jsonutil.DecodeJSON(rawLocal.Value, localMountsTable); err != nil {
		t.Fatal(err)
	}
	// This requires some explanation: because we're directly munging the mount
	// table, the table initially when core unseals contains cubbyhole as per
	// above, but then we overwrite it with our own table with one local entry,
	// so we should now only expect the noop2 entry
	if len(localMountsTable.Entries) != 1 || localMountsTable.Entries[0].Path != "noop2/" {
		t.Fatalf("expected one entry in local mount table, got %#v", localMountsTable)
	}

	oldMounts := c.mounts
	if err := c.loadMounts(); err != nil {
		t.Fatal(err)
	}
	compEntries := c.mounts.Entries[:0]
	// Filter out required mounts
	for _, v := range c.mounts.Entries {
		if v.Type == "kv" {
			compEntries = append(compEntries, v)
		}
	}
	c.mounts.Entries = compEntries

	if !reflect.DeepEqual(oldMounts, c.mounts) {
		t.Fatalf("expected\n%#v\ngot\n%#v\n", oldMounts, c.mounts)
	}

	if len(c.mounts.Entries) != 2 {
		t.Fatalf("expected two mount entries, got %#v", localMountsTable)
	}
}

func TestCore_Unmount(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
	err := c.unmount("secret")
	if err != nil {
		t.Fatalf("err: %v", err)
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
	if !reflect.DeepEqual(c.mounts.sortEntriesByPath(), c2.mounts.sortEntriesByPath()) {
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
	view := c.router.MatchingStorageByAPIPath("test/")

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
	if err := c.unmount("test/"); err != nil {
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

	// View should be empty
	out, err := logical.CollectKeys(view)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("bad: %#v", out)
	}
}

func TestCore_Remount(t *testing.T) {
	c, keys, _ := TestCoreUnsealed(t)
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
	view := c.router.MatchingStorageByAPIPath("test/")

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
	out, err := logical.CollectKeys(view)
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
	c, _, _ := TestCoreUnsealed(t)
	table := c.defaultMountTable()
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

	// We filter out local entries here since the logic is rather dumb
	// (straight JSON comparison) and doesn't seal well with the separate
	// locations
	newEntries := mt.Entries[:0]
	for _, entry := range mt.Entries {
		if !entry.Local {
			newEntries = append(newEntries, entry)
		}
	}
	mt.Entries = newEntries

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

	var persistFunc func(*MountTable, bool) error

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
	if err := persistFunc(mt, false); err == nil {
		t.Fatal("expected error")
	}

	if len(mt.Entries) > 0 {
		mt.Type = origTableType
		mt.Entries[0].Table = "bar"
		if err := persistFunc(mt, false); err == nil {
			t.Fatal("expected error")
		}

		mt.Entries[0].Table = mt.Type
		if err := persistFunc(mt, false); err != nil {
			t.Fatal(err)
		}
	}
}

func verifyDefaultTable(t *testing.T, table *MountTable) {
	if len(table.Entries) != 4 {
		t.Fatalf("bad: %v", table.Entries)
	}
	table.sortEntriesByPath()
	for _, entry := range table.Entries {
		switch entry.Path {
		case "cubbyhole/":
			if entry.Type != "cubbyhole" {
				t.Fatalf("bad: %v", entry)
			}
		case "secret/":
			if entry.Type != "kv" {
				t.Fatalf("bad: %v", entry)
			}
		case "sys/":
			if entry.Type != "system" {
				t.Fatalf("bad: %v", entry)
			}
		case "identity/":
			if entry.Type != "identity" {
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

func TestSingletonMountTableFunc(t *testing.T) {
	c, _, _ := TestCoreUnsealed(t)

	mounts, auth := c.singletonMountTables()

	if len(mounts.Entries) != 2 {
		t.Fatalf("length of mounts is wrong; expected 2, got %d", len(mounts.Entries))
	}

	for _, entry := range mounts.Entries {
		switch entry.Type {
		case "system":
		case "identity":
		default:
			t.Fatalf("unknown type %s", entry.Type)
		}
	}

	if len(auth.Entries) != 1 {
		t.Fatal("length of auth is wrong")
	}

	if auth.Entries[0].Type != "token" {
		t.Fatal("unexpected entry type for auth")
	}
}
