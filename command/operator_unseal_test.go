package command

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func testOperatorUnsealCommand(tb testing.TB) (*cli.MockUi, *OperatorUnsealCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorUnsealCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorUnsealCommand_Run(t *testing.T) {
	t.Parallel()

	t.Run("error_non_terminal", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorUnsealCommand(t)
		cmd.client = client
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

		ui, cmd := testOperatorUnsealCommand(t)
		cmd.client = client
		cmd.testOutput = ioutil.Discard

		// Reset and check output
		code := cmd.Run([]string{
			"-reset",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}
		expected := "0/3"
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

		for _, key := range keys {
			ui, cmd := testOperatorUnsealCommand(t)
			cmd.client = client
			cmd.testOutput = ioutil.Discard

			// Reset and check output
			code := cmd.Run([]string{
				key,
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
			}
		}

		status, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if status.Sealed {
			t.Error("expected unsealed")
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testOperatorUnsealCommand(t)
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

		_, cmd := testOperatorUnsealCommand(t)
		assertNoTabs(t, cmd)
	})
}

func TestOperatorUnsealCommand_Format(t *testing.T) {
	defer func() {
		os.Setenv(EnvVaultCLINoColor, "")
	}()

	client, keys, closer := testVaultServerUnseal(t)
	defer closer()

	// Seal so we can unseal
	if err := client.Sys().Seal(); err != nil {
		t.Fatal(err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	args, format, _ := setupEnv([]string{"operator", "unseal", "-format", "json"})
	if format != "json" {
		t.Fatalf("expected %q, got %q", "json", format)
	}

	// Unseal with one key
	code := RunCustom(append(args, []string{
		keys[0],
	}...), runOpts)
	if exp := 0; code != exp {
		t.Errorf("expected %d to be %d: %s", code, exp, stderr.String())
	}

	if !json.Valid(stdout.Bytes()) {
		t.Error("expected output to be valid JSON")
	}
}
