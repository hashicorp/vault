package command

import (
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"io"
	"strings"
	"testing"
)

func testKVMetadataPatchCommand(tb testing.TB) (*cli.MockUi, *KVMetadataPatchCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &KVMetadataPatchCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func kvMetadataPatchWithRetry(t *testing.T, client *api.Client, args []string, stdin *io.PipeReader) (int, string) {
	t.Helper()

	return retryKVCommand(t, func() (int, string) {
		ui, cmd := testKVMetadataPatchCommand(t)
		cmd.client = client

		if stdin != nil {
			cmd.testStdin = stdin
		}

		code := cmd.Run(args)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		return code, combined
	})
}

func kvMetadataPutWithRetry(t *testing.T, client *api.Client, args []string, stdin *io.PipeReader) (int, string) {
	t.Helper()

	return retryKVCommand(t, func() (int, string) {
		ui, cmd := testKVMetadataPutCommand(t)
		cmd.client = client

		if stdin != nil {
			cmd.testStdin = stdin
		}

		code := cmd.Run(args)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		return code, combined
	})
}

func TestKvMetadataPatchEmptyArgs (t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount error: %#v", err)
	}

	args := make([]string, 0)
	code, combined := kvMetadataPatchWithRetry(t, client, args, nil)

	expectedCode := 1
	expectedOutput := "Not enough arguments"

	if code != expectedCode {
		t.Fatalf("expected code to be %d but was %d for patch cmd with args %#v", expectedCode, code, args)
	}

	if !strings.Contains(combined, expectedOutput) {
		t.Fatalf("expected output to be %q but was %q for patch cmd with args %#v", expectedOutput, combined, args)
	}
}