// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
)

func testAuditEnableCommand(tb testing.TB) (*cli.MockUi, *AuditEnableCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &AuditEnableCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestAuditEnableCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"empty",
			nil,
			"Error enabling audit device: audit type missing. Valid types include 'file', 'socket' and 'syslog'.",
			1,
		},
		{
			"not_a_valid_type",
			[]string{"nope_definitely_not_a_valid_type_like_ever"},
			"",
			2,
		},
		{
			"enable",
			[]string{"file", "file_path=discard"},
			"Success! Enabled the file audit device at: file/",
			0,
		},
		{
			"enable_path",
			[]string{
				"-path", "audit_path",
				"file",
				"file_path=discard",
			},
			"Success! Enabled the file audit device at: audit_path/",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			ui, cmd := testAuditEnableCommand(t)
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

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testAuditEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-path", "audit_enable_integration/",
			"-description", "The best kind of test",
			"file",
			"file_path=discard",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Enabled the file audit device at: audit_enable_integration/"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		audits, err := client.Sys().ListAudit()
		if err != nil {
			t.Fatal(err)
		}

		auditInfo, ok := audits["audit_enable_integration/"]
		if !ok {
			t.Fatalf("expected audit to exist")
		}
		if exp := "file"; auditInfo.Type != exp {
			t.Errorf("expected %q to be %q", auditInfo.Type, exp)
		}
		if exp := "The best kind of test"; auditInfo.Description != exp {
			t.Errorf("expected %q to be %q", auditInfo.Description, exp)
		}

		filePath, ok := auditInfo.Options["file_path"]
		if !ok || filePath != "discard" {
			t.Errorf("missing some options: %#v", auditInfo)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testAuditEnableCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"pki",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error enabling audit device: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testAuditEnableCommand(t)
		assertNoTabs(t, cmd)
	})

	t.Run("mount_all", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerAllBackends(t)
		defer closer()

		for _, name := range []string{"file", "socket", "syslog"} {
			ui, cmd := testAuditEnableCommand(t)
			cmd.client = client

			args := []string{name}
			switch name {
			case "file":
				args = append(args, "file_path=discard")
			case "socket":
				args = append(args, "address=127.0.0.1:8888", "skip_test=true")
			case "syslog":
				if _, exists := os.LookupEnv("WSLENV"); exists {
					t.Log("skipping syslog test on WSL")
					continue
				}
			}
			code := cmd.Run(args)
			if exp := 0; code != exp {
				t.Errorf("type %s, expected %d to be %d - %s", name, code, exp, ui.OutputWriter.String()+ui.ErrorWriter.String())
			}
		}
	})
}
