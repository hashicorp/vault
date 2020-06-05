package command

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testKVMetadataPutCommand(tb testing.TB) (*cli.MockUi, *KVMetadataPutCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &KVMetadataPutCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestKvMetadataPutCommandDeleteVersionAfter(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatal(err)
	}

	ui, cmd := testKVMetadataPutCommand(t)
	cmd.client = client

	// Set a limit of 1s first.
	code := cmd.Run([]string{"-delete-version-after=1s", "kv/secret/my-secret"})
	if code != 0 {
		t.Errorf("expected %d but received %d", 0, code)
	}

	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
	if !strings.Contains(combined, "Success! Data written to: kv/metadata/secret/my-secret\n") {
		t.Errorf("expected %q but received %q", "Success! Data written to: kv/metadata/secret/my-secret\n", combined)
	}

	secret, err := client.Logical().Read("kv/metadata/secret/my-secret")
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["delete_version_after"] != "1s" {
		t.Fatalf("expected 1s but received %q", secret.Data["delete_version_after"])
	}

	// Now verify that we can return it to 0s.
	ui, cmd = testKVMetadataPutCommand(t)
	cmd.client = client

	// Set a limit of 1s first.
	code = cmd.Run([]string{"-delete-version-after=0", "kv/secret/my-secret"})
	if code != 0 {
		t.Errorf("expected %d but received %d", 0, code)
	}

	combined = ui.OutputWriter.String() + ui.ErrorWriter.String()
	if !strings.Contains(combined, "Success! Data written to: kv/metadata/secret/my-secret\n") {
		t.Errorf("expected %q but received %q", "Success! Data written to: kv/metadata/secret/my-secret\n", combined)
	}

	secret, err = client.Logical().Read("kv/metadata/secret/my-secret")
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["delete_version_after"] != "0s" {
		t.Fatalf("expected 0s but received %q", secret.Data["delete_version_after"])
	}
}
