// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package version

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func replaceVersion(v, vp, vm string) func() {
	origV := Version
	origVP := VersionPrerelease
	origVM := VersionMetadata

	Version = v
	VersionPrerelease = vp
	VersionMetadata = vm

	return func() {
		Version = origV
		VersionPrerelease = origVP
		VersionMetadata = origVM
	}
}

// TestGetVersion verifies that each call returns the build version without
// sharing display overrides or adding them to JSON.
func TestGetVersion(t *testing.T) {
	// This test cannot be parallelized because it messes with some global
	// variables that determine the version information.
	restoreVersionFunc := replaceVersion("1.2.3", "", "")
	defer restoreVersionFunc()

	// Test the general case
	vi := GetVersion()
	assert.Equal(t, "1.2.3", vi.Version)
	assert.Equal(t, "", vi.VersionPrerelease)
	assert.Equal(t, "", vi.VersionMetadata)
	assert.Equal(t, "", vi.Revision)
	assert.Equal(t, "", vi.BuildDate)

	originalJSON, err := json.Marshal(vi)
	require.NoError(t, err, "build information must serialize before a display override")
	vi.BinaryName = "Vault Client"
	fresh := GetVersion()
	require.Empty(t, fresh.BinaryName, "a display override must not leak to later version instances")
	require.Equal(t, "Vault v1.2.3", fresh.FullVersionNumber(false), "new instances must retain the default product name")
	overriddenJSON, err := json.Marshal(vi)
	require.NoError(t, err, "build information must serialize after a display override")
	require.Equal(t, originalJSON, overriddenJSON, "the display name must not change serialized version information")

	// Test the git describe case
	origGitDescribe := GitDescribe
	GitDescribe = "git-describe"
	vi = GetVersion()
	assert.Equal(t, "git-describe", vi.Version)

	GitDescribe = origGitDescribe
}

func TestVersionNumber(t *testing.T) {
	// This test cannot be parallelized because it messes with some global
	// variables that determine the version information.
	restoreVersionFunc := replaceVersion("unknown", "unknown", "")
	defer restoreVersionFunc()

	// Test the unknown version case
	vi := GetVersion()
	assert.Equal(t, "(version unknown)", vi.VersionNumber())

	replaceVersion("1.2.3", "", "")

	// Test the pre-release case
	vi = GetVersion()
	vi.VersionPrerelease = "rc1"
	assert.Equal(t, "1.2.3-rc1", vi.VersionNumber())

	// Test the pre-release and metadata version case
	vi.VersionMetadata = "ent"
	assert.Equal(t, "1.2.3-rc1+ent", vi.VersionNumber())

	// Test the metadata only version case
	vi.VersionPrerelease = ""
	assert.Equal(t, "1.2.3+ent", vi.VersionNumber())
}

// TestFullVersionNumber verifies that display names don't change version numbers
// or build details and still identify the binary when the version is unknown.
func TestFullVersionNumber(t *testing.T) {
	// Keep these serial because the unknown-version cases change package globals.
	for _, tc := range []struct {
		name         string
		info         VersionInfo
		unknownBuild bool
		revision     bool
		wantFull     string
		wantNumber   string
	}{
		{
			name:     "zero_value",
			wantFull: "Vault v",
		},
		{
			name:       "unset_name",
			info:       VersionInfo{Version: "1.2.3"},
			wantFull:   "Vault v1.2.3",
			wantNumber: "1.2.3",
		},
		{
			name:       "explicit_default",
			info:       VersionInfo{BinaryName: "Vault", Version: "1.2.3"},
			wantFull:   "Vault v1.2.3",
			wantNumber: "1.2.3",
		},
		{
			name:       "prerelease",
			info:       VersionInfo{Version: "1.2.3", VersionPrerelease: "rc1"},
			wantFull:   "Vault v1.2.3-rc1",
			wantNumber: "1.2.3-rc1",
		},
		{
			name:       "enterprise",
			info:       VersionInfo{Version: "1.2.3", VersionMetadata: "ent"},
			wantFull:   "Vault v1.2.3+ent",
			wantNumber: "1.2.3+ent",
		},
		{
			name:       "revision",
			info:       VersionInfo{Version: "1.2.3", Revision: "ab1234f"},
			revision:   true,
			wantFull:   "Vault v1.2.3 (ab1234f)",
			wantNumber: "1.2.3",
		},
		{
			name:       "date_without_revision",
			info:       VersionInfo{Version: "1.2.3", Revision: "ab1234f", BuildDate: "2023-10-20"},
			wantFull:   "Vault v1.2.3, built 2023-10-20",
			wantNumber: "1.2.3",
		},
		{
			name: "client_all_fields",
			info: VersionInfo{
				BinaryName:        "Vault Client",
				Version:           "1.2.3",
				VersionPrerelease: "rc1",
				VersionMetadata:   "ent.hsm.Vault.custom",
				Revision:          "build-Vault-123",
				BuildDate:         "2026-10-27T12:34:56Z",
			},
			revision:   true,
			wantFull:   "Vault Client v1.2.3-rc1+ent.hsm.Vault.custom (build-Vault-123), built 2026-10-27T12:34:56Z",
			wantNumber: "1.2.3-rc1+ent.hsm.Vault.custom",
		},
		{
			name:       "custom_name_git_describe",
			info:       VersionInfo{BinaryName: "Custom Vault", Version: "v1.2.3-rc1-7-gabc123"},
			wantFull:   "Custom Vault vv1.2.3-rc1-7-gabc123",
			wantNumber: "v1.2.3-rc1-7-gabc123",
		},
		{
			name:         "unknown_default",
			info:         VersionInfo{Version: "ignored", VersionMetadata: "ent", Revision: "ignored", BuildDate: "ignored"},
			unknownBuild: true,
			revision:     true,
			wantFull:     "Vault (version unknown)",
			wantNumber:   "(version unknown)",
		},
		{
			name:         "unknown_client",
			info:         VersionInfo{BinaryName: "Vault Client", Version: "ignored", VersionMetadata: "ent", Revision: "ignored", BuildDate: "ignored"},
			unknownBuild: true,
			revision:     true,
			wantFull:     "Vault Client (version unknown)",
			wantNumber:   "(version unknown)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buildVersion, prerelease := "1.2.3", ""
			if tc.unknownBuild {
				buildVersion, prerelease = "unknown", "unknown"
			}
			t.Cleanup(replaceVersion(buildVersion, prerelease, ""))
			info := tc.info
			require.Equal(t, tc.wantFull, info.FullVersionNumber(tc.revision), "display formatting must preserve the product name and build fields")
			require.Equal(t, tc.wantNumber, info.VersionNumber(), "numeric version output must not include the display name")
			require.Equal(t, tc.info, info, "formatting must not mutate version information")
		})
	}
}
