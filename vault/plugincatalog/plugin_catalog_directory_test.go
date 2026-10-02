// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/stretchr/testify/require"
)

// TestPluginCatalog_Get_CommandContainment verifies that PluginCatalog.get only
// hands out a runnable command for entries whose command stays inside the
// configured plugin directory, even when the entry was written straight to
// storage (as happens with a raft snapshot restore) and so never went through
// the registration checks in Set. Entries that escape the directory, lexically
// or through symlinks, must be refused, while legitimate layouts and entries
// whose binary is simply missing must keep working.
func TestPluginCatalog_Get_CommandContainment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// setup prepares the filesystem under root, where pluginDir is the
		// configured plugin directory (root/plugins, created and empty). It
		// returns the command to store in the catalog entry.
		setup func(t *testing.T, root, pluginDir string) string
		// wantCommand is the expected returned command relative to the plugin
		// directory, when the lookup should succeed.
		wantCommand string
		wantErr     bool
	}{
		{
			name: "binary directly in the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my-plugin"), []byte("x"), 0o755))
				return "my-plugin"
			},
			wantCommand: "my-plugin",
		},
		{
			name: "binary in an extracted artifact subdirectory",
			setup: func(t *testing.T, root, pluginDir string) string {
				sub := filepath.Join(pluginDir, "my-plugin_1.0.0_linux_amd64")
				require.NoError(t, os.Mkdir(sub, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(sub, "my-plugin"), []byte("x"), 0o755))
				return filepath.Join("my-plugin_1.0.0_linux_amd64", "my-plugin")
			},
			wantCommand: filepath.Join("my-plugin_1.0.0_linux_amd64", "my-plugin"),
		},
		{
			name: "missing binary is not an error",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "gone"
			},
			wantCommand: "gone",
		},
		{
			name: "missing directory component is not an error",
			setup: func(t *testing.T, root, pluginDir string) string {
				return filepath.Join("no-such-dir", "my-plugin")
			},
			wantCommand: filepath.Join("no-such-dir", "my-plugin"),
		},
		{
			name: "absolute looking command is confined by the join",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "/bin/true"
			},
			wantCommand: filepath.Join("bin", "true"),
		},
		{
			name: "dot dot segments that stay inside the directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my-plugin"), []byte("x"), 0o755))
				return filepath.Join("sub", "..", "my-plugin")
			},
			wantCommand: "my-plugin",
		},
		{
			name: "symlink to another file inside the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "real"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "real"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantCommand: "link",
		},
		{
			name: "relative parent reference escaping the directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(root, "outside"), []byte("x"), 0o755))
				return filepath.Join("..", "outside")
			},
			wantErr: true,
		},
		{
			name: "parent reference escaping to a missing target",
			setup: func(t *testing.T, root, pluginDir string) string {
				return filepath.Join("..", "outside")
			},
			wantErr: true,
		},
		{
			name: "more parent references than the directory is deep",
			setup: func(t *testing.T, root, pluginDir string) string {
				target := filepath.Join(root, "outside")
				require.NoError(t, os.WriteFile(target, []byte("x"), 0o755))
				return strings.Repeat("../", 64) + strings.TrimPrefix(target, "/")
			},
			wantErr: true,
		},
		{
			name: "sibling directory sharing the plugin directory name as a prefix",
			setup: func(t *testing.T, root, pluginDir string) string {
				evil := pluginDir + "-evil"
				require.NoError(t, os.Mkdir(evil, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(evil, "my-plugin"), []byte("x"), 0o755))
				return filepath.Join("..", filepath.Base(evil), "my-plugin")
			},
			wantErr: true,
		},
		{
			name: "symlink to a sibling directory sharing the plugin directory name as a prefix",
			setup: func(t *testing.T, root, pluginDir string) string {
				evil := pluginDir + "-evil"
				require.NoError(t, os.Mkdir(evil, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(evil, "my-plugin"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(filepath.Join(evil, "my-plugin"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "symlink pointing outside of the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				target := filepath.Join(root, "outside")
				require.NoError(t, os.WriteFile(target, []byte("x"), 0o755))
				require.NoError(t, os.Symlink(target, filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "symlinked subdirectory pointing outside of the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				outside := filepath.Join(root, "outside-dir")
				require.NoError(t, os.Mkdir(outside, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(outside, "my-plugin"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(outside, filepath.Join(pluginDir, "sub")))
				return filepath.Join("sub", "my-plugin")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			root, err := filepath.EvalSymlinks(t.TempDir())
			require.NoError(t, err)
			pluginDir := filepath.Join(root, "plugins")
			require.NoError(t, os.Mkdir(pluginDir, 0o755))

			command := tt.setup(t, root, pluginDir)
			catalog := newDirectoryTestCatalog(t, pluginDir)
			putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
				Name:    "test-plugin",
				Type:    consts.PluginTypeSecrets,
				Command: command,
				Sha256:  []byte("sha"),
			})

			got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
			if tt.wantErr {
				require.ErrorIs(t, err, ErrPluginOutsideDirectory)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, filepath.Join(pluginDir, tt.wantCommand), got.Command)
		})
	}
}

// TestPluginCatalog_Get_SymlinkedPluginDirectory verifies that a plugin
// directory that is itself reached through a symlink (common on macOS and in
// container images) does not cause legitimate plugin binaries to be rejected,
// because both the directory and the command are resolved before comparing.
func TestPluginCatalog_Get_SymlinkedPluginDirectory(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	realDir := filepath.Join(root, "real-plugins")
	require.NoError(t, os.Mkdir(realDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(realDir, "my-plugin"), []byte("x"), 0o755))
	linkDir := filepath.Join(root, "link-plugins")
	require.NoError(t, os.Symlink(realDir, linkDir))

	catalog := newDirectoryTestCatalog(t, linkDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "my-plugin",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, filepath.Join(linkDir, "my-plugin"), got.Command)
}

// TestPluginCatalog_Get_NoPluginDirectory verifies that when no plugin
// directory is configured a stored non-container entry is never returned as an
// external plugin, regardless of its command.
func TestPluginCatalog_Get_NoPluginDirectory(t *testing.T) {
	t.Parallel()

	catalog := newDirectoryTestCatalog(t, "")
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "../../bin/sh",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.NoError(t, err)
	require.Nil(t, got)
}

// TestPluginCatalog_Get_OCIEntryUnchanged verifies that container plugin
// entries are not subject to the plugin directory containment check, since
// they do not reference a host binary.
func TestPluginCatalog_Get_OCIEntryUnchanged(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	pluginDir := filepath.Join(root, "plugins")
	require.NoError(t, os.Mkdir(pluginDir, 0o755))

	catalog := newDirectoryTestCatalog(t, pluginDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:     "test-plugin",
		Type:     consts.PluginTypeSecrets,
		OCIImage: "example.com/some/image",
		Command:  "../not-used-for-containers",
		Sha256:   []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "example.com/some/image", got.OCIImage)
	require.Equal(t, "../not-used-for-containers", got.Command)
}

// TestPluginCatalog_Get_UnsafeEntryShadowingBuiltin verifies that an unsafe
// stored entry that reuses the name of a builtin plugin yields an error rather
// than being executed or silently returned. An error (rather than nil) is
// required so that mount setup skips just that mount instead of failing unseal.
func TestPluginCatalog_Get_UnsafeEntryShadowingBuiltin(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	pluginDir := filepath.Join(root, "plugins")
	require.NoError(t, os.Mkdir(pluginDir, 0o755))

	catalog := newDirectoryTestCatalog(t, pluginDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "consul",
		Type:    consts.PluginTypeSecrets,
		Command: "../../bin/sh",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "consul", consts.PluginTypeSecrets, "")
	require.ErrorIs(t, err, ErrPluginOutsideDirectory)
	require.Nil(t, got)
}

// TestPluginCatalog_Get_ExistingRegistrationsUnaffected verifies that entries
// registered through Set, which is where the containment rules originate, are
// still returned by get after the read-time check was added, so that catalogs
// populated before the check existed keep working.
func TestPluginCatalog_Get_ExistingRegistrationsUnaffected(t *testing.T) {
	t.Parallel()

	catalog := testPluginCatalog(t)
	file, err := os.CreateTemp(catalog.directory, "registered")
	require.NoError(t, err)
	require.NoError(t, file.Close())
	require.NoError(t, os.Chmod(file.Name(), 0o755))

	command := filepath.Base(file.Name())
	require.NoError(t, catalog.Set(t.Context(), pluginutil.SetPluginInput{
		Name:    "registered",
		Type:    consts.PluginTypeDatabase,
		Command: command,
		Sha256:  []byte{'1'},
	}))

	got, err := catalog.Get(t.Context(), "registered", consts.PluginTypeDatabase, "")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, filepath.Join(catalog.directory, command), got.Command)
}

// newDirectoryTestCatalog builds a catalog backed by in-memory storage with the
// given plugin directory.
func newDirectoryTestCatalog(t *testing.T, pluginDir string) *PluginCatalog {
	t.Helper()

	logger := hclog.NewNullLogger()
	storage, err := inmem.NewInmem(nil, logger)
	require.NoError(t, err)

	catalog, err := SetupPluginCatalog(t.Context(), &PluginCatalogInput{
		Logger:               logger,
		BuiltinRegistry:      corehelpers.NewMockBuiltinRegistry(),
		CatalogView:          logical.NewLogicalStorage(storage),
		PluginDirectory:      pluginDir,
		PluginRuntimeCatalog: testPluginRuntimeCatalog(t),
	})
	require.NoError(t, err)
	return catalog
}

// putRawCatalogEntry writes a plugin entry directly to the catalog's storage,
// bypassing Set and its validation, the same way a snapshot restore does.
func putRawCatalogEntry(t *testing.T, catalog *PluginCatalog, entry *pluginutil.PluginRunner) {
	t.Helper()

	buf, err := json.Marshal(entry)
	require.NoError(t, err)
	err = catalog.catalogView.Put(t.Context(), &logical.StorageEntry{
		Key:   entry.Type.String() + "/" + entry.Name,
		Value: buf,
	})
	require.NoError(t, err)
}
