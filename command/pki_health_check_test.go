package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"

	"github.com/stretchr/testify/require"
)

func testPKIHealthCheckCommand(tb testing.TB) (*cli.MockUi, *PKIHealthCheckCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &PKIHealthCheckCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestPKIHC_Run(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
	}); err != nil {
		t.Fatalf("pki mount error: %#v", err)
	}

	if resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "Root X1",
		"ttl":         "876h",
	}); err != nil || resp == nil {
		t.Fatalf("failed to prime CA: %v", err)
	}

	if _, err := client.Logical().Read("pki/crl/rotate"); err != nil {
		t.Fatalf("failed to rotate CRLs: %v", err)
	}

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	code := RunCustom([]string{"pki", "health-check", "-format=json", "pki"}, runOpts)
	combined := stdout.String() + stderr.String()

	var results map[string][]map[string]interface{}
	if err := json.Unmarshal([]byte(combined), &results); err != nil {
		t.Fatalf("failed to decode json (ret %v): %v\njson:\n%v", code, err, combined)
	}

	t.Log(combined)

	expected := map[string][]map[string]interface{}{
		"ca_validity_period": {
			{
				"status": "critical",
			},
		},
		"crl_validity_period": {
			{
				"status": "ok",
			},
			{
				"status": "ok",
			},
		},
	}

	for test, subtest := range expected {
		actual, ok := results[test]
		require.True(t, ok, fmt.Sprintf("expected top-level test %v to be present", test))
		require.NotNil(t, actual, fmt.Sprintf("expected top-level test %v to be non-empty; wanted wireframe format %v", test, subtest))
		require.Equal(t, len(subtest), len(actual), fmt.Sprintf("top-level test %v has different number of results %v in wireframe, %v in test output\nwireframe: %v\noutput: %v\n", test, len(subtest), len(actual), subtest, actual))

		for index, subset := range subtest {
			for key, value := range subset {
				a_value, present := actual[index][key]
				require.True(t, present)
				if value != nil {
					require.Equal(t, value, a_value)
				}
			}
		}
	}
}
