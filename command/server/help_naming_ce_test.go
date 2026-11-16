// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestServerNaming_EventsNotFoundHint checks that when the server has no events
// endpoint, `events subscribe` tells users to check sys/experiments with the
// binary they ran. CE servers don't serve event subscriptions, so the endpoint
// is missing on a CE test cluster.
func TestServerNaming_EventsNotFoundHint(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	ui := cli.NewMockUi()
	cmd, err := initCommands(ui, ui, &base.RunOptions{
		BinaryName: "renamed-vault", Client: cluster.Cores[0].Client, HCPTokenHelper: noHCPTokenHelper{},
	})["events subscribe"]()
	require.NoError(t, err, "command factory")
	require.Equal(t, 1, cmd.Run([]string{"kv-v2/data-write"}), "missing events endpoint")
	output := strings.Join(strings.Fields(ui.ErrorWriter.String()), " ")
	require.Contains(t, output, "events endpoint not found", "missing events endpoint")
	require.Contains(t, output, "check `renamed-vault read sys/experiments`", "hint uses the binary name")
}
