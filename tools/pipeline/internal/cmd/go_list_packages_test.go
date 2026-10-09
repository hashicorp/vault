// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// TestGoListPackagesCmd_Flags verifies that --module is repeatable without
// comma splitting and that --tags and --go-work reach the request.
func TestGoListPackagesCmd_Flags(t *testing.T) {
	goListPackagesReq = &golang.ListPackagesReq{}

	cmd := newGoListPackagesCmd()
	suppressOutput(t, cmd)
	require.NoError(t, cmd.ParseFlags([]string{
		"--go-work", "/repo/go.work",
		"--tags", ",deadlock",
		"--module", ".",
		"--module", "sdk,api",
	}))

	require.Equal(t, &golang.ListPackagesReq{
		GoWork:  "/repo/go.work",
		Tags:    ",deadlock",
		Modules: []string{".", "sdk,api"},
	}, goListPackagesReq)
}
