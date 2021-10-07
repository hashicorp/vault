package command

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
)

func testKVPutCommand(tb testing.TB) (*cli.MockUi, *KVPutCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &KVPutCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func retryKVCommand(t *testing.T, client *api.Client, args []string) (code int, combined string) {
	t.Helper()

	// Loop until return message does not indicate upgrade, or timeout.
	timeout := time.After(20 * time.Second)
	for true {
		ui, cmd := testKVPutCommand(t)
		cmd.client = client
		code = cmd.Run(args)
		combined = ui.OutputWriter.String() + ui.ErrorWriter.String()

		// This is an error if a v1 mount, but test case case doesn't
		// currently contain the information to know the difference.
		if strings.Contains(combined, "Upgrading from non-versioned to versioned") {
			select {
			case <-timeout:
				t.Errorf("timeout expired waiting for upgrade: %q", combined)
				return code, combined
			default:
			}
			continue
		}
		break
	}
	return code, combined
}

func TestKVPutCommand(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"empty_kvs",
			[]string{"secret/write/foo"},
			"Must supply data",
			1,
		},
		{
			"kvs_no_value",
			[]string{"secret/write/foo", "foo"},
			"Failed to parse K=V data",
			1,
		},
		{
			"single_value",
			[]string{"secret/write/foo", "foo=bar"},
			"Success!",
			0,
		},
		{
			"multi_value",
			[]string{"secret/write/foo", "foo=bar", "zip=zap"},
			"Success!",
			0,
		},
		{
			"v2_single_value",
			[]string{"kv/write/foo", "foo=bar"},
			"created_time",
			0,
		},
		{
			"v2_multi_value",
			[]string{"kv/write/foo", "foo=bar", "zip=zap"},
			"created_time",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().Mount("kv/", &api.MountInput{
				Type: "kv-v2",
			}); err != nil {
				t.Fatal(err)
			}

			code, combined := retryKVCommand(t, client, tc.args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}
			if !strings.Contains(combined, tc.out) {
				t.Errorf("expected %q to contain %q", combined, tc.out)
			}
		})
	}

	t.Run("v2_cas", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		if err := client.Sys().Mount("kv/", &api.MountInput{
			Type: "kv-v2",
		}); err != nil {
			t.Fatal(err)
		}

		// Only have to potentially retry the first time.
		code, combined := retryKVCommand(t, client, []string{
			"-cas", "0", "kv/write/cas", "bar=baz",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}
		if !strings.Contains(combined, "created_time") {
			t.Errorf("expected %q to contain %q", combined, "created_time")
		}

		ui, cmd := testKVPutCommand(t)
		cmd.client = client
		code = cmd.Run([]string{
			"-cas", "1", "kv/write/cas", "bar=baz",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}
		combined = ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, "created_time") {
			t.Errorf("expected %q to contain %q", combined, "created_time")
		}

		ui, cmd = testKVPutCommand(t)
		cmd.client = client
		code = cmd.Run([]string{
			"-cas", "1", "kv/write/cas", "bar=baz",
		})
		if code != 2 {
			t.Fatalf("expected 2 to be %d", code)
		}
		combined = ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, "check-and-set parameter did not match the current version") {
			t.Errorf("expected %q to contain %q", combined, "check-and-set parameter did not match the current version")
		}
	})

	t.Run("v1_data", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testKVPutCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/write/data", "bar=baz",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, "Success!") {
			t.Errorf("expected %q to contain %q", combined, "created_time")
		}

		ui, rcmd := testReadCommand(t)
		rcmd.client = client
		code = rcmd.Run([]string{
			"secret/write/data",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}
		combined = ui.OutputWriter.String() + ui.ErrorWriter.String()
		if strings.Contains(combined, "data") {
			t.Errorf("expected %q not to contain %q", combined, "data")
		}
	})

	t.Run("stdin_full", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(`{"foo":"bar"}`))
			stdinW.Close()
		}()

		_, cmd := testKVPutCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"secret/write/stdin_full", "-",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}

		secret, err := client.Logical().Read("secret/write/stdin_full")
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil || secret.Data == nil {
			t.Fatal("expected secret to have data")
		}
		if exp, act := "bar", secret.Data["foo"].(string); exp != act {
			t.Errorf("expected %q to be %q", act, exp)
		}
	})

	t.Run("stdin_value", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte("bar"))
			stdinW.Close()
		}()

		_, cmd := testKVPutCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"secret/write/stdin_value", "foo=-",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}

		secret, err := client.Logical().Read("secret/write/stdin_value")
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil || secret.Data == nil {
			t.Fatal("expected secret to have data")
		}
		if exp, act := "bar", secret.Data["foo"].(string); exp != act {
			t.Errorf("expected %q to be %q", act, exp)
		}
	})

	t.Run("integration", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		_, cmd := testKVPutCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/write/integration", "foo=bar", "zip=zap",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}

		secret, err := client.Logical().Read("secret/write/integration")
		if err != nil {
			t.Fatal(err)
		}
		if secret == nil || secret.Data == nil {
			t.Fatal("expected secret to have data")
		}
		if exp, act := "bar", secret.Data["foo"].(string); exp != act {
			t.Errorf("expected %q to be %q", act, exp)
		}
		if exp, act := "zap", secret.Data["zip"].(string); exp != act {
			t.Errorf("expected %q to be %q", act, exp)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testKVPutCommand(t)
		assertNoTabs(t, cmd)
	})
}

func testKVGetCommand(tb testing.TB) (*cli.MockUi, *KVGetCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &KVGetCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestKVGetCommand(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar"},
			"Too many arguments",
			1,
		},
		{
			"not_found",
			[]string{"secret/nope/not/once/never"},
			"",
			2,
		},
		{
			"default",
			[]string{"secret/read/foo"},
			"foo",
			0,
		},
		{
			"v1_field",
			[]string{"-field", "foo", "secret/read/foo"},
			"bar",
			0,
		},
		{
			"v2_field",
			[]string{"-field", "foo", "kv/read/foo"},
			"bar",
			0,
		},

		{
			"v2_not_found",
			[]string{"kv/nope/not/once/never"},
			"",
			2,
		},

		{
			"v2_read",
			[]string{"kv/read/foo"},
			"foo",
			0,
		},
		{
			"v2_read",
			[]string{"kv/read/foo"},
			"version",
			0,
		},
		{
			"v2_read_version",
			[]string{"--version", "1", "kv/read/foo"},
			"foo",
			0,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()
				if err := client.Sys().Mount("kv/", &api.MountInput{
					Type: "kv-v2",
				}); err != nil {
					t.Fatal(err)
				}

				// Give time for the upgrade code to run/finish
				time.Sleep(time.Second)

				if _, err := client.Logical().Write("secret/read/foo", map[string]interface{}{
					"foo": "bar",
				}); err != nil {
					t.Fatal(err)
				}

				if _, err := client.Logical().Write("kv/data/read/foo", map[string]interface{}{
					"data": map[string]interface{}{
						"foo": "bar",
					},
				}); err != nil {
					t.Fatal(err)
				}

				ui, cmd := testKVGetCommand(t)
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
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testKVGetCommand(t)
		assertNoTabs(t, cmd)
	})
}

func testKVMetadataGetCommand(tb testing.TB) (*cli.MockUi, *KVMetadataGetCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &KVMetadataGetCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestKVMetadataGetCommand(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"v1",
			[]string{"secret/foo"},
			"Metadata not supported on KV Version 1",
			1,
		},
		{
			"metadata_exists",
			[]string{"kv/foo"},
			"current_version",
			0,
		},
		{
			"versions_exist",
			[]string{"kv/foo"},
			"deletion_time",
			0,
		},
	}

	t.Run("validations", func(t *testing.T) {
		t.Parallel()

		for _, tc := range cases {
			tc := tc

			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				client, closer := testVaultServer(t)
				defer closer()
				if err := client.Sys().Mount("kv/", &api.MountInput{
					Type: "kv-v2",
				}); err != nil {
					t.Fatal(err)
				}

				// Give time for the upgrade code to run/finish
				time.Sleep(time.Second)

				if _, err := client.Logical().Write("kv/data/foo", map[string]interface{}{
					"data": map[string]interface{}{
						"foo": "bar",
					},
				}); err != nil {
					t.Fatal(err)
				}

				ui, cmd := testKVMetadataGetCommand(t)
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
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testKVMetadataGetCommand(t)
		assertNoTabs(t, cmd)
	})
}

func testKVPatchCommand(tb testing.TB) (*cli.MockUi, *KVPatchCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &KVPatchCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestKVPatchCommand_ArgValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"not_enough_args",
			[]string{},
			"Not enough arguments",
			1,
		},
		{
			"empty_kvs",
			[]string{"kv/patch/foo"},
			"Must supply data",
			1,
		},
		{
			"kvs_no_value",
			[]string{"kv/patch/foo", "foo"},
			"Failed to parse K=V data",
			1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().Mount("kv/", &api.MountInput{
				Type: "kv-v2",
			}); err != nil {
				t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
			}

			ui, cmd := testKVPatchCommand(t)
			cmd.client = client

			code := cmd.Run(tc.args)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d for cmd %#v with args %#v\n", tc.code, code, cmd, tc.args)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

			if !strings.Contains(combined, tc.out) {
				t.Fatalf("expected output to be %q but was %q for cmd %#v with args %#v\n", tc.out, combined, cmd, tc.args)
			}
		})
	}
}

func TestKvPatchCommand_StdinFull(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	if _, err := client.Logical().Write("kv/data/patch/foo", map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "a",
		},
	}); err != nil {
		t.Fatalf("write failed, err: %#v\n", err)
	}

	stdinR, stdinW := io.Pipe()
	go func() {
		stdinW.Write([]byte(`{"foo":"bar"}`))
		stdinW.Close()
	}()

	_, cmd := testKVPatchCommand(t)
	cmd.client = client

	cmd.testStdin = stdinR

	args := []string{"kv/patch/foo", "-"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("expected code to be 0 but was %d for cmd %#v with args %#v\n", code, cmd, args)
	}

	secret, err := client.Logical().Read("kv/data/patch/foo")
	if err != nil {
		t.Fatalf("read failed, err: %#v\n", err)
	}

	if secret == nil || secret.Data == nil {
		t.Fatal("expected secret to have data")
	}

	secretDataRaw, ok := secret.Data["data"]

	if !ok {
		t.Fatalf("expected secret to have nested data key, data: %#v", secret.Data)
	}

	secretData := secretDataRaw.(map[string]interface{})
	foo, ok := secretData["foo"].(string)
	if !ok {
		t.Fatal("expected foo to be a string but it wasn't")
	}

	if exp, act := "bar", foo; exp != act {
		t.Fatalf("expected %q to be %q, data: %#v\n", act, exp, secret.Data)
	}
}

func TestKvPatchCommand_StdinValue(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	if _, err := client.Logical().Write("kv/data/patch/foo", map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "a",
		},
	}); err != nil {
		t.Fatalf("write failed, err: %#v\n", err)
	}

	stdinR, stdinW := io.Pipe()
	go func() {
		stdinW.Write([]byte("bar"))
		stdinW.Close()
	}()

	_, cmd := testKVPatchCommand(t)
	cmd.client = client

	cmd.testStdin = stdinR

	args := []string{"kv/patch/foo", "foo=-"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("expected code to be 0 but was %d for cmd %#v with args %#v\n", code, cmd, args)
	}

	secret, err := client.Logical().Read("kv/data/patch/foo")
	if err != nil {
		t.Fatalf("read failed, err: %#v\n", err)
	}

	if secret == nil || secret.Data == nil {
		t.Fatal("expected secret to have data")
	}

	secretDataRaw, ok := secret.Data["data"]

	if !ok {
		t.Fatalf("expected secret to have nested data key, data: %#v\n", secret.Data)
	}

	secretData := secretDataRaw.(map[string]interface{})

	if exp, act := "bar", secretData["foo"].(string); exp != act {
		t.Fatalf("expected %q to be %q, data: %#v\n", act, exp, secret.Data)
	}
}

func TestKVPatchCommand_RWMethodNotExists(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	ui, cmd := testKVPatchCommand(t)
	cmd.client = client

	args := []string{"-method", "rw", "kv/patch/foo", "foo=a"}
	code := cmd.Run(args)
	if code != 2 {
		t.Fatalf("expected code to be 2 but was %d for cmd %#v with args %#v\n", code, cmd, args)
	}

	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

	expectedOutputSubstr := "No value found"
	if !strings.Contains(combined, expectedOutputSubstr) {
		t.Fatalf("expected output %q to contain %q for cmd %#v with args %#v\n", combined, expectedOutputSubstr, cmd, args)
	}
}

func TestKVPatchCommand_RWMethodSucceeds(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	if _, err := client.Logical().Write("kv/data/patch/foo", map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "a",
			"bar": "b",
		},
	}); err != nil {
		t.Fatalf("write failed, err: %#v\n", err)
	}

	ui, cmd := testKVPatchCommand(t)
	cmd.client = client

	// Test single value
	args := []string{"-method", "rw", "kv/patch/foo", "foo=aa"}
	code := cmd.Run(args)
	if code != 0 {
		t.Fatalf("expected code to be 0 but was %d for cmd %#v with args %#v\n", code, cmd, args)
	}

	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

	expectedOutputSubstr := "created_time"
	if !strings.Contains(combined, expectedOutputSubstr) {
		t.Fatalf("expected output %q to contain %q for cmd %#v with args %#v\n", combined, expectedOutputSubstr, cmd, args)
	}

	// Test multi value
	ui, cmd = testKVPatchCommand(t)
	cmd.client = client

	args = []string{"-method", "rw", "kv/patch/foo", "foo=aaa", "bar=bbb"}
	code = cmd.Run(args)
	if code != 0 {
		t.Fatalf("expected code to be 0 but was %d for cmd %#v with args %#v\n", code, cmd, args)
	}

	combined = ui.OutputWriter.String() + ui.ErrorWriter.String()

	if !strings.Contains(combined, expectedOutputSubstr) {
		t.Fatalf("expected output %q to contain %q for cmd %#v with args %#v\n", combined, expectedOutputSubstr, cmd, args)
	}
}

func TestKVPatchCommand_CAS(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		expected string
		out      string
		code     int
	}{
		{
			"right version",
			[]string{"-cas", "1", "kv/foo", "bar=quux"},
			"quux",
			"",
			0,
		},
		{
			"wrong version",
			[]string{"-cas", "2", "kv/foo", "bar=wibble"},
			"baz",
			"check-and-set parameter did not match the current version",
			2,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().Mount("kv/", &api.MountInput{
				Type: "kv-v2",
			}); err != nil {
				t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
			}

			// create a policy with patch capability
			policy := `path "kv/*" { capabilities = ["create", "update", "read", "patch"] }`
			secretAuth, err := createTokenForPolicy(t, client, policy)
			if err != nil {
				t.Fatalf("policy/token creation failed for policy %s, err: %#v\n", policy, err)
			}

			kvClient, err := client.Clone()
			if err != nil {
				t.Fatal(err)
			}

			kvClient.SetToken(secretAuth.ClientToken)

			_, err = kvClient.Logical().Write("kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"bar": "baz"}})
			if err != nil {
				t.Fatal(err)
			}

			ui, cmd := testKVPatchCommand(t)
			cmd.client = kvClient

			code := cmd.Run(tc.args)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d", tc.code, code)
			}

			if tc.out != "" {
				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("expected %q to contain %q", combined, tc.out)
				}
			}

			secret, err := kvClient.Logical().Read("kv/data/foo")
			bar := secret.Data["data"].(map[string]interface{})["bar"]
			if bar != tc.expected {
				t.Fatalf("expected bar to be %q but it was %q", tc.expected, bar)
			}
		})
	}
}

func TestKVPatchCommand_Methods(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		expected string
		code     int
	}{
		{
			"rw",
			[]string{"-method", "rw", "kv/foo", "bar=quux"},
			"quux",
			0,
		},
		{
			"patch",
			[]string{"-method", "patch", "kv/foo", "bar=wibble"},
			"wibble",
			0,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().Mount("kv/", &api.MountInput{
				Type: "kv-v2",
			}); err != nil {
				t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
			}

			// create a policy with patch capability
			policy := `path "kv/*" { capabilities = ["create", "update", "read", "patch"] }`
			secretAuth, err := createTokenForPolicy(t, client, policy)
			if err != nil {
				t.Fatalf("policy/token creation failed for policy %s, err: %#v\n", policy, err)
			}

			kvClient, err := client.Clone()
			if err != nil {
				t.Fatal(err)
			}

			kvClient.SetToken(secretAuth.ClientToken)

			_, err = kvClient.Logical().Write("kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"bar": "baz"}})
			if err != nil {
				t.Fatal(err)
			}

			_, cmd := testKVPatchCommand(t)
			cmd.client = kvClient

			code := cmd.Run(tc.args)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d", tc.code, code)
			}

			secret, err := kvClient.Logical().Read("kv/data/foo")
			bar := secret.Data["data"].(map[string]interface{})["bar"]
			if bar != tc.expected {
				t.Fatalf("expected bar to be %q but it was %q", tc.expected, bar)
			}
		})
	}
}

func TestKVPatchCommand_403Fallback(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		expected string
		code     int
	}{
		// if no -method is specified, and patch fails, it should fall back to rw and succeed
		{
			"unspecified",
			[]string{"kv/foo", "bar=quux"},
			`add the "patch" capability to your ACL policy`,
			0,
		},
		// if -method=patch is specified, and patch fails, it should not fall back, and just error
		{
			"specifying patch",
			[]string{"-method", "patch", "kv/foo", "bar=quux"},
			"permission denied",
			2,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, closer := testVaultServer(t)
			defer closer()

			if err := client.Sys().Mount("kv/", &api.MountInput{
				Type: "kv-v2",
			}); err != nil {
				t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
			}

			// create a policy without patch capability
			policy := `path "kv/*" { capabilities = ["create", "update", "read"] }`
			secretAuth, err := createTokenForPolicy(t, client, policy)
			if err != nil {
				t.Fatalf("policy/token creation failed for policy %s, err: %#v\n", policy, err)
			}

			kvClient, err := client.Clone()
			if err != nil {
				t.Fatal(err)
			}

			kvClient.SetToken(secretAuth.ClientToken)

			// Write a value then attempt to patch it
			_, err = kvClient.Logical().Write("kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"bar": "baz"}})
			if err != nil {
				t.Fatal(err)
			}

			ui, cmd := testKVPatchCommand(t)
			cmd.client = kvClient
			code := cmd.Run(tc.args)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d", tc.code, code)
			}

			combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
			if !strings.Contains(combined, tc.expected) {
				t.Errorf("expected %q to contain %q", combined, tc.expected)
			}
		})
	}
}

func createTokenForPolicy(t *testing.T, client *api.Client, policy string) (*api.SecretAuth, error) {
	t.Helper()

	if err := client.Sys().PutPolicy("policy", policy); err != nil {
		return nil, err
	}

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"policy"},
		TTL:      "30m",
	})
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Auth == nil || secret.Auth.ClientToken == "" {
		return nil, fmt.Errorf("missing auth data: %#v", secret)
	}

	return secret.Auth, err
}

func TestKVPatchCommand_RWMethodNoRead(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	// The rw patch method requires the read capability - create a
	// policy that does not have it in order to force command to fail
	policy := `path "kv/data/patch/foo" { capabilities = ["create", "update"] }`
	secretAuth, err := createTokenForPolicy(t, client, policy)
	if err != nil {
		t.Fatalf("policy/token creation failed for policy %s, err: %#v\n", policy, err)
	}

	client.SetToken(secretAuth.ClientToken)

	if _, err := client.Logical().Write("kv/data/patch/foo", map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "a",
			"bar": "b",
		},
	}); err != nil {
		t.Fatalf("write failed, err: %#v\n", err)
	}

	ui, cmd := testKVPatchCommand(t)
	cmd.client = client

	args := []string{"kv/patch/foo", "foo=aa"}
	code := cmd.Run(args)
	expectedCode := 2

	if code != expectedCode {
		t.Fatalf("expected code to be %d but was %d for cmd %#v with args %#v\n", expectedCode, code, cmd, args)
	}

	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

	expectedOutputSubstr := "Error doing pre-read"
	if !strings.Contains(combined, expectedOutputSubstr) {
		t.Fatalf("expected output %q to contain %q for cmd %#v with args %#v\n", combined, expectedOutputSubstr, cmd, args)
	}
}

func TestKVPatchCommand_RWMethodNoUpdate(t *testing.T) {
	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	// The rw patch method requires the update capability - create a
	// policy that does not have it in order to force command to fail
	policy := `path "kv/data/patch/foo" { capabilities = ["create", "read"] }`
	secretAuth, err := createTokenForPolicy(t, client, policy)
	if err != nil {
		t.Fatalf("policy/token creation failed for policy %s, err: %#v\n", policy, err)
	}

	client.SetToken(secretAuth.ClientToken)

	if _, err := client.Logical().Write("kv/data/patch/foo", map[string]interface{}{
		"data": map[string]interface{}{
			"foo": "a",
			"bar": "b",
		},
	}); err != nil {
		t.Fatalf("write failed, err: %#v\n", err)
	}

	ui, cmd := testKVPatchCommand(t)
	cmd.client = client

	args := []string{"kv/patch/foo", "foo=aa"}
	code := cmd.Run(args)
	expectedCode := 2

	if code != expectedCode {
		t.Fatalf("expected code to be %d but was %d for cmd %#v with args %#v\n", expectedCode, code, cmd, args)
	}

	combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

	expectedOutputSubstr := "Error writing data"
	if !strings.Contains(combined, expectedOutputSubstr) {
		t.Fatalf("expected output %q to contain %q for cmd %#v with args %#v\n", combined, expectedOutputSubstr, cmd, args)
	}
}
