// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"encoding/json"
	"os"
	"path"
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
// hands out a runnable command for entries that satisfy the same command rules
// registration (Set) enforces, even when the entry was written straight to
// storage (as happens with a raft snapshot restore) and so never went through
// registration. The binary must be directly in the directory the entry's
// layout implies: the plugin directory, the extracted artifact directory, or
// the runtime directory of a downloaded plugin. Entries whose binary does not
// exist yet are allowed, as long as nothing on the existing part of the path
// is a symlink.
func TestPluginCatalog_Get_CommandContainment(t *testing.T) {
	t.Parallel()

	artifactDir := GetExtractedArtifactDir("test-plugin", "1.0.0")

	tests := []struct {
		name string
		// setup prepares the filesystem under root, where pluginDir is the
		// configured plugin directory (root/plugins, created and empty). It
		// returns the command to store in the catalog entry.
		setup    func(t *testing.T, root, pluginDir string) string
		version  string
		download bool
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
			name: "missing binary directly in the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "gone"
			},
			wantCommand: "gone",
		},
		{
			name: "leading current directory component",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my-plugin"), []byte("x"), 0o755))
				return "./my-plugin"
			},
			wantCommand: "my-plugin",
		},
		{
			name: "symlink to another file directly in the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "real"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "real"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantCommand: "link",
		},
		{
			name: "binary in its extracted artifact directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Mkdir(filepath.Join(pluginDir, artifactDir), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, artifactDir, "test-plugin"), []byte("x"), 0o755))
				return filepath.Join(artifactDir, "test-plugin")
			},
			version:     "1.0.0",
			wantCommand: filepath.Join(artifactDir, "test-plugin"),
		},
		{
			name: "missing extracted artifact directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				return filepath.Join(artifactDir, "test-plugin")
			},
			version:     "1.0.0",
			wantCommand: filepath.Join(artifactDir, "test-plugin"),
		},
		{
			name: "downloaded plugin in its runtime directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				runtimeDir := filepath.Join(pluginDir, ".runtime", artifactDir)
				require.NoError(t, os.MkdirAll(runtimeDir, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(runtimeDir, "test-plugin"), []byte("x"), 0o755))
				return filepath.Join(".runtime", artifactDir, "test-plugin")
			},
			version:     "1.0.0",
			download:    true,
			wantCommand: filepath.Join(".runtime", artifactDir, "test-plugin"),
		},
		{
			name: "downloaded plugin not extracted yet",
			setup: func(t *testing.T, root, pluginDir string) string {
				return filepath.Join(".runtime", artifactDir, "test-plugin")
			},
			version:     "1.0.0",
			download:    true,
			wantCommand: filepath.Join(".runtime", artifactDir, "test-plugin"),
		},
		{
			name: "empty command names the plugin directory itself",
			setup: func(t *testing.T, root, pluginDir string) string {
				return ""
			},
			wantErr: true,
		},
		{
			name: "dot command names the plugin directory itself",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "."
			},
			wantErr: true,
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
			name: "parent references that would stay inside the directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my-plugin"), []byte("x"), 0o755))
				return "sub/../my-plugin"
			},
			wantErr: true,
		},
		{
			name: "double dot inside a file name",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my..plugin"), []byte("x"), 0o755))
				return "my..plugin"
			},
			wantErr: true,
		},
		{
			name: "absolute looking command lands in a subdirectory",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "/bin/true"
			},
			wantErr: true,
		},
		{
			name: "binary in a subdirectory that is not an artifact directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Mkdir(filepath.Join(pluginDir, "sub"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "sub", "my-plugin"), []byte("x"), 0o755))
				return filepath.Join("sub", "my-plugin")
			},
			wantErr: true,
		},
		{
			name: "artifact directory of a different version",
			setup: func(t *testing.T, root, pluginDir string) string {
				other := GetExtractedArtifactDir("test-plugin", "2.0.0")
				require.NoError(t, os.Mkdir(filepath.Join(pluginDir, other), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, other, "test-plugin"), []byte("x"), 0o755))
				return filepath.Join(other, "test-plugin")
			},
			version: "1.0.0",
			wantErr: true,
		},
		{
			name: "downloaded plugin outside its runtime directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "test-plugin"), []byte("x"), 0o755))
				return "test-plugin"
			},
			version:  "1.0.0",
			download: true,
			wantErr:  true,
		},
		{
			name: "version with parent references",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "test-plugin"
			},
			version:  "1/../../../../../../outside",
			download: true,
			wantErr:  true,
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
			name: "symlink to a file in a subdirectory of the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Mkdir(filepath.Join(pluginDir, "sub"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "sub", "real"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "sub", "real"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "symlink chain that leaves the plugin directory on its last hop",
			setup: func(t *testing.T, root, pluginDir string) string {
				target := filepath.Join(root, "outside")
				require.NoError(t, os.WriteFile(target, []byte("x"), 0o755))
				require.NoError(t, os.Symlink(target, filepath.Join(pluginDir, "hop")))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "hop"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "broken symlink pointing outside of the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Symlink(filepath.Join(root, "missing"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "broken symlink pointing inside the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "missing"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "broken symlink chain",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Symlink(filepath.Join(root, "missing"), filepath.Join(pluginDir, "hop")))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "hop"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "artifact directory symlinked outside of the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				outside := filepath.Join(root, "outside-dir")
				require.NoError(t, os.Mkdir(outside, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(outside, "test-plugin"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(outside, filepath.Join(pluginDir, artifactDir)))
				return filepath.Join(artifactDir, "test-plugin")
			},
			version: "1.0.0",
			wantErr: true,
		},
		{
			name: "artifact directory symlinked outside with the binary missing",
			setup: func(t *testing.T, root, pluginDir string) string {
				outside := filepath.Join(root, "outside-dir")
				require.NoError(t, os.Mkdir(outside, 0o755))
				require.NoError(t, os.Symlink(outside, filepath.Join(pluginDir, artifactDir)))
				return filepath.Join(artifactDir, "test-plugin")
			},
			version: "1.0.0",
			wantErr: true,
		},
		{
			name: "artifact directory symlinked to another directory inside the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				other := filepath.Join(pluginDir, "other")
				require.NoError(t, os.Mkdir(other, 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(other, "test-plugin"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(other, filepath.Join(pluginDir, artifactDir)))
				return filepath.Join(artifactDir, "test-plugin")
			},
			version: "1.0.0",
			wantErr: true,
		},
		{
			name: "broken artifact directory symlink",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Symlink(filepath.Join(root, "missing-dir"), filepath.Join(pluginDir, artifactDir)))
				return filepath.Join(artifactDir, "test-plugin")
			},
			version: "1.0.0",
			wantErr: true,
		},
		{
			name: "runtime directory symlinked outside before extraction",
			setup: func(t *testing.T, root, pluginDir string) string {
				outside := filepath.Join(root, "outside-runtime")
				require.NoError(t, os.Mkdir(outside, 0o755))
				require.NoError(t, os.Symlink(outside, filepath.Join(pluginDir, ".runtime")))
				return filepath.Join(".runtime", artifactDir, "test-plugin")
			},
			version:  "1.0.0",
			download: true,
			wantErr:  true,
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
				Name:     "test-plugin",
				Type:     consts.PluginTypeSecrets,
				Version:  tt.version,
				Command:  command,
				Download: tt.download,
				Sha256:   []byte("sha"),
			})

			got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, tt.version)
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

// TestPluginCatalog_SetAndGetAgree verifies that registration (Set) and lookup
// (get) accept and refuse the same commands when the binary exists. The same
// command is registered through Set in one catalog and written straight to
// storage in another, the way a snapshot restore would, and both must agree.
// A disagreement means a stored entry could do something registration forbids.
func TestPluginCatalog_SetAndGetAgree(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		// setup prepares pluginDir under root and returns the command.
		setup   func(t *testing.T, root, pluginDir string) string
		wantErr bool
	}{
		{
			name: "binary directly in the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my-plugin"), []byte("x"), 0o755))
				return "my-plugin"
			},
		},
		{
			name: "symlink to another file directly in the plugin directory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "real"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "real"), filepath.Join(pluginDir, "link")))
				return "link"
			},
		},
		{
			name: "binary in a subdirectory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Mkdir(filepath.Join(pluginDir, "sub"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "sub", "my-plugin"), []byte("x"), 0o755))
				return filepath.Join("sub", "my-plugin")
			},
			wantErr: true,
		},
		{
			name: "symlink to a file in a subdirectory",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Mkdir(filepath.Join(pluginDir, "sub"), 0o755))
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "sub", "real"), []byte("x"), 0o755))
				require.NoError(t, os.Symlink(filepath.Join(pluginDir, "sub", "real"), filepath.Join(pluginDir, "link")))
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
			name: "broken symlink",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.Symlink(filepath.Join(root, "missing"), filepath.Join(pluginDir, "link")))
				return "link"
			},
			wantErr: true,
		},
		{
			name: "parent reference",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(root, "outside"), []byte("x"), 0o755))
				return filepath.Join("..", "outside")
			},
			wantErr: true,
		},
		{
			name: "double dot inside a file name",
			setup: func(t *testing.T, root, pluginDir string) string {
				require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my..plugin"), []byte("x"), 0o755))
				return "my..plugin"
			},
			wantErr: true,
		},
		{
			name: "directory instead of a binary",
			setup: func(t *testing.T, root, pluginDir string) string {
				return "."
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

			registered := newDirectoryTestCatalog(t, pluginDir)
			setErr := registered.Set(t.Context(), pluginutil.SetPluginInput{
				Name:    "test-plugin",
				Type:    consts.PluginTypeDatabase,
				Command: command,
				Sha256:  []byte{'1'},
			})

			restored := newDirectoryTestCatalog(t, pluginDir)
			putRawCatalogEntry(t, restored, &pluginutil.PluginRunner{
				Name:    "test-plugin",
				Type:    consts.PluginTypeDatabase,
				Command: command,
				Sha256:  []byte{'1'},
			})
			got, getErr := restored.Get(t.Context(), "test-plugin", consts.PluginTypeDatabase, "")

			if tt.wantErr {
				require.Error(t, setErr, "registration must refuse %q", command)
				require.ErrorIs(t, getErr, ErrPluginOutsideDirectory, "lookup must refuse %q", command)
				require.Nil(t, got)
				return
			}
			require.NoError(t, setErr, "registration must accept %q", command)
			require.NoError(t, getErr, "lookup must accept %q", command)
			require.NotNil(t, got)
		})
	}
}

// TestPluginCatalog_Get_SymlinkedPluginDirectory verifies that a plugin
// directory that is itself reached through a symlink (common on macOS and in
// container images) does not cause legitimate plugin binaries to be rejected,
// by registration or by lookup, because the directory is resolved before
// comparing.
func TestPluginCatalog_Get_SymlinkedPluginDirectory(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	realDir := filepath.Join(root, "real-plugins")
	require.NoError(t, os.Mkdir(realDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(realDir, "my-plugin"), []byte("x"), 0o755))
	linkDir := filepath.Join(root, "link-plugins")
	require.NoError(t, os.Symlink(realDir, linkDir))

	registered := newDirectoryTestCatalog(t, linkDir)
	require.NoError(t, registered.Set(t.Context(), pluginutil.SetPluginInput{
		Name:    "test-plugin",
		Type:    consts.PluginTypeDatabase,
		Command: "my-plugin",
		Sha256:  []byte{'1'},
	}))

	restored := newDirectoryTestCatalog(t, linkDir)
	putRawCatalogEntry(t, restored, &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "my-plugin",
		Sha256:  []byte("sha"),
	})

	got, err := restored.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, filepath.Join(linkDir, "my-plugin"), got.Command)
}

// TestPluginCatalog_Get_ParentReferenceInName verifies that a stored entry
// whose name contains "..", which registration refuses, is refused by lookup
// too, even when its command is otherwise valid.
func TestPluginCatalog_Get_ParentReferenceInName(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	pluginDir := filepath.Join(root, "plugins")
	require.NoError(t, os.Mkdir(pluginDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "my-plugin"), []byte("x"), 0o755))

	catalog := newDirectoryTestCatalog(t, pluginDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "test..plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "my-plugin",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test..plugin", consts.PluginTypeSecrets, "")
	require.ErrorIs(t, err, ErrPluginOutsideDirectory)
	require.ErrorIs(t, err, consts.ErrPathContainsParentReferences)
	require.Nil(t, got)
}

// TestPluginCatalog_Get_BrokenSymlinkTargetCreatedLater verifies that a
// symlink in the plugin directory whose target does not exist yet is refused,
// and stays refused once a file appears at the target outside the directory.
// The binary is started by path some time after the lookup, and the operating
// system follows the symlink again at that point, so accepting a broken
// symlink would let whatever is later created at its target be executed.
func TestPluginCatalog_Get_BrokenSymlinkTargetCreatedLater(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	pluginDir := filepath.Join(root, "plugins")
	require.NoError(t, os.Mkdir(pluginDir, 0o755))
	target := filepath.Join(root, "later")
	require.NoError(t, os.Symlink(target, filepath.Join(pluginDir, "link")))

	catalog := newDirectoryTestCatalog(t, pluginDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "link",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.ErrorIs(t, err, ErrPluginOutsideDirectory, "broken symlink must be refused before its target exists")
	require.Nil(t, got)

	require.NoError(t, os.WriteFile(target, []byte("x"), 0o755))

	got, err = catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.ErrorIs(t, err, ErrPluginOutsideDirectory, "symlink must be refused once its target exists outside the plugin directory")
	require.Nil(t, got)
}

// TestPluginCatalog_Get_EncodedAndLookalikeCommands verifies that commands are
// never decoded. URL-encoded sequences and characters that only look like dots
// or slashes are literal file name characters: a command made of a single
// such name is used as a file directly in the plugin directory. Commands that
// contain a real ".." or "/" are judged on those characters alone, exactly as
// registration would judge them. A real file exists at root/outside/plugin,
// where a decoding interpretation of the commands would land.
func TestPluginCatalog_Get_EncodedAndLookalikeCommands(t *testing.T) {
	t.Parallel()

	literal := map[string]string{
		"URL encoded dots and slashes":        "%2E%2E%2Foutside%2Fplugin",
		"double URL encoded dots and slashes": "%252e%252e%252foutside%252fplugin",
		"fullwidth dots and solidus":          "\uff0e\uff0e\uff0foutside\uff0fplugin",
		"fullwidth reverse solidus":           "\uff0e\uff0e\uff3coutside\uff3cplugin",
		"two dot leader and division slash":   "\u2025\u2215outside\u2215plugin",
		"one dot leaders and fraction slash":  "\u2024\u2024\u2044outside\u2044plugin",
	}
	refused := map[string]string{
		"URL encoded dots with a real slash":   "%2e%2e/outside/plugin",
		"real dots with URL encoded slashes":   "..%2foutside%2fplugin",
		"backslash parent reference":           `..\outside\plugin`,
		"leading backslash parent references":  `\..\..\outside\plugin`,
		"backslash parent reference after dir": `sub\..\..\outside\plugin`,
		"parent reference with trailing space": `.. \outside\plugin`,
		"three dots":                           `...\outside\plugin`,
		"lookalike dots with a real slash":     "\u2025/outside/plugin",
	}

	run := func(t *testing.T, command string, wantErr bool) {
		root, err := filepath.EvalSymlinks(t.TempDir())
		require.NoError(t, err)
		pluginDir := filepath.Join(root, "plugins")
		require.NoError(t, os.Mkdir(pluginDir, 0o755))
		require.NoError(t, os.Mkdir(filepath.Join(root, "outside"), 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(root, "outside", "plugin"), []byte("x"), 0o755))

		catalog := newDirectoryTestCatalog(t, pluginDir)
		putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
			Name:    "test-plugin",
			Type:    consts.PluginTypeSecrets,
			Command: command,
			Sha256:  []byte("sha"),
		})

		got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
		if wantErr {
			require.ErrorIs(t, err, ErrPluginOutsideDirectory, "command %q", command)
			require.Nil(t, got)
			return
		}
		require.NoError(t, err, "command %q", command)
		require.NotNil(t, got)
		require.Equal(t, filepath.Join(pluginDir, command), got.Command, "command must be used literally")
	}

	for name, command := range literal {
		t.Run("literal/"+name, func(t *testing.T) {
			t.Parallel()
			run(t, command, false)
		})
	}
	for name, command := range refused {
		t.Run("refused/"+name, func(t *testing.T) {
			t.Parallel()
			run(t, command, true)
		})
	}
}

// TestPluginCatalog_Get_InvalidUTF8Command verifies that a command holding an
// overlong UTF-8 encoding of "../" (0xC0 0xAE is an invalid two-byte form of
// ".") is not decoded into a parent reference. The catalog stores entries as
// JSON, which replaces invalid UTF-8 with U+FFFD, and the command that comes
// back must still be a file directly in the plugin directory.
func TestPluginCatalog_Get_InvalidUTF8Command(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	pluginDir := filepath.Join(root, "plugins")
	require.NoError(t, os.Mkdir(pluginDir, 0o755))

	catalog := newDirectoryTestCatalog(t, pluginDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "\xc0\xae\xc0\xae\xc0\xafoutside",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, pluginDir, filepath.Dir(got.Command))
}

// TestPluginCatalog_Get_JSONEscapedTraversal verifies that the checks run on
// the command after the stored JSON has been decoded. JSON string escapes such
// as \u002e (".") and \/ ("/") are decoded by the JSON parser, so a stored
// entry can spell "../" without containing those bytes; the decoded value
// must still be refused.
func TestPluginCatalog_Get_JSONEscapedTraversal(t *testing.T) {
	t.Parallel()

	commands := map[string]string{
		"unicode escaped dots":                      `\u002e\u002e/outside`,
		"unicode escaped dots and escaped slash":    `\u002e\u002e\/outside`,
		"unicode escaped slash":                     `..\u002foutside`,
		"every character unicode escaped":           `\u002e\u002e\u002f\u006f\u0075\u0074\u0073\u0069\u0064\u0065`,
		"escaped slashes deeper than the directory": `..\/..\/..\/..\/..\/..\/..\/..\/bin\/sh`,
	}

	for name, rawCommand := range commands {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root, err := filepath.EvalSymlinks(t.TempDir())
			require.NoError(t, err)
			pluginDir := filepath.Join(root, "plugins")
			require.NoError(t, os.Mkdir(pluginDir, 0o755))
			require.NoError(t, os.WriteFile(filepath.Join(root, "outside"), []byte("x"), 0o755))

			catalog := newDirectoryTestCatalog(t, pluginDir)
			err = catalog.catalogView.Put(t.Context(), &logical.StorageEntry{
				Key:   "secret/test-plugin",
				Value: []byte(`{"name":"test-plugin","type":"secret","command":"` + rawCommand + `","sha256":"c2hh"}`),
			})
			require.NoError(t, err)

			got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
			require.ErrorIs(t, err, ErrPluginOutsideDirectory, "decoded command must be refused")
			require.Nil(t, got)
		})
	}
}

// TestPluginCatalog_Get_NULByteInCommand verifies that a command containing a
// NUL byte is refused. Such a path cannot be passed to the operating system,
// and it must not be mistaken for a missing binary and accepted.
func TestPluginCatalog_Get_NULByteInCommand(t *testing.T) {
	t.Parallel()

	root, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	pluginDir := filepath.Join(root, "plugins")
	require.NoError(t, os.Mkdir(pluginDir, 0o755))

	catalog := newDirectoryTestCatalog(t, pluginDir)
	putRawCatalogEntry(t, catalog, &pluginutil.PluginRunner{
		Name:    "test-plugin",
		Type:    consts.PluginTypeSecrets,
		Command: "my-plugin\x00",
		Sha256:  []byte("sha"),
	})

	got, err := catalog.Get(t.Context(), "test-plugin", consts.PluginTypeSecrets, "")
	require.Error(t, err)
	require.Nil(t, got)
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
	key := path.Join(entry.Type.String(), entry.Name)
	if entry.Version != "" {
		key = path.Join(key, entry.Version)
	}
	err = catalog.catalogView.Put(t.Context(), &logical.StorageEntry{
		Key:   key,
		Value: buf,
	})
	require.NoError(t, err)
}
