// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package version

import (
	"bytes"
	"fmt"
	"time"
)

type VersionInfo struct {
	// BinaryName is only for FullVersionNumber output, which defaults to "Vault"
	// when this is empty. We leave it out of JSON.
	BinaryName        string `json:"-"`
	Revision          string `json:"revision,omitempty"`
	Version           string `json:"version,omitempty"`
	VersionPrerelease string `json:"version_prerelease,omitempty"`
	VersionMetadata   string `json:"version_metadata,omitempty"`
	BuildDate         string `json:"build_date,omitempty"`
}

func GetVersion() *VersionInfo {
	ver := Version
	rel := VersionPrerelease
	md := VersionMetadata
	if GitDescribe != "" {
		ver = GitDescribe
	}

	return &VersionInfo{
		Revision:          GitCommit,
		Version:           ver,
		VersionPrerelease: rel,
		VersionMetadata:   md,
		BuildDate:         BuildDate,
	}
}

func GetVaultBuildDate() (time.Time, error) {
	buildDate, err := time.Parse(time.RFC3339, BuildDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse build date based on RFC3339: %w", err)
	}
	return buildDate, nil
}

func (c *VersionInfo) VersionNumber() string {
	if Version == "unknown" && VersionPrerelease == "unknown" {
		return "(version unknown)"
	}

	version := c.Version

	if c.VersionPrerelease != "" {
		version = fmt.Sprintf("%s-%s", version, c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, c.VersionMetadata)
	}

	return version
}

func (c *VersionInfo) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer
	binaryName := c.BinaryName
	if binaryName == "" {
		binaryName = "Vault"
	}

	if Version == "unknown" && VersionPrerelease == "unknown" {
		return fmt.Sprintf("%s (version unknown)", binaryName)
	}

	fmt.Fprintf(&versionString, "%s v%s", binaryName, c.Version)
	if c.VersionPrerelease != "" {
		fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		fmt.Fprintf(&versionString, "+%s", c.VersionMetadata)
	}

	if rev && c.Revision != "" {
		fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	if c.BuildDate != "" {
		fmt.Fprintf(&versionString, ", built %s", c.BuildDate)
	}

	return versionString.String()
}
