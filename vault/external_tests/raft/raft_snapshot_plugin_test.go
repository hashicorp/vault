// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package rafttests_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// errParentRefs and errOutside identify a refusal because of
	// plugin_directory containment. errParentRefs comes from the registration
	// API's ".." check, errOutside from the checks that resolve symlinks, both
	// at registration and when a stored entry is looked up.
	errParentRefs = "path cannot contain parent references"
	errOutside    = "outside of configured plugin directory"

	// fixturePolicy exists only in the snapshot. Seeing it after a restore
	// proves the restore finished and its storage is being served.
	fixturePolicyName = "snapshot-fixture"
	fixturePolicy     = `path "snapshot-fixture" { capabilities = ["read"] }`
)

// TestRaftSnapshotPluginDirectory_Restore verifies that plugin catalog entries
// whose command escapes plugin_directory cannot be used to run a host binary
// after a raft snapshot restore of the same cluster.
//
// The plugin registration API refuses such commands (for example
// "../../bin/sh"), but a restore replaces the storage behind the catalog
// without running that validation. The entries are written with sys/raw before
// the snapshot is taken, so this also covers any other direct writer of
// catalog storage. Each entry escapes in a different way and is used through a
// different route that starts a plugin: a secrets mount, an auth mount and a
// database connection. The probe binary writes a per-entry marker file if it
// ever runs, so a failure identifies exactly which combination executed.
func TestRaftSnapshotPluginDirectory_Restore(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	pluginDir := newPluginDir(t)
	probe := newPluginProbe(t)
	escapes := newEscapes(t, ctx, pluginDir, probe)

	storage := newRaftStorage(t, "node")
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		Physical: storage, HAPhysical: storage,
		PluginDirectory: pluginDir, EnableRaw: true, DisableMlock: true,
	})
	client := cluster.Cores[0].Client

	requireRegistrationRejected(t, ctx, client, probe, escapes)
	snapshot := snapshotWithPlantedEscapes(t, ctx, client, probe, escapes)

	require.NoError(t, client.Sys().RaftSnapshotRestoreWithContext(ctx, bytes.NewReader(snapshot), false))
	waitForRestoredStorage(t, ctx, client)

	assertEscapesRefused(t, ctx, client, escapes)
	requireBuiltinMountsWork(t, ctx, client)
}

// TestRaftSnapshotPluginDirectory_ForceRestoreUnrelatedCluster verifies that
// plugin catalog entries whose command escapes plugin_directory cannot be used
// to run a host binary after a force restore of a snapshot taken from a
// different cluster.
//
// A force restore skips the seal-hash check that ties a snapshot to the
// cluster it came from, so whoever controls the source cluster controls the
// target's entire storage. With Shamir seals the target seals itself after the
// restore and can be unsealed with the source's key shares, which the builder
// of the snapshot knows. The entries, routes and markers are the same as in
// TestRaftSnapshotPluginDirectory_Restore.
func TestRaftSnapshotPluginDirectory_ForceRestoreUnrelatedCluster(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	pluginDir := newPluginDir(t)
	probe := newPluginProbe(t)
	escapes := newEscapes(t, ctx, pluginDir, probe)

	// EnableRaw is only needed on the source, to build the snapshot. The target
	// receives the entries through the restore API alone.
	sourceStorage := newRaftStorage(t, "source")
	source := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		Physical: sourceStorage, HAPhysical: sourceStorage,
		PluginDirectory: pluginDir, EnableRaw: true, DisableMlock: true,
	})
	targetStorage := newRaftStorage(t, "target")
	target := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		Physical: targetStorage, HAPhysical: targetStorage,
		PluginDirectory: pluginDir, DisableMlock: true,
	})
	client := target.Cores[0].Client

	requireRegistrationRejected(t, ctx, client, probe, escapes)
	snapshot := snapshotWithPlantedEscapes(t, ctx, source.Cores[0].Client, probe, escapes)

	// The snapshot contains the source's leader address. Stop the source so the
	// target can't forward to that unrelated live cluster while it becomes
	// active.
	source.StopCore(t, 0)

	// A normal restore must reject the unrelated cluster's seal keys. This is
	// the check that force skips.
	err := client.Sys().RaftSnapshotRestoreWithContext(ctx, bytes.NewReader(snapshot), false)
	require.ErrorContains(t, err, "could not verify hash file")

	require.NoError(t, client.Sys().RaftSnapshotRestoreWithContext(ctx, bytes.NewReader(snapshot), true))

	// The restored storage carries the source's keyring, so the target can't
	// decrypt it with its own keys and seals itself.
	require.Eventually(t, func() bool {
		status, err := client.Sys().SealStatusWithContext(ctx)
		return err == nil && status.Sealed
	}, 30*time.Second, 100*time.Millisecond, "force restore should seal after replacing the keyring")

	// Unseal with the source's key shares and switch to its root token. The
	// builder of the snapshot knows both.
	for _, key := range source.BarrierKeys {
		_, err := client.Sys().UnsealWithContext(ctx, hex.EncodeToString(key))
		require.NoError(t, err)
	}
	client.SetToken(source.RootToken)
	waitForRestoredStorage(t, ctx, client)

	assertEscapesRefused(t, ctx, client, escapes)
	requireBuiltinMountsWork(t, ctx, client)
}

// TestRaftSnapshotPluginDirectory_RestoredMountTable verifies that a snapshot
// whose mount tables already use a catalog entry that escapes plugin_directory
// does not run that binary when the restore brings the mounts up, even though
// no request ever mounts the plugin.
//
// Setting up mounts is part of unsealing, which a restore performs itself, so
// a snapshot that references the plugin from a mount would otherwise run it
// straight away and again on every later unseal. The test covers a secrets
// mount, an auth mount, and a mount of the builtin "transit" engine whose name
// is shadowed by an escaping catalog entry.
//
// A refused mount must be skipped rather than failing the unseal, which would
// take the node down: the node must unseal, keep the mounts and their data,
// and refuse requests to them.
func TestRaftSnapshotPluginDirectory_RestoredMountTable(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	pluginDir := newPluginDir(t)
	probe := newPluginProbe(t)
	command := probe.relativeTo(t, pluginDir)

	storage := newRaftStorage(t, "node")
	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		Physical: storage, HAPhysical: storage,
		PluginDirectory: pluginDir, EnableRaw: true, DisableMlock: true,
	})
	client := cluster.Cores[0].Client

	type restoredMount struct {
		// path is the mount path as stored in the mount table, which has no
		// "auth/" prefix for auth mounts.
		path string
		// table is the storage key of the mount table holding the mount.
		table string
		// builtinType creates the mount on the source, so that it has real
		// storage, a UUID and an accessor.
		builtinType string
		// pluginName is the catalog entry the mount uses in the snapshot.
		pluginName string
		pluginType consts.PluginType
		marker     string
		uuid       string
	}
	mounts := []*restoredMount{
		{path: "restored-secret/", table: "core/mounts", builtinType: "kv", pluginName: "mounted-secret", pluginType: consts.PluginTypeSecrets},
		{path: "restored-auth/", table: "core/auth", builtinType: "userpass", pluginName: "mounted-auth", pluginType: consts.PluginTypeCredential},
		{path: "shadowed-builtin/", table: "core/mounts", builtinType: "transit", pluginName: "transit", pluginType: consts.PluginTypeSecrets},
	}

	// --- Build the snapshot -----------------------------------------------

	for _, m := range mounts {
		m.marker = probe.newMarker(t, ctx)
		switch m.pluginType {
		case consts.PluginTypeSecrets:
			require.NoError(t, client.Sys().MountWithContext(ctx, m.path, &api.MountInput{Type: m.builtinType}))
		case consts.PluginTypeCredential:
			require.NoError(t, client.Sys().EnableAuthWithOptionsWithContext(ctx, m.path, &api.EnableAuthOptions{Type: m.builtinType}))
		default:
			t.Fatalf("%s: unhandled plugin type %s", m.path, m.pluginType)
		}
		plantCatalogEntry(t, ctx, client, m.pluginName, m.pluginType, command, probe, m.marker)
	}

	// Data that must survive the mount being skipped.
	_, err := client.Logical().WriteWithContext(ctx, "restored-secret/probe", map[string]interface{}{"key": "value"})
	require.NoError(t, err)

	// Point each mount at its escaping catalog entry in the stored mount
	// tables. The tables in memory are only reloaded by the restore, so
	// nothing runs yet.
	for _, table := range []string{"core/mounts", "core/auth"} {
		rawPath := "sys/raw/" + table
		secret, err := client.Logical().ReadWithContext(ctx, rawPath)
		require.NoError(t, err)
		require.NotNil(t, secret)

		// UseNumber keeps durations in nanoseconds exact through the round trip.
		decoder := json.NewDecoder(strings.NewReader(secret.Data["value"].(string)))
		decoder.UseNumber()
		var decoded map[string]interface{}
		require.NoError(t, decoder.Decode(&decoded))
		entries, ok := decoded["entries"].([]interface{})
		require.True(t, ok, table)

		for _, raw := range entries {
			entry := raw.(map[string]interface{})
			for _, m := range mounts {
				if m.table != table || entry["path"] != m.path {
					continue
				}
				entry["type"] = m.pluginName
				// Without a version the lookup reads the unversioned catalog
				// entry, as it does for any external plugin. A builtin version
				// would resolve to the builtin and never reach the entry.
				delete(entry, "plugin_version")
				delete(entry, "running_plugin_version")
				delete(entry, "running_sha256")
				m.uuid = entry["uuid"].(string)
			}
		}

		encoded, err := json.Marshal(decoded)
		require.NoError(t, err)
		_, err = client.Logical().WriteWithContext(ctx, rawPath, map[string]interface{}{"value": string(encoded)})
		require.NoError(t, err)
	}
	for _, m := range mounts {
		require.NotEmpty(t, m.uuid, "%s: mount not found in %s", m.path, m.table)
	}

	require.NoError(t, client.Sys().PutPolicyWithContext(ctx, fixturePolicyName, fixturePolicy))
	var snapshot bytes.Buffer
	require.NoError(t, client.Sys().RaftSnapshotWithContext(ctx, &snapshot))
	require.NoError(t, client.Sys().DeletePolicyWithContext(ctx, fixturePolicyName))

	// --- Restore: this alone sets up the mounts ---------------------------

	require.NoError(t, client.Sys().RaftSnapshotRestoreWithContext(ctx, bytes.NewReader(snapshot.Bytes()), false))
	waitForRestoredStorage(t, ctx, client)

	for _, m := range mounts {
		// The security assertion: setting up the restored mount must not have
		// run the binary outside plugin_directory. assert reports every mount
		// that ran it in one pass.
		assert.NoFileExists(t, m.marker, "%s: restored mount table executed a binary outside plugin_directory", m.path)

		// Mount setup only logs why it skipped a mount. The catalog read API
		// uses the same lookup and shows the refusal is because of the path.
		apiType, err := api.ParsePluginType(m.pluginType.String())
		require.NoError(t, err)
		_, err = client.Sys().GetPluginWithContext(ctx, &api.GetPluginInput{Name: m.pluginName, Type: apiType})
		assert.ErrorContains(t, err, errOutside, m.path)
	}

	// --- The refused mounts are skipped, not dropped ----------------------

	status, err := client.Sys().SealStatusWithContext(ctx)
	require.NoError(t, err)
	require.False(t, status.Sealed)

	// The mounts stay in the mount tables with the restored plugin name...
	secretMounts, err := client.Sys().ListMountsWithContext(ctx)
	require.NoError(t, err)
	authMounts, err := client.Sys().ListAuthWithContext(ctx)
	require.NoError(t, err)
	for _, m := range mounts {
		listed := secretMounts
		if m.pluginType == consts.PluginTypeCredential {
			listed = authMounts
		}
		if assert.Contains(t, listed, m.path) {
			assert.Equal(t, m.pluginName, listed[m.path].Type, m.path)
		}
	}

	// ...but have no backend, so requests fail instead of starting the plugin.
	_, err = client.Logical().WriteWithContext(ctx, "restored-secret/probe", map[string]interface{}{"key": "other"})
	require.ErrorContains(t, err, "backend is nil")

	// Their data is kept, so the mount works again once the entry is fixed.
	stored, err := client.Logical().ReadWithContext(ctx, "sys/raw/logical/"+mounts[0].uuid+"/probe")
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.JSONEq(t, `{"key":"value"}`, stored.Data["value"].(string))

	requireBuiltinMountsWork(t, ctx, client)
}

// TestRaftSnapshotPluginProcess is a subprocess fixture that records execution
// without running a shell or touching anything outside the parent test's tempdir.
//
// It does nothing when run as part of the normal test suite, because the marker
// environment variable is unset. The TestRaftSnapshotPluginDirectory tests run
// the test binary with -test.run selecting only this test and the variable set,
// so the presence of the marker file is proof that the binary was executed.
func TestRaftSnapshotPluginProcess(t *testing.T) {
	t.Parallel()
	marker := os.Getenv("VAULT_SNAPSHOT_PLUGIN_TEST_MARKER")
	if marker == "" {
		return
	}
	require.NoError(t, os.WriteFile(marker, []byte("executed\n"), 0o600))
}

// pluginProbe stands in for an attacker-chosen host binary. It is the test
// binary itself, which lives outside every plugin directory, re-run so that
// only TestRaftSnapshotPluginProcess executes and writes a marker file.
type pluginProbe struct {
	executable string
	// checksum is the probe's real sha256. go-plugin compares the catalog's
	// sha256 with the binary before running it, so using the real value rules
	// out a checksum mismatch as the reason the probe did not run.
	checksum []byte
	args     []string
}

func newPluginProbe(t *testing.T) *pluginProbe {
	t.Helper()

	executable, err := os.Executable()
	require.NoError(t, err)
	executable, err = filepath.EvalSymlinks(executable)
	require.NoError(t, err)

	binary, err := os.Open(executable)
	require.NoError(t, err)
	hash := sha256.New()
	_, err = io.Copy(hash, binary)
	require.NoError(t, binary.Close())
	require.NoError(t, err)

	return &pluginProbe{
		executable: executable,
		checksum:   hash.Sum(nil),
		args:       []string{"-test.run=^TestRaftSnapshotPluginProcess$"},
	}
}

// env returns the plugin environment that makes the probe write marker.
func (p *pluginProbe) env(marker string) []string {
	return []string{"VAULT_SNAPSHOT_PLUGIN_TEST_MARKER=" + marker}
}

// newMarker returns a fresh marker path after proving that running the probe
// with it creates the file, so that a later absence of the file means the
// probe did not run.
func (p *pluginProbe) newMarker(t *testing.T, ctx context.Context) string {
	t.Helper()

	marker := filepath.Join(t.TempDir(), "executed")
	cmd := exec.CommandContext(ctx, p.executable, p.args...)
	cmd.Env = append(os.Environ(), p.env(marker)...)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "%s", output)
	require.FileExists(t, marker)
	require.NoError(t, os.Remove(marker))
	return marker
}

// relativeTo returns the probe's path relative to pluginDir. It starts with
// "..", which is exactly what plugin_directory is meant to prevent.
func (p *pluginProbe) relativeTo(t *testing.T, pluginDir string) string {
	t.Helper()

	relative, err := filepath.Rel(pluginDir, p.executable)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(relative, ".."+string(filepath.Separator)))
	return relative
}

// escape is a catalog entry whose command reaches the probe outside
// plugin_directory. It is used through the route named by its plugin type.
type escape struct {
	name       string
	pluginType consts.PluginType
	command    string
	// registrationErr is what the registration API says about this command,
	// proving the API would never have accepted it.
	registrationErr string
	marker          string
}

// newEscapes returns one entry for each way of escaping pluginDir, each with
// its own marker.
func newEscapes(t *testing.T, ctx context.Context, pluginDir string, probe *pluginProbe) []*escape {
	t.Helper()

	relative := probe.relativeTo(t, pluginDir)

	// More ".." segments than the plugin directory is deep. Joining this with
	// pluginDir is cleaned to the probe's absolute path, so an attacker does
	// not need to know how deep the directory is.
	excessive := strings.Repeat(".."+string(filepath.Separator), 64) +
		strings.TrimPrefix(probe.executable, string(filepath.Separator))

	// A command with no ".." at all: a symlink inside pluginDir that points at
	// the probe outside it. Only a check on the resolved path catches this.
	const symlinkName = "symlink-escape"
	require.NoError(t, os.Symlink(probe.executable, filepath.Join(pluginDir, symlinkName)))

	escapes := []*escape{
		{name: "secret-relative", pluginType: consts.PluginTypeSecrets, command: relative, registrationErr: errParentRefs},
		{name: "auth-relative", pluginType: consts.PluginTypeCredential, command: relative, registrationErr: errParentRefs},
		{name: "database-relative", pluginType: consts.PluginTypeDatabase, command: relative, registrationErr: errParentRefs},
		{name: "secret-excessive-parents", pluginType: consts.PluginTypeSecrets, command: excessive, registrationErr: errParentRefs},
		{name: "secret-symlink", pluginType: consts.PluginTypeSecrets, command: symlinkName, registrationErr: errOutside},
	}
	for _, e := range escapes {
		e.marker = probe.newMarker(t, ctx)
	}
	return escapes
}

// newPluginDir returns an empty plugin directory. Symlinks are resolved
// because temp paths can traverse one (e.g. /var on macOS), which would
// otherwise make later path comparisons misleading.
func newPluginDir(t *testing.T) string {
	t.Helper()

	pluginDir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return pluginDir
}

// newRaftStorage returns a raft backend that is closed when the test ends.
// Create it before the cluster that uses it, so the cluster is cleaned up
// first.
func newRaftStorage(t *testing.T, nodeID string) *raft.RaftBackend {
	t.Helper()

	storage, err := raft.NewRaftBackend(map[string]string{
		"path": t.TempDir(), "node_id": nodeID, "performance_multiplier": "1",
	}, corehelpers.NewTestLogger(t))
	require.NoError(t, err)
	raftStorage := storage.(*raft.RaftBackend)
	t.Cleanup(func() { require.NoError(t, raftStorage.Close()) })
	return raftStorage
}

// requireRegistrationRejected proves that the registration API refuses every
// escape, so the entries can only reach the catalog by bypassing it.
func requireRegistrationRejected(t *testing.T, ctx context.Context, client *api.Client, probe *pluginProbe, escapes []*escape) {
	t.Helper()

	for _, e := range escapes {
		apiType, err := api.ParsePluginType(e.pluginType.String())
		require.NoError(t, err)
		err = client.Sys().RegisterPluginWithContext(ctx, &api.RegisterPluginInput{
			Name: e.name, Type: apiType,
			Command: e.command, Args: probe.args, Env: probe.env(e.marker), SHA256: hex.EncodeToString(probe.checksum),
		})
		require.ErrorContains(t, err, e.registrationErr, e.name)
		require.NoFileExists(t, e.marker, e.name)
	}
}

// plantCatalogEntry writes a catalog entry straight to storage with sys/raw,
// bypassing the registration checks, and returns its raw path.
func plantCatalogEntry(t *testing.T, ctx context.Context, client *api.Client, name string, pluginType consts.PluginType, command string, probe *pluginProbe, marker string) string {
	t.Helper()

	entry, err := json.Marshal(&pluginutil.PluginRunner{
		Name: name, Type: pluginType,
		Command: command, Args: probe.args, Env: probe.env(marker), Sha256: probe.checksum,
	})
	require.NoError(t, err)
	rawPath := "sys/raw/core/plugin-catalog/" + pluginType.String() + "/" + name
	_, err = client.Logical().WriteWithContext(ctx, rawPath, map[string]interface{}{"value": string(entry)})
	require.NoError(t, err, name)
	return rawPath
}

// snapshotWithPlantedEscapes plants the escapes and the fixture policy in the
// cluster behind client, snapshots it, and removes them again. When the
// snapshot is restored into the same cluster, the restore is then the only
// place they can come from.
func snapshotWithPlantedEscapes(t *testing.T, ctx context.Context, client *api.Client, probe *pluginProbe, escapes []*escape) []byte {
	t.Helper()

	rawPaths := make([]string, 0, len(escapes))
	for _, e := range escapes {
		rawPaths = append(rawPaths, plantCatalogEntry(t, ctx, client, e.name, e.pluginType, e.command, probe, e.marker))
	}
	require.NoError(t, client.Sys().PutPolicyWithContext(ctx, fixturePolicyName, fixturePolicy))

	var snapshot bytes.Buffer
	require.NoError(t, client.Sys().RaftSnapshotWithContext(ctx, &snapshot))

	for _, rawPath := range rawPaths {
		_, err := client.Logical().DeleteWithContext(ctx, rawPath)
		require.NoError(t, err)
	}
	require.NoError(t, client.Sys().DeletePolicyWithContext(ctx, fixturePolicyName))
	return snapshot.Bytes()
}

// waitForRestoredStorage waits until the node is active and serves the
// fixture policy, which exists only in the snapshot. The restore API responds
// before the restore has completed.
func waitForRestoredStorage(t *testing.T, ctx context.Context, client *api.Client) {
	t.Helper()

	require.Eventually(t, func() bool {
		leader, err := client.Sys().LeaderWithContext(ctx)
		if err != nil || !leader.IsSelf {
			return false
		}
		restored, err := client.Sys().GetPolicyWithContext(ctx, fixturePolicyName)
		return err == nil && restored == fixturePolicy
	}, 30*time.Second, 100*time.Millisecond)
}

// assertEscapesRefused uses every escape through its route and asserts that
// the probe did not run and that the request was refused because of path
// containment.
func assertEscapesRefused(t *testing.T, ctx context.Context, client *api.Client, escapes []*escape) {
	t.Helper()

	// A builtin database mount that hosts the connection the database escape
	// configures. It is not part of the exploit.
	require.NoError(t, client.Sys().MountWithContext(ctx, "database-probe", &api.MountInput{Type: "database"}))

	for _, e := range escapes {
		// Each route looks the plugin up in the restored catalog and starts it.
		var err error
		switch e.pluginType {
		case consts.PluginTypeSecrets:
			err = client.Sys().MountWithContext(ctx, e.name, &api.MountInput{Type: e.name})
		case consts.PluginTypeCredential:
			err = client.Sys().EnableAuthWithOptionsWithContext(ctx, e.name, &api.EnableAuthOptions{Type: e.name})
		case consts.PluginTypeDatabase:
			_, err = client.Logical().WriteWithContext(ctx, "database-probe/config/"+e.name, map[string]interface{}{
				"plugin_name":       e.name,
				"connection_url":    "unused",
				"allowed_roles":     "*",
				"verify_connection": false,
			})
		default:
			t.Fatalf("%s: unhandled plugin type %s", e.name, e.pluginType)
		}

		// The security assertion: the binary outside plugin_directory must not
		// have run. assert reports every bypassing route in one pass.
		assert.NoFileExists(t, e.marker, "%s: restored catalog executed a binary outside plugin_directory; error: %v", e.name, err)

		// The probe does not implement the plugin protocol, so a request error
		// alone proves nothing. The marker check above does.
		if !assert.Error(t, err, e.name) {
			continue
		}

		// The request must have been refused because of the path, not because
		// of some unrelated startup failure that happens to prevent execution.
		var responseError *api.ResponseError
		if assert.ErrorAs(t, err, &responseError, e.name) {
			assert.Regexp(t, errOutside+"|"+errParentRefs, strings.Join(responseError.Errors, "\n"),
				"%s: must fail because of path containment, not an unrelated startup error", e.name)
		}
	}
}

// requireBuiltinMountsWork asserts that the node is unsealed and still serves
// builtin mounts, which do not depend on the catalog's external entries, so a
// bad entry is contained rather than breaking the node.
func requireBuiltinMountsWork(t *testing.T, ctx context.Context, client *api.Client) {
	t.Helper()

	status, err := client.Sys().SealStatusWithContext(ctx)
	require.NoError(t, err)
	require.False(t, status.Sealed)

	require.NoError(t, client.Sys().MountWithContext(ctx, "sibling", &api.MountInput{Type: "kv"}))
	_, err = client.Logical().WriteWithContext(ctx, "sibling/probe", map[string]interface{}{"key": "value"})
	require.NoError(t, err)
	secret, err := client.Logical().ReadWithContext(ctx, "sibling/probe")
	require.NoError(t, err)
	require.NotNil(t, secret)
	require.Equal(t, "value", secret.Data["key"])
}
