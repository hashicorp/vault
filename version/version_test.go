// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func replaceVersion(v, vp string) func() {
	origV := Version
	origVP := VersionPrerelease

	Version = v
	VersionPrerelease = vp

	return func() {
		Version = origV
		VersionPrerelease = origVP
	}
}

func TestGetVersion(t *testing.T) {
	// This test cannot be parallelized because it messes with some global
	// variables that determine the version information.
	restoreVersionFunc := replaceVersion("1.2.3", "")
	defer restoreVersionFunc()

	// Test the general case
	vi := GetVersion()
	assert.Equal(t, "1.2.3", vi.Version)
	assert.Equal(t, "", vi.VersionPrerelease)
	assert.Equal(t, "", vi.VersionMetadata)
	assert.Equal(t, "", vi.Revision)
	assert.Equal(t, "", vi.BuildDate)

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
	restoreVersionFunc := replaceVersion("unknown", "unknown")
	defer restoreVersionFunc()

	// Test the unknown version case
	vi := GetVersion()
	assert.Equal(t, "(version unknown)", vi.VersionNumber())

	replaceVersion("1.2.3", "")

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

func TestFullVersionNumber(t *testing.T) {
	// This test cannot be parallelized because it messes with some global
	// variables that determine the version information.
	restoreVersionFunc := replaceVersion("unknown", "unknown")
	defer restoreVersionFunc()

	// Test the unknown version case
	vi := GetVersion()
	assert.Equal(t, "Vault (version unknown)", vi.FullVersionNumber(false))

	// Test the no pre-release, metadata, revision, build date case
	replaceVersion("1.2.3", "")
	vi = GetVersion()
	assert.Equal(t, "Vault v1.2.3", vi.FullVersionNumber(false))

	// Test the pre-release case
	vi.VersionPrerelease = "rc1"
	assert.Equal(t, "Vault v1.2.3-rc1", vi.FullVersionNumber(false))

	// Test the metadata case
	vi.VersionPrerelease = ""
	vi.VersionMetadata = "ent"
	assert.Equal(t, "Vault v1.2.3+ent", vi.FullVersionNumber(false))

	// Test the revision case
	vi.VersionMetadata = ""
	vi.Revision = "ab1234f"
	assert.Equal(t, "Vault v1.2.3 (ab1234f)", vi.FullVersionNumber(true))

	// Test the build date case
	vi.BuildDate = "2023-10-20"
	assert.Equal(t, "Vault v1.2.3, built 2023-10-20", vi.FullVersionNumber(false))

	// Test the case where all of the things are set
	vi.VersionPrerelease = "rc1"
	vi.VersionMetadata = "ent"
	assert.Equal(t, "Vault v1.2.3-rc1+ent (ab1234f), built 2023-10-20", vi.FullVersionNumber(true))
}
