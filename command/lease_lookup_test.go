// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testLeaseLookupCommand(tb testing.TB) (*cli.MockUi, *LeaseLookupCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &LeaseLookupCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

// testLeaseLookupCommandMountAndLease mounts a leased secret backend and returns
// the leaseID of an item.
func testLeaseLookupCommandMountAndLease(tb testing.TB, client *api.Client) string {
	if err := client.Sys().Mount("testing", &api.MountInput{
		Type: "generic-leased",
	}); err != nil {
		tb.Fatal(err)
	}

	if _, err := client.Logical().Write("testing/foo", map[string]interface{}{
		"key":   "value",
		"lease": "5m",
	}); err != nil {
		tb.Fatal(err)
	}

	// Read the secret back to get the leaseID
	secret, err := client.Logical().Read("testing/foo")
	if err != nil {
		tb.Fatal(err)
	}
	if secret == nil || secret.LeaseID == "" {
		tb.Fatalf("missing secret or lease: %#v", secret)
	}

	return secret.LeaseID
}

// TestLeaseLookupCommand_Run tests basic lookup
func TestLeaseLookupCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		_ = testLeaseLookupCommandMountAndLease(t, client)

		ui, cmd := testLeaseLookupCommand(t)
		cmd.client = client

		code := cmd.Run(nil)
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		expectedMsg := "Missing ID!"
		if !strings.Contains(combined, expectedMsg) {
			t.Errorf("expected %q to contain %q", combined, expectedMsg)
		}
	})

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		leaseID := testLeaseLookupCommandMountAndLease(t, client)

		_, cmd := testLeaseLookupCommand(t)
		cmd.client = client

		code := cmd.Run([]string{leaseID})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testLeaseLookupCommand(t)
		assertNoTabs(t, cmd)
	})
}
