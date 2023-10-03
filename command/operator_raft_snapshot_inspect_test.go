// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/mitchellh/cli"
)

func testOperatorRaftSnapshotInspectCommand(tb testing.TB) (*cli.MockUi, *OperatorRaftSnapshotInspectCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorRaftSnapshotInspectCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func createSnapshot(tb testing.TB) (*os.File, func(), error) {
	// Create new raft backend
	r, raftDir := raft.GetRaft(tb, true, false)
	defer os.RemoveAll(raftDir)

	// Write some data
	for i := 0; i < 100; i++ {
		err := r.Put(context.Background(), &physical.Entry{
			Key:   fmt.Sprintf("key-%d", i),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		})
		if err != nil {
			return nil, nil, fmt.Errorf("Error adding data to snapshot %s", err)
		}
	}

	// Create temporary file to save snapshot to
	snap, err := os.CreateTemp("", "temp_snapshot.snap")
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating temporary file %s", err)
	}

	cleanup := func() {
		err := os.RemoveAll(snap.Name())
		if err != nil {
			tb.Errorf("Error deleting temporary snapshot %s", err)
		}
	}

	// Save snapshot
	err = r.Snapshot(snap, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("Error saving raft snapshot %s", err)
	}

	return snap, cleanup, nil
}

func TestOperatorRaftSnapshotInspectCommand_Run(t *testing.T) {
	t.Parallel()

	file1, cleanup1, err := createSnapshot(t)
	if err != nil {
		t.Fatalf("Error creating snapshot %s", err)
	}

	file2, cleanup2, err := createSnapshot(t)
	if err != nil {
		t.Fatalf("Error creating snapshot %s", err)
	}

	cases := []struct {
		name    string
		args    []string
		out     string
		code    int
		cleanup func()
	}{
		{
			"too_many_args",
			[]string{"test.snap", "test"},
			"Too many arguments",
			1,
			nil,
		},
		{
			"default",
			[]string{file1.Name()},
			"ID           bolt-snapshot",
			0,
			cleanup1,
		},
		{
			"all_flags",
			[]string{"-kvdetails", "-kvdepth", "10", "-kvfilter", "key", file2.Name()},
			"Key Name      Count",
			0,
			cleanup2,
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

				ui, cmd := testOperatorRaftSnapshotInspectCommand(t)

				cmd.client = client
				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}

				if tc.cleanup != nil {
					tc.cleanup()
				}
			})
		}
	})
}
