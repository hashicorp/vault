// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package version

import (
	_ "embed"
	"strings"
)

var (
	// The git commit that was compiled. This will be filled in by the compiler.
	GitCommit   string
	GitDescribe string

	// The compilation date. This will be filled in by the compiler.
	BuildDate string

	// Whether cgo is enabled or not; set at build time
	CgoEnabled bool

	// Version and VersionPrerelease info are now being embedded directly from the VERSION file.
	// VersionMetadata is being passed in via ldflags in CI, otherwise the default set here is used.
	//go:embed VERSION
	fullVersion                   string
	Version, VersionPrerelease, _ = strings.Cut(strings.TrimSpace(fullVersion), "-")
	VersionMetadata               = ""
)
