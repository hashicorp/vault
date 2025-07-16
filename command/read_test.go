// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func testReadCommand(tb testing.TB) (*cli.MockUi, *ReadCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &ReadCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestReadCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"proper_args",
			[]string{"foo", "bar=baz"},
			"No value found at foo\n",
			2,
		},
		{
			"not_found",
			[]string{"nope/not/once/never"},
			"",
			2,
		},
		{
			"default",
			[]string{"secret/read/foo"},
			"foo",
			0,
		},
		{
			"field",
			[]string{
				"-field", "foo",
				"secret/read/foo",
			},
			"bar",
			0,
		},
		{
			"field_not_found",
			[]string{
				"-field", "not-a-real-field",
				"secret/read/foo",
			},
			"not present in secret",
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

				if _, err := client.Logical().Write("secret/read/foo", map[string]interface{}{
					"foo": "bar",
				}); err != nil {
					t.Fatal(err)
				}

				ui, cmd := testReadCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("expected %d to be %d", code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("%s: expected %q to contain %q", tc.name, combined, tc.out)
				}
			})
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testReadCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/foo",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error reading secret/foo: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_data_object_from_api_response", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testReadCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"sys/health",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		expected := []string{
			"cluster_id", "cluster_name", "initialized", "performance_standby", "replication_dr_mode", "replication_performance_mode", "sealed",
			"server_time_utc", "standby", "version",
		}
		for _, expectedField := range expected {
			if !strings.Contains(combined, expectedField) {
				t.Errorf("expected %q to contain %q", combined, expected)
			}
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testReadCommand(t)
		assertNoTabs(t, cmd)
	})
}

// TestRead_Snapshot tests that the read_snapshot_id query parameter is added
// to the request when the -snapshot-id flag is used.
func TestRead_Snapshot(t *testing.T) {
	t.Parallel()
	mockVaultServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		snapID := r.URL.Query().Get("read_snapshot_id")
		if snapID != "abcd" {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(`{"secret":{"data":{"foo":"bar"}}}`))
	}))
	defer mockVaultServer.Close()

	cfg := api.DefaultConfig()
	cfg.Address = mockVaultServer.URL
	client, err := api.NewClient(cfg)
	require.NoError(t, err)

	ui, cmd := testReadCommand(t)
	cmd.client = client

	// a read command with a snapshot id shouldn't error
	code := cmd.Run([]string{
		"-snapshot-id", "abcd", "path/to/item",
	})
	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
	require.Equal(t, 0, code, combined)

	// check that the raw flag also works with a snapshot id
	code = cmd.Run([]string{
		"-format", "raw", "-snapshot-id", "abcd", "path/to/item",
	})
	combined = ui.OutputWriter.String() + ui.ErrorWriter.String()
	require.Equal(t, 0, code, combined)
}
