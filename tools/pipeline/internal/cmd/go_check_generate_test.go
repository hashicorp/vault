// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/golang"
	"github.com/stretchr/testify/require"
)

// TestGoCheckGenerateCmd_Flags verifies that the generate check binds its
// flags and rejects the tools build tag before running anything.
func TestGoCheckGenerateCmd_Flags(t *testing.T) {
	goCheckGenerateReq = &golang.CheckGenerateReq{}

	cmd := newGoCheckGenerateCmd()
	suppressOutput(t, cmd)
	cmd.SetArgs([]string{"--tags", "ent,tools", "--module", "sdk"})
	require.ErrorContains(t, cmd.Execute(), "tools build tag")
	require.Equal(t, &golang.CheckGenerateReq{Tags: "ent,tools", Modules: []string{"sdk"}}, goCheckGenerateReq)
}

// TestGoCheckGenerateRerun verifies that the suggested fix reruns the check
// with every flag that selects what's generated, quoted for a shell.
func TestGoCheckGenerateRerun(t *testing.T) {
	t.Parallel()

	require.Equal(t, "pipeline go check generate", goCheckGenerateRerun(&golang.CheckGenerateReq{}))
	require.Equal(t,
		"pipeline go check generate --go-work '/my repo/go.work' --tags ent,enterprise --module sdk --module 'it'\\''s'",
		goCheckGenerateRerun(&golang.CheckGenerateReq{
			GoWork:  "/my repo/go.work",
			Tags:    "ent enterprise",
			Modules: []string{"sdk", "it's"},
		}),
	)
}
