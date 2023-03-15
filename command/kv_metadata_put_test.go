// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-test/deep"
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

func TestKvMetadataPutCommand_DeleteVersionAfter(t *testing.T) {
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

func TestKvMetadataPutCommand_CustomMetadata(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	basePath := t.Name() + "/"
	secretPath := basePath + "secret/my-secret"

	if err := client.Sys().Mount(basePath, &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount error: %#v", err)
	}

	ui, cmd := testKVMetadataPutCommand(t)
	cmd.client = client

	exitStatus := cmd.Run([]string{"-custom-metadata=foo=abc", "-custom-metadata=bar=123", secretPath})

	if exitStatus != 0 {
		t.Fatalf("Expected 0 exit status but received %d", exitStatus)
	}

	metaFullPath := basePath + "metadata/secret/my-secret"
	commandOutput := ui.OutputWriter.String() + ui.ErrorWriter.String()
	expectedOutput := "Success! Data written to: " + metaFullPath

	if !strings.Contains(commandOutput, expectedOutput) {
		t.Fatalf("Expected command output %q but received %q", expectedOutput, commandOutput)
	}

	metadata, err := client.Logical().Read(metaFullPath)
	if err != nil {
		t.Fatalf("Metadata read error: %#v", err)
	}

	// JSON output from read decoded into map[string]interface{}
	expectedCustomMetadata := map[string]interface{}{
		"foo": "abc",
		"bar": "123",
	}

	if diff := deep.Equal(metadata.Data["custom_metadata"], expectedCustomMetadata); len(diff) > 0 {
		t.Fatal(diff)
	}

	ui, cmd = testKVMetadataPutCommand(t)
	cmd.client = client

	// Overwrite entire custom metadata with a single key
	exitStatus = cmd.Run([]string{"-custom-metadata=baz=abc123", secretPath})

	if exitStatus != 0 {
		t.Fatalf("Expected 0 exit status but received %d", exitStatus)
	}

	commandOutput = ui.OutputWriter.String() + ui.ErrorWriter.String()

	if !strings.Contains(commandOutput, expectedOutput) {
		t.Fatalf("Expected command output %q but received %q", expectedOutput, commandOutput)
	}

	metadata, err = client.Logical().Read(metaFullPath)

	if err != nil {
		t.Fatalf("Metadata read error: %#v", err)
	}

	expectedCustomMetadata = map[string]interface{}{
		"baz": "abc123",
	}

	if diff := deep.Equal(metadata.Data["custom_metadata"], expectedCustomMetadata); len(diff) > 0 {
		t.Fatal(diff)
	}
}

func TestKvMetadataPutCommand_UnprovidedFlags(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	basePath := t.Name() + "/"
	secretPath := basePath + "my-secret"

	if err := client.Sys().Mount(basePath, &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount error: %#v", err)
	}

	_, cmd := testKVMetadataPutCommand(t)
	cmd.client = client

	args := []string{"-cas-required=true", "-max-versions=10", secretPath}
	code, _ := kvMetadataPutWithRetry(t, client, args, nil)

	if code != 0 {
		t.Fatalf("expected 0 exit status but received %d", code)
	}

	args = []string{"-custom-metadata=foo=bar", secretPath}
	code, _ = kvMetadataPutWithRetry(t, client, args, nil)

	if code != 0 {
		t.Fatalf("expected 0 exit status but received %d", code)
	}

	secret, err := client.Logical().Read(basePath + "metadata/" + "my-secret")
	if err != nil {
		t.Fatal(err)
	}

	if secret.Data["cas_required"] != true {
		t.Fatalf("expected cas_required to be true but received %#v", secret.Data["cas_required"])
	}

	if secret.Data["max_versions"] != json.Number("10") {
		t.Fatalf("expected max_versions to be 10 but received %#v", secret.Data["max_versions"])
	}
}
