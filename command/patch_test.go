// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"io"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testPatchCommand(tb testing.TB) (*cli.MockUi, *PatchCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PatchCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPatchCommand_Run(t *testing.T) {
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
			"empty_kvs",
			[]string{"secret/write/foo"},
			"Must supply data or use -force",
			1,
		},
		{
			"force_kvs",
			[]string{"-force", "pki/roles/example"},
			"allow_localhost",
			0,
		},
		{
			"force_f_kvs",
			[]string{"-f", "pki/roles/example"},
			"allow_localhost",
			0,
		},
		{
			"kvs_no_value",
			[]string{"pki/roles/example", "foo"},
			"Failed to parse K=V data",
			1,
		},
		{
			"single_value",
			[]string{"pki/roles/example", "allow_localhost=true"},
			"allow_localhost",
			0,
		},
		{
			"multi_value",
			[]string{"pki/roles/example", "allow_localhost=true", "allowed_domains=true"},
			"allow_localhost",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().Mount("pki", &api.MountInput{
				Type: "pki",
			}); err != nil {
				t.Fatalf("pki mount error: %#v", err)
			}

			if _, err := client.Logical().Write("pki/roles/example", nil); err != nil {
				t.Fatalf("failed to prime role: %v", err)
			}

			if _, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
				"key_type":    "ec",
				"common_name": "Root X1",
			}); err != nil {
				t.Fatalf("failed to prime CA: %v", err)
			}

			ui, cmd := testPatchCommand(t)
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

	t.Run("stdin_full", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().Mount("pki", &api.MountInput{
			Type: "pki",
		}); err != nil {
			t.Fatalf("pki mount error: %#v", err)
		}

		if _, err := client.Logical().Write("pki/roles/example", nil); err != nil {
			t.Fatalf("failed to prime role: %v", err)
		}

		if _, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
			"key_type":    "ec",
			"common_name": "Root X1",
		}); err != nil {
			t.Fatalf("failed to prime CA: %v", err)
		}

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(`{"allow_localhost":"false","allow_wildcard_certificates":"false"}`))
			stdinW.Close()
		}()

		ui, cmd := testPatchCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"pki/roles/example", "-",
		})
		if code != 0 {
			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			t.Fatalf("expected retcode=%d to be 0\nOutput:\n%v", code, combined)
		}

		secret, err := client.Logical().Read("pki/roles/example")
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil || secret.Data == nil {
			t.Fatal("expected secret to have data")
		}
		if exp, act := false, secret.Data["allow_localhost"].(bool); exp != act {
			t.Errorf("expected allowed_localhost=%v to be %v", act, exp)
		}
		if exp, act := false, secret.Data["allow_wildcard_certificates"].(bool); exp != act {
			t.Errorf("expected allow_wildcard_certificates=%v to be %v", act, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testPatchCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"foo/bar", "a=b",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error writing data to foo/bar: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testPatchCommand(t)
		assertNoTabs(t, cmd)
	})
}
