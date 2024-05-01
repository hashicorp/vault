// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package command_testonly

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/helper/timeutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func testOperatorUsageCommand(tb testing.TB) (*cli.MockUi, *command.OperatorUsageCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &command.OperatorUsageCommand{
		BaseCommand: &command.BaseCommand{
			UI: ui,
		},
	}
}

// TestOperatorUsageCommandRun writes mock activity log data and runs the
// operator usage command. The test verifies that the output contains the
// expected values per client type.
// This test cannot be run in parallel because it sets the VAULT_TOKEN env
// var
func TestOperatorUsageCommandRun(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	defer cluster.Cleanup()
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	client := cluster.Cores[0].Client
	_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{"enabled": "enable"})
	require.NoError(t, err)

	now := time.Now().UTC()

	_, err = clientcountutil.NewActivityLogData(client).
		NewPreviousMonthData(1).
		NewClientsSeen(6, clientcountutil.WithClientType("entity")).
		NewClientsSeen(4, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(2, clientcountutil.WithClientType("secret-sync")).
		NewClientsSeen(7, clientcountutil.WithClientType("pki-acme")).
		NewCurrentMonthData().
		NewClientsSeen(3, clientcountutil.WithClientType("entity")).
		NewClientsSeen(4, clientcountutil.WithClientType("non-entity-token")).
		NewClientsSeen(5, clientcountutil.WithClientType("secret-sync")).
		NewClientsSeen(8, clientcountutil.WithClientType("pki-acme")).
		Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES, generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES)
	require.NoError(t, err)

	ui, cmd := testOperatorUsageCommand(t)

	t.Setenv("VAULT_TOKEN", client.Token())
	start := timeutil.MonthsPreviousTo(1, now).Format(time.RFC3339)
	end := timeutil.EndOfMonth(now).UTC().Format(time.RFC3339)
	// Reset and check output
	code := cmd.Run([]string{
		"-address", client.Address(),
		"-tls-skip-verify",
		"-start-time", start,
		"-end-time", end,
	})
	require.Equal(t, 0, code, ui.ErrorWriter.String())
	output := ui.OutputWriter.String()
	outputLines := strings.Split(output, "\n")
	require.Equal(t, fmt.Sprintf("Period start: %s", start), outputLines[0])
	require.Equal(t, fmt.Sprintf("Period end: %s", end), outputLines[1])

	require.Contains(t, outputLines[3], "Secret sync")
	require.Contains(t, outputLines[3], "ACME clients")
	nsCounts := strings.Fields(outputLines[5])
	require.Equal(t, "[root]", nsCounts[0])
	require.Equal(t, "9", nsCounts[1])
	require.Equal(t, "8", nsCounts[2])
	require.Equal(t, "7", nsCounts[3])
	require.Equal(t, "15", nsCounts[4])
	require.Equal(t, "39", nsCounts[5])

	totalCounts := strings.Fields(outputLines[7])
	require.Equal(t, "Total", totalCounts[0])
	require.Equal(t, nsCounts[1:], totalCounts[1:])
}
