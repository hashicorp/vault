// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/version"
)

func testVersionHistoryCommand(tb testing.TB) (*cli.MockUi, *clientcmd.VersionHistoryCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &clientcmd.VersionHistoryCommand{
		BaseCommand: &base.BaseCommand{
			UI: ui,
		},
	}
}

func TestVersionHistoryCommand_TableOutput(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	ui, cmd := testVersionHistoryCommand(t)
	cmd.SetClient(client)

	code := cmd.Run([]string{})

	if expectedCode := 0; code != expectedCode {
		t.Fatalf("expected %d to be %d: %s", code, expectedCode, ui.ErrorWriter.String())
	}

	if errorString := ui.ErrorWriter.String(); !strings.Contains(errorString, `Note:
Use of this command requires a server running Vault 1.10.0 or greater.
Version tracking was added in 1.9.0. Earlier versions have not been tracked.
`) {
		t.Errorf("expected %q to contain %q", errorString, `Note:
Use of this command requires a server running Vault 1.10.0 or greater.
Version tracking was added in 1.9.0. Earlier versions have not been tracked.
`)
	}

	output := ui.OutputWriter.String()

	if !strings.Contains(output, version.Version) {
		t.Errorf("expected %q to contain version %q", output, version.Version)
	}
}

func TestVersionHistoryCommand_JsonOutput(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &base.RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	args := []string{"version-history", "-format", "json"}
	format := "json"
	if format != "json" {
		t.Fatalf("expected format to be %q, actual %q", "json", format)
	}

	code := RunCustom(args, runOpts)

	if expectedCode := 0; code != expectedCode {
		t.Fatalf("expected %d to be %d: %s", code, expectedCode, stderr.String())
	}

	if stderrString := stderr.String(); !strings.Contains(stderrString, `Note:
Use of this command requires a server running Vault 1.10.0 or greater.
Version tracking was added in 1.9.0. Earlier versions have not been tracked.
`) {
		t.Errorf("expected %q to contain %q", stderrString, `Note:
Use of this command requires a server running Vault 1.10.0 or greater.
Version tracking was added in 1.9.0. Earlier versions have not been tracked.
`)
	}

	stdoutBytes := stdout.Bytes()

	if !json.Valid(stdoutBytes) {
		t.Fatalf("expected output %q to be valid JSON", stdoutBytes)
	}

	var versionHistoryResp map[string]interface{}
	err := json.Unmarshal(stdoutBytes, &versionHistoryResp)
	if err != nil {
		t.Fatalf("failed to unmarshal json from STDOUT, err: %s", err.Error())
	}

	var respData map[string]interface{}
	var ok bool
	var keys []interface{}
	var keyInfo map[string]interface{}

	if respData, ok = versionHistoryResp["data"].(map[string]interface{}); !ok {
		t.Fatalf("expected data key to be map, actual: %#v", versionHistoryResp["data"])
	}

	if keys, ok = respData["keys"].([]interface{}); !ok {
		t.Fatalf("expected keys to be array, actual: %#v", respData["keys"])
	}

	if keyInfo, ok = respData["key_info"].(map[string]interface{}); !ok {
		t.Fatalf("expected key_info to be map, actual: %#v", respData["key_info"])
	}

	if len(keys) != 1 {
		t.Fatalf("expected single version history entry for %q", version.Version)
	}

	if keyInfo[version.Version] == nil {
		t.Fatalf("expected version %s to be present in key_info, actual: %#v", version.Version, keyInfo)
	}
}
