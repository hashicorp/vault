// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/stretchr/testify/require"
)

func testStatusCommand(tb testing.TB) (*cli.MockUi, *StatusCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &StatusCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

// TestStatusCommand_RaftCluster creates a raft cluster and verifies that a
// follower has "Removed From Cluster" returned as false in the status command.
// The test then removes that follower, and checks that "Removed From Cluster"
// is now true
func TestStatusCommand_RaftCluster(t *testing.T) {
	t.Parallel()
	cluster := testVaultRaftCluster(t)
	defer cluster.Cleanup()

	toRemove := cluster.Cores[1]
	expectRemovedFromCluster := func(expectCode int, removed bool) {
		ui, cmd := testStatusCommand(t)
		cmd.client = toRemove.Client
		code := cmd.Run(nil)
		require.Equal(t, expectCode, code)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		require.Regexp(t, fmt.Sprintf(".*Removed From Cluster\\s+%t.*", removed), combined)
	}

	expectRemovedFromCluster(0, false)

	_, err := cluster.Cores[0].Client.Logical().Write("sys/storage/raft/remove-peer",
		map[string]interface{}{
			"server_id": toRemove.NodeID,
		})
	require.NoError(t, err)
	testhelpers.RetryUntil(t, 10*time.Second, func() error {
		if !toRemove.Sealed() {
			return errors.New("core not sealed")
		}
		return nil
	})
	expectRemovedFromCluster(2, true)
}

func TestStatusCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		args   []string
		sealed bool
		out    string
		code   int
	}{
		{
			"unsealed",
			nil,
			false,
			"Sealed          false",
			0,
		},
		{
			"sealed",
			nil,
			true,
			"Sealed             true",
			2,
		},
		{
			"args",
			[]string{"foo"},
			false,
			"Too many arguments",
			1,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()

				if tc.sealed {
					if err := client.Sys().Seal(); err != nil {
						t.Fatal(err)
					}
				}

				ui, cmd := testStatusCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			})
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testStatusCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error checking seal status: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testStatusCommand(t)
		assertNoTabs(t, cmd)
	})
}
