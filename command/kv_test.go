package command

import (
	"context"
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

func retryKVCommand(t *testing.T, cmdFunc func() (int, string)) (int, string) {
	t.Helper()

	var code int
	var combined string

	// Loop until return message does not indicate upgrade, or timeout.
	timeout := time.After(20 * time.Second)
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-timeout:
			t.Errorf("timeout expired waiting for upgrade: %q", combined)
			return code, combined
		case <-ticker:
			code, combined = cmdFunc()

			// This is an error if a v1 mount, but test case doesn't
			// currently contain the information to know the difference.
			if !strings.Contains(combined, "Upgrading from non-versioned to versioned") {
				return code, combined
			}
		}
	}
}

func kvPutWithRetry(t *testing.T, client *api.Client, args []string) (int, string) {
	t.Helper()

	return retryKVCommand(t, func() (int, string) {
		ui, cmd := testKVPutCommand(t)
		cmd.client = client

		code := cmd.Run(args)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		return code, combined
	})
}

func kvPatchWithRetry(t *testing.T, client *api.Client, args []string, stdin *io.PipeReader) (int, string) {
	t.Helper()

	return retryKVCommand(t, func() (int, string) {
		ui, cmd := testKVPatchCommand(t)
		cmd.client = client

		if stdin != nil {
			cmd.testStdin = stdin
		}

		code := cmd.Run(args)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		return code, combined
	})
}

func TestKVPutCommand(t *testing.T) {
	t.Parallel()

	v2ExpectedFields := []string{"created_time", "custom_metadata", "deletion_time", "deletion_time", "version"}

	cases := []struct {
		name       string
		args       []string
		outStrings []string
		code       int
	}{
		{
			"not_enough_args",
			[]string{},
			[]string{"Not enough arguments"},
			1,
		},
		{
			"empty_kvs",
			[]string{"secret/write/foo"},
			[]string{"Must supply data"},
			1,
		},
		{
			"kvs_no_value",
			[]string{"secret/write/foo", "foo"},
			[]string{"Failed to parse K=V data"},
			1,
		},
		{
			"single_value",
			[]string{"secret/write/foo", "foo=bar"},
			[]string{"Success!"},
			0,
		},
		{
			"multi_value",
			[]string{"secret/write/foo", "foo=bar", "zip=zap"},
			[]string{"Success!"},
			0,
		},
		{
			"v1_mount_flag_syntax",
			[]string{"-mount", "secret", "write/foo", "foo=bar"},
			[]string{"Success!"},
			0,
		},
		{
			"v1_mount_flag_syntax_key_same_as_mount",
			[]string{"-mount", "secret", "secret", "foo=bar"},
			[]string{"Success!"},
			0,
		},
		{
			"v2_single_value",
			[]string{"kv/write/foo", "foo=bar"},
			v2ExpectedFields,
			0,
		},
		{
			"v2_multi_value",
			[]string{"kv/write/foo", "foo=bar", "zip=zap"},
			v2ExpectedFields,
			0,
		},
		{
			"v2_secret_path",
			[]string{"kv/write/foo", "foo=bar"},
			[]string{"== Secret Path ==", "kv/data/write/foo"},
			0,
		},
		{
			"v2_mount_flag_syntax",
			[]string{"-mount", "kv", "write/foo", "foo=bar"},
			v2ExpectedFields,
			0,
		},
		{
			"v2_mount_flag_syntax_key_same_as_mount",
			[]string{"-mount", "kv", "kv", "foo=bar"},
			v2ExpectedFields,
			0,
		},
		{
			"v2_single_value_backslash",
			[]string{"kv/write/foo", "foo=\\"},
			[]string{"== Secret Path ==", "kv/data/write/foo"},
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

			code, combined := kvPutWithRetry(t, client, tc.args)
			if code != tc.code {
				t.Errorf("expected %d to be %d", code, tc.code)
			}

			for _, str := range tc.outStrings {
				if !strings.Contains(combined, str) {
					t.Errorf("expected %q to contain %q", combined, str)
				}
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
		code, combined := kvPutWithRetry(t, client, []string{
			"-cas", "0", "kv/write/cas", "bar=baz",
		})
		if code != 0 {
			t.Fatalf("expected 0 to be %d", code)
		}

		for _, str := range v2ExpectedFields {
			if !strings.Contains(combined, str) {
				t.Errorf("expected %q to contain %q", combined, str)
			}
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

		for _, str := range v2ExpectedFields {
			if !strings.Contains(combined, str) {
				t.Errorf("expected %q to contain %q", combined, str)
			}
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

	baseV2ExpectedFields := []string{"created_time", "custom_metadata", "deletion_time", "deletion_time", "version"}

	cases := []struct {
		name       string
		args       []string
		outStrings []string
		code       int
	}{
		{
			"not_enough_args",
			[]string{},
			[]string{"Not enough arguments"},
			1,
		},
		{
			"too_many_args",
			[]string{"foo", "bar"},
			[]string{"Too many arguments"},
			1,
		},
		{
			"not_found",
			[]string{"secret/nope/not/once/never"},
			[]string{"No value found at secret/nope/not/once/never"},
			2,
		},
		{
			"default",
			[]string{"secret/read/foo"},
			[]string{"foo"},
			0,
		},
		{
			"v1_field",
			[]string{"-field", "foo", "secret/read/foo"},
			[]string{"bar"},
			0,
		},
		{
			"v1_mount_flag_syntax",
			[]string{"-mount", "secret", "read/foo"},
			[]string{"foo"},
			0,
		},
		{
			"v2_field",
			[]string{"-field", "foo", "kv/read/foo"},
			[]string{"bar"},
			0,
		},
		{
			"v2_mount_flag_syntax",
			[]string{"-mount", "kv", "read/foo"},
			append(baseV2ExpectedFields, "foo"),
			0,
		},
		{
			"v2_mount_flag_syntax_leading_slash",
			[]string{"-mount", "kv", "/read/foo"},
			append(baseV2ExpectedFields, "foo"),
			0,
		},
		{
			"v1_mount_flag_syntax_key_same_as_mount",
			[]string{"-mount", "kv", "kv"},
			append(baseV2ExpectedFields, "foo"),
			0,
		},
		{
			"v2_mount_flag_syntax_key_same_as_mount",
			[]string{"-mount", "kv", "kv"},
			append(baseV2ExpectedFields, "foo"),
			0,
		},
		{
			"v2_not_found",
			[]string{"kv/nope/not/once/never"},
			[]string{"No value found at kv/data/nope/not/once/never"},
			2,
		},
		{
			"v2_read",
			[]string{"kv/read/foo"},
			append(baseV2ExpectedFields, "foo"),
			0,
		},
		{
			"v2_read_leading_slash",
			[]string{"/kv/read/foo"},
			append(baseV2ExpectedFields, "foo"),
			0,
		},
		{
			"v2_read_version",
			[]string{"--version", "1", "kv/read/foo"},
			append(baseV2ExpectedFields, "foo"),
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

				// create KV entries to test -mount flag where secret key is same as mount path
				if _, err := client.Logical().Write("secret/secret", map[string]interface{}{
					"foo": "bar",
				}); err != nil {
					t.Fatal(err)
				}

				if _, err := client.Logical().Write("kv/data/kv", map[string]interface{}{
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

				for _, str := range tc.outStrings {
					if !strings.Contains(combined, str) {
						t.Errorf("expected %q to contain %q", combined, str)
					}
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

	expectedTopLevelFields := []string{
		"cas_required",
		"created_time",
		"current_version",
		"custom_metadata",
		"delete_version_after",
		"max_versions",
		"oldest_version",
		"updated_time",
	}

	expectedVersionFields := []string{
		"created_time", // field is redundant
		"deletion_time",
		"destroyed",
	}

	cases := []struct {
		name       string
		args       []string
		outStrings []string
		code       int
	}{
		{
			"v1",
			[]string{"secret/foo"},
			[]string{"Metadata not supported on KV Version 1"},
			1,
		},
		{
			"metadata_exists",
			[]string{"kv/foo"},
			expectedTopLevelFields,
			0,
		},
		// ensure that all top-level and version-level fields are output along with version num
		{
			"versions_exist",
			[]string{"kv/foo"},
			append(expectedTopLevelFields, expectedVersionFields[:]...),
			0,
		},
		{
			"mount_flag_syntax",
			[]string{"-mount", "kv", "foo"},
			expectedTopLevelFields,
			0,
		},
		{
			"mount_flag_syntax_key_same_as_mount",
			[]string{"-mount", "kv", "kv"},
			expectedTopLevelFields,
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

				// create KV entry to test -mount flag where secret key is same as mount path
				if _, err := client.Logical().Write("kv/data/kv", map[string]interface{}{
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
				for _, str := range tc.outStrings {
					if !strings.Contains(combined, str) {
						t.Errorf("expected %q to contain %q", combined, str)
					}
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
		{
			"mount_flag_syntax",
			[]string{"-mount", "kv"},
			"Not enough arguments",
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

			code, combined := kvPatchWithRetry(t, client, tc.args, nil)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d for patch cmd with args %#v\n", tc.code, code, tc.args)
			}

			if !strings.Contains(combined, tc.out) {
				t.Fatalf("expected output to be %q but was %q for patch cmd with args %#v\n", tc.out, combined, tc.args)
			}
		})
	}
}

// expectedPatchFields produces a deterministic slice of
// expected fields for patch command output since const
// slices are not supported
func expectedPatchFields() []string {
	return []string{
		"created_time",
		"custom_metadata",
		"deletion_time",
		"destroyed",
		"version",
	}
}

func TestKVPatchCommand_StdinFull(t *testing.T) {
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

	cases := [][]string{
		{"kv/patch/foo", "-"},
		{"-mount", "kv", "patch/foo", "-"},
	}
	for i, args := range cases {
		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i)))
			stdinW.Close()
		}()
		code, combined := kvPatchWithRetry(t, client, args, stdinR)

		for _, str := range expectedPatchFields() {
			if !strings.Contains(combined, str) {
				t.Errorf("expected %q to contain %q", combined, str)
			}
		}

		if code != 0 {
			t.Fatalf("expected code to be 0 but was %d for patch cmd with args %#v\n", code, args)
		}

		secret, err := client.Logical().ReadWithContext(context.Background(), "kv/data/patch/foo")
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
		foo, ok := secretData[fmt.Sprintf("foo%d", i)].(string)
		if !ok {
			t.Fatal("expected foo to be a string but it wasn't")
		}

		if exp, act := fmt.Sprintf("bar%d", i), foo; exp != act {
			t.Fatalf("expected %q to be %q, data: %#v\n", act, exp, secret.Data)
		}
	}
}

func TestKVPatchCommand_StdinValue(t *testing.T) {
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

	cases := [][]string{
		{"kv/patch/foo", "foo=-"},
		{"-mount", "kv", "patch/foo", "foo=-"},
	}

	for i, args := range cases {
		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(fmt.Sprintf("bar%d", i)))
			stdinW.Close()
		}()

		code, combined := kvPatchWithRetry(t, client, args, stdinR)
		if code != 0 {
			t.Fatalf("expected code to be 0 but was %d for patch cmd with args %#v\n", code, args)
		}

		for _, str := range expectedPatchFields() {
			if !strings.Contains(combined, str) {
				t.Errorf("expected %q to contain %q", combined, str)
			}
		}

		secret, err := client.Logical().ReadWithContext(context.Background(), "kv/data/patch/foo")
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

		if exp, act := fmt.Sprintf("bar%d", i), secretData["foo"].(string); exp != act {
			t.Fatalf("expected %q to be %q, data: %#v\n", act, exp, secret.Data)
		}
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

	cases := [][]string{
		{"-method", "rw", "kv/patch/foo", "foo=a"},
		{"-method", "rw", "-mount", "kv", "patch/foo", "foo=a"},
	}

	for _, args := range cases {
		code, combined := kvPatchWithRetry(t, client, args, nil)

		if code != 2 {
			t.Fatalf("expected code to be 2 but was %d for patch cmd with args %#v\n", code, args)
		}

		expectedOutputSubstr := "No value found"
		if !strings.Contains(combined, expectedOutputSubstr) {
			t.Fatalf("expected output %q to contain %q for patch cmd with args %#v\n", combined, expectedOutputSubstr, args)
		}
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

	// Test single value
	args := []string{"-method", "rw", "kv/patch/foo", "foo=aa"}
	code, combined := kvPatchWithRetry(t, client, args, nil)

	if code != 0 {
		t.Fatalf("expected code to be 0 but was %d for patch cmd with args %#v\n", code, args)
	}

	for _, str := range expectedPatchFields() {
		if !strings.Contains(combined, str) {
			t.Errorf("expected %q to contain %q", combined, str)
		}
	}

	// Test that full path was output
	for _, str := range []string{"== Secret Path ==", "kv/data/patch/foo"} {
		if !strings.Contains(combined, str) {
			t.Errorf("expected %q to contain %q", combined, str)
		}
	}

	// Test multi value
	args = []string{"-method", "rw", "kv/patch/foo", "foo=aaa", "bar=bbb"}
	code, combined = kvPatchWithRetry(t, client, args, nil)

	if code != 0 {
		t.Fatalf("expected code to be 0 but was %d for patch cmd with args %#v\n", code, args)
	}

	for _, str := range expectedPatchFields() {
		if !strings.Contains(combined, str) {
			t.Errorf("expected %q to contain %q", combined, str)
		}
	}
}

func TestKVPatchCommand_CAS(t *testing.T) {
	cases := []struct {
		name       string
		key        string
		args       []string
		expected   string
		outStrings []string
		code       int
	}{
		{
			"right version",
			"foo",
			[]string{"-cas", "1", "kv/foo", "bar=quux"},
			"quux",
			expectedPatchFields(),
			0,
		},
		{
			"wrong version",
			"foo",
			[]string{"-cas", "2", "kv/foo", "bar=wibble"},
			"baz",
			[]string{"check-and-set parameter did not match the current version"},
			2,
		},
		{
			"mount_flag_syntax",
			"foo",
			[]string{"-mount", "kv", "-cas", "1", "foo", "bar=quux"},
			"quux",
			expectedPatchFields(),
			0,
		},
		{
			"v2_mount_flag_syntax_key_same_as_mount",
			"kv",
			[]string{"-mount", "kv", "-cas", "1", "kv", "bar=quux"},
			"quux",
			expectedPatchFields(),
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

			data := map[string]interface{}{
				"bar": "baz",
			}

			_, err = kvClient.Logical().Write("kv/data/"+tc.key, map[string]interface{}{"data": data})
			if err != nil {
				t.Fatal(err)
			}

			code, combined := kvPatchWithRetry(t, kvClient, tc.args, nil)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d", tc.code, code)
			}

			for _, str := range tc.outStrings {
				if !strings.Contains(combined, str) {
					t.Errorf("expected %q to contain %q", combined, str)
				}
			}

			secret, err := kvClient.Logical().ReadWithContext(context.Background(), "kv/data/"+tc.key)
			if err != nil {
				t.Fatal(err)
			}
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

			code, _ := kvPatchWithRetry(t, kvClient, tc.args, nil)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d", tc.code, code)
			}

			secret, err := kvClient.Logical().ReadWithContext(context.Background(), "kv/data/foo")
			if err != nil {
				t.Fatal(err)
			}
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

			code, combined := kvPatchWithRetry(t, kvClient, tc.args, nil)

			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d", tc.code, code)
			}

			if !strings.Contains(combined, tc.expected) {
				t.Errorf("expected %q to contain %q", combined, tc.expected)
			}
		})
	}
}

func TestKVPatchCommand_RWMethodPolicyVariations(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		policy   string
		expected string
		code     int
	}{
		// if the policy doesn't have read capability and -method=rw is specified, it fails
		{
			"no read",
			[]string{"-method", "rw", "kv/foo", "bar=quux"},
			`path "kv/*" { capabilities = ["create", "update"] }`,
			"permission denied",
			2,
		},
		// if the policy doesn't have update capability and -method=rw is specified, it fails
		{
			"no update",
			[]string{"-method", "rw", "kv/foo", "bar=quux"},
			`path "kv/*" { capabilities = ["create", "read"] }`,
			"permission denied",
			2,
		},
		// if the policy has both read and update and -method=rw is specified, it succeeds
		{
			"read and update",
			[]string{"-method", "rw", "kv/foo", "bar=quux"},
			`path "kv/*" { capabilities = ["create", "read", "update"] }`,
			"",
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

			secretAuth, err := createTokenForPolicy(t, client, tc.policy)
			if err != nil {
				t.Fatalf("policy/token creation failed for policy %s, err: %#v\n", tc.policy, err)
			}

			client.SetToken(secretAuth.ClientToken)

			putArgs := []string{"kv/foo", "foo=bar", "bar=baz"}
			code, combined := kvPutWithRetry(t, client, putArgs)
			if code != 0 {
				t.Errorf("write failed, expected %d to be 0, output: %s", code, combined)
			}

			code, combined = kvPatchWithRetry(t, client, tc.args, nil)
			if code != tc.code {
				t.Fatalf("expected code to be %d but was %d for patch cmd with args %#v\n", tc.code, code, tc.args)
			}

			if code != 0 {
				if !strings.Contains(combined, tc.expected) {
					t.Fatalf("expected output %q to contain %q for patch cmd with args %#v\n", combined, tc.expected, tc.args)
				}
			}
		})
	}
}

func TestPadEqualSigns(t *testing.T) {
	t.Parallel()

	header := "Test Header"

	cases := []struct {
		name          string
		totalPathLen  int
		expectedCount int
	}{
		{
			name:          "path with even length",
			totalPathLen:  20,
			expectedCount: 4,
		},
		{
			name:          "path with odd length",
			totalPathLen:  19,
			expectedCount: 3,
		},
		{
			name:          "smallest possible path",
			totalPathLen:  8,
			expectedCount: 2,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			padded := padEqualSigns(header, tc.totalPathLen)

			signs := strings.Split(padded, fmt.Sprintf(" %s ", header))
			if len(signs[0]) != len(signs[1]) {
				t.Fatalf("expected an equal number of equal signs on both sides")
			}
			for _, sign := range signs {
				count := strings.Count(sign, "=")
				if count != tc.expectedCount {
					t.Fatalf("expected %d equal signs but there were %d", tc.expectedCount, count)
				}
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
