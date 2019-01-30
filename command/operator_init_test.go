// +build !race

package command

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func testOperatorInitCommand(tb testing.TB) (*cli.MockUi, *OperatorInitCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorInitCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorInitCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"too_many_args",
			[]string{"foo"},
			"Too many arguments",
			1,
		},
		{
			"pgp_keys_multi",
			[]string{
				"-pgp-keys", "keybase:hashicorp",
				"-pgp-keys", "keybase:jefferai",
			},
			"can only be specified once",
			1,
		},
		{
			"root_token_pgp_key_multi",
			[]string{
				"-root-token-pgp-key", "keybase:hashicorp",
				"-root-token-pgp-key", "keybase:jefferai",
			},
			"can only be specified once",
			1,
		},
		{
			"root_token_pgp_key_multi_inline",
			[]string{
				"-root-token-pgp-key", "keybase:hashicorp,keybase:jefferai",
			},
			"can only specify one pgp key",
			1,
		},
		{
			"recovery_pgp_keys_multi",
			[]string{
				"-recovery-pgp-keys", "keybase:hashicorp",
				"-recovery-pgp-keys", "keybase:jefferai",
			},
			"can only be specified once",
			1,
		},
		{
			"key_shares_pgp_less",
			[]string{
				"-key-shares", "10",
				"-pgp-keys", "keybase:jefferai,keybase:sethvargo",
			},
			"incorrect number",
			2,
		},
		{
			"key_shares_pgp_more",
			[]string{
				"-key-shares", "1",
				"-pgp-keys", "keybase:jefferai,keybase:sethvargo",
			},
			"incorrect number",
			2,
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

				ui, cmd := testOperatorInitCommand(t)
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

	t.Run("status", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerUninit(t)
		defer closer()

		ui, cmd := testOperatorInitCommand(t)
		cmd.client = client

		// Verify the non-init response code
		code := cmd.Run([]string{
			"-status",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		// Now init to verify the init response code
		if _, err := client.Sys().Init(&api.InitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		}); err != nil {
			t.Fatal(err)
		}

		// Verify the init response code
		ui, cmd = testOperatorInitCommand(t)
		cmd.client = client
		code = cmd.Run([]string{
			"-status",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}
	})

	t.Run("default", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerUninit(t)
		defer closer()

		ui, cmd := testOperatorInitCommand(t)
		cmd.client = client

		code := cmd.Run([]string{})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		init, err := client.Sys().InitStatus()
		if err != nil {
			t.Fatal(err)
		}
		if !init {
			t.Error("expected initialized")
		}

		re := regexp.MustCompile(`Unseal Key \d+: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < 5 || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		keys := make([]string, len(match))
		for i := range match {
			keys[i] = match[i][1]
		}

		// Try unsealing with those keys - only use 3, which is the default
		// threshold.
		for i, key := range keys[:3] {
			resp, err := client.Sys().Unseal(key)
			if err != nil {
				t.Fatal(err)
			}

			exp := (i + 1) % 3 // 1, 2, 0
			if resp.Progress != exp {
				t.Errorf("expected %d to be %d", resp.Progress, exp)
			}
		}

		status, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if status.Sealed {
			t.Errorf("expected vault to be unsealed: %#v", status)
		}
	})

	t.Run("custom_shares_threshold", func(t *testing.T) {
		t.Parallel()

		keyShares, keyThreshold := 20, 15

		client, closer := testVaultServerUninit(t)
		defer closer()

		ui, cmd := testOperatorInitCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-key-shares", strconv.Itoa(keyShares),
			"-key-threshold", strconv.Itoa(keyThreshold),
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		init, err := client.Sys().InitStatus()
		if err != nil {
			t.Fatal(err)
		}
		if !init {
			t.Error("expected initialized")
		}

		re := regexp.MustCompile(`Unseal Key \d+: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < keyShares || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		keys := make([]string, len(match))
		for i := range match {
			keys[i] = match[i][1]
		}

		// Try unsealing with those keys - only use 3, which is the default
		// threshold.
		for i, key := range keys[:keyThreshold] {
			resp, err := client.Sys().Unseal(key)
			if err != nil {
				t.Fatal(err)
			}

			exp := (i + 1) % keyThreshold
			if resp.Progress != exp {
				t.Errorf("expected %d to be %d", resp.Progress, exp)
			}
		}

		status, err := client.Sys().SealStatus()
		if err != nil {
			t.Fatal(err)
		}
		if status.Sealed {
			t.Errorf("expected vault to be unsealed: %#v", status)
		}
	})

	t.Run("pgp", func(t *testing.T) {
		t.Parallel()

		tempDir, pubFiles, err := getPubKeyFiles(t)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		client, closer := testVaultServerUninit(t)
		defer closer()

		ui, cmd := testOperatorInitCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-key-shares", "4",
			"-key-threshold", "2",
			"-pgp-keys", fmt.Sprintf("%s,@%s, %s,     %s           ",
				pubFiles[0], pubFiles[1], pubFiles[2], pubFiles[3]),
			"-root-token-pgp-key", pubFiles[0],
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		re := regexp.MustCompile(`Unseal Key \d+: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < 4 || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		keys := make([]string, len(match))
		for i := range match {
			keys[i] = match[i][1]
		}

		// Try unsealing with one key
		decryptedKey := testPGPDecrypt(t, pgpkeys.TestPrivKey1, keys[0])
		if _, err := client.Sys().Unseal(decryptedKey); err != nil {
			t.Fatal(err)
		}

		// Decrypt the root token
		reToken := regexp.MustCompile(`Root Token: (.+)`)
		match = reToken.FindAllStringSubmatch(output, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("no match")
		}
		root := match[0][1]
		decryptedRoot := testPGPDecrypt(t, pgpkeys.TestPrivKey1, root)

		if l, exp := len(decryptedRoot), vault.TokenLength+2; l != exp {
			t.Errorf("expected %d to be %d", l, exp)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testOperatorInitCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-key-shares=1",
			"-key-threshold=1",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error initializing: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testOperatorInitCommand(t)
		assertNoTabs(t, cmd)
	})
}
