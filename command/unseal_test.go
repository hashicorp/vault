package command

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testUnsealCommand(tb testing.TB) (*cli.MockUi, *UnsealCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &UnsealCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestUnsealCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("error_non_terminal", func(t *testing.T) {
		t.Parallel()

		ui, cmd := testUnsealCommand(t)
		cmd.testOutput = ioutil.Discard

		code := cmd.Run(nil)
		if exp := 1; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "is not a terminal"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("reset", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Seal so we can unseal
		if err := client.Sys().Seal(); err != nil {
			t.Fatal(err)
		}

		// Enter an unseal key
		if _, err := client.Sys().Unseal(keys[0]); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testUnsealCommand(t)
		cmd.client = client
		cmd.testOutput = ioutil.Discard

		// Reset and check output
		code := cmd.Run([]string{
			"-reset",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}
		expected := "Unseal Progress: 0"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("full", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Seal so we can unseal
		if err := client.Sys().Seal(); err != nil {
			t.Fatal(err)
		}

		for i, key := range keys {
			ui, cmd := testUnsealCommand(t)
			cmd.client = client
			cmd.testOutput = ioutil.Discard

			// Reset and check output
			code := cmd.Run([]string{
				key,
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d", code, exp)
			}
			expected := fmt.Sprintf("Unseal Progress: %d", (i+1)%3) // 1, 2, 0
			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, expected) {
				t.Errorf("expected %q to contain %q", combined, expected)
			}
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testUnsealCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"abcd",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error unsealing: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testUnsealCommand(t)
		assertNoTabs(t, cmd)
	})
}
