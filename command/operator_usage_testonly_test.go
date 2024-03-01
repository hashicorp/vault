// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package command

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/helper/timeutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil"
	"github.com/hashicorp/vault/sdk/helper/clientcountutil/generation"
	"github.com/stretchr/testify/require"
)

func testOperatorUsageCommand(tb testing.TB) (*cli.MockUi, *OperatorUsageCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorUsageCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorUsageCommandRun(t *testing.T) {
	t.Parallel()

	t.Run("client types", func(t *testing.T) {
		t.Parallel()

		client, _, closer := testVaultServerUnseal(t)
		defer closer()

		_, err := client.Logical().Write("sys/internal/counters/config", map[string]interface{}{"enabled": "enable"})
		require.NoError(t, err)

		now := time.Now().UTC()

		_, err = clientcountutil.NewActivityLogData(client).
			NewPreviousMonthData(1).
			NewClientsSeen(6, clientcountutil.WithClientType("entity")).
			NewClientsSeen(4, clientcountutil.WithClientType("non-entity-token")).
			NewClientsSeen(2, clientcountutil.WithClientType("secret-sync")).
			NewCurrentMonthData().
			NewClientsSeen(3, clientcountutil.WithClientType("entity")).
			NewClientsSeen(4, clientcountutil.WithClientType("non-entity-token")).
			NewClientsSeen(5, clientcountutil.WithClientType("secret-sync")).
			Write(context.Background(), generation.WriteOptions_WRITE_ENTITIES, generation.WriteOptions_WRITE_PRECOMPUTED_QUERIES)
		require.NoError(t, err)

		ui, cmd := testOperatorUsageCommand(t)
		cmd.client = client

		start := timeutil.MonthsPreviousTo(1, now).Format(time.RFC3339)
		end := timeutil.EndOfMonth(now).UTC().Format(time.RFC3339)
		// Reset and check output
		code := cmd.Run([]string{
			"-start-time", start,
			"-end-time", end,
		})
		require.Equal(t, 0, code)
		output := ui.OutputWriter.String()
		outputLines := strings.Split(output, "\n")
		require.Equal(t, fmt.Sprintf("Period start: %s", start), outputLines[0])
		require.Equal(t, fmt.Sprintf("Period end: %s", end), outputLines[1])

		require.Contains(t, outputLines[3], "Secret sync")
		nsCounts := strings.Fields(outputLines[5])
		require.Equal(t, "[root]", nsCounts[0])
		require.Equal(t, "9", nsCounts[1])
		require.Equal(t, "8", nsCounts[2])
		require.Equal(t, "7", nsCounts[3])
		require.Equal(t, "24", nsCounts[4])

		totalCounts := strings.Fields(outputLines[7])
		require.Equal(t, "Total", totalCounts[0])
		require.Equal(t, nsCounts[1:], totalCounts[1:])
	})
}
