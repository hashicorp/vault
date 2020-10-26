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

	basePath := t.Name() + "/"
	if err := client.Sys().Mount(basePath, &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatal(err)
	}

	ui, cmd := testKVMetadataPutCommand(t)
	cmd.client = client

	// Set a limit of 1s first.
	code := cmd.Run([]string{"-delete-version-after=1s", basePath + "secret/my-secret"})
	if code != 0 {
		t.Fatalf("expected %d but received %d", 0, code)
	}

	metaFullPath := basePath + "metadata/secret/my-secret"
	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
	success := "Success! Data written to: " + metaFullPath
	if !strings.Contains(combined, success) {
		t.Fatalf("expected %q but received %q", success, combined)
	}

	secret, err := client.Logical().Read(metaFullPath)
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
	code = cmd.Run([]string{"-delete-version-after=0", basePath + "secret/my-secret"})
	if code != 0 {
		t.Errorf("expected %d but received %d", 0, code)
	}

	combined = ui.OutputWriter.String() + ui.ErrorWriter.String()
	if !strings.Contains(combined, success) {
		t.Errorf("expected %q but received %q", success, combined)
	}

	secret, err = client.Logical().Read(metaFullPath)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["delete_version_after"] != "0s" {
		t.Fatalf("expected 0s but received %q", secret.Data["delete_version_after"])
	}
}
