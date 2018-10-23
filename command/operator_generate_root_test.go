// +build !race

package command

import (
	"encoding/base64"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func testOperatorGenerateRootCommand(tb testing.TB) (*cli.MockUi, *OperatorGenerateRootCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorGenerateRootCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorGenerateRootCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"init_invalid_otp",
			[]string{
				"-init",
				"-otp", "not-a-valid-otp",
			},
			"OTP string is wrong length",
			2,
		},
		{
			"init_pgp_multi",
			[]string{
				"-init",
				"-pgp-key", "keybase:hashicorp",
				"-pgp-key", "keybase:jefferai",
			},
			"can only be specified once",
			1,
		},
		{
			"init_pgp_multi_inline",
			[]string{
				"-init",
				"-pgp-key", "keybase:hashicorp,keybase:jefferai",
			},
			"can only specify one pgp key",
			1,
		},
		{
			"init_pgp_otp",
			[]string{
				"-init",
				"-pgp-key", "keybase:hashicorp",
				"-otp", "abcd1234",
			},
			"cannot specify both -otp and -pgp-key",
			1,
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

				ui, cmd := testOperatorGenerateRootCommand(t)
				cmd.client = client

				code := cmd.Run(tc.args)
				if code != tc.code {
					t.Errorf("%s: expected %d to be %d", tc.name, code, tc.code)
				}

				combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
				if !strings.Contains(combined, tc.out) {
					t.Errorf("%s: expected %q to contain %q", tc.name, combined, tc.out)
				}
			})
		}
	})

	t.Run("generate_otp", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		_, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-generate-otp",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}
	})

	t.Run("decode", func(t *testing.T) {
		t.Parallel()

		encoded := "Bxg9JQQqOCNKBRICNwMIRzo2J3cWCBRi"
		otp := "3JhHkONiyiaNYj14nnD9xZQS"

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		// Simulate piped output to print raw output
		old := os.Stdout
		_, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		os.Stdout = w

		code := cmd.Run([]string{
			"-decode", encoded,
			"-otp", otp,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		w.Close()
		os.Stdout = old

		expected := "4RUmoevJ3lsLni9sTXcNnRE1"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if combined != expected {
			t.Errorf("expected %q to be %q", combined, expected)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		// Initialize a generation
		if _, err := client.Sys().GenerateRootInit("", ""); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-cancel",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Root token generation canceled"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		status, err := client.Sys().GenerateRootStatus()
		if err != nil {
			t.Fatal(err)
		}

		if status.Started {
			t.Errorf("expected status to be canceled: %#v", status)
		}
	})

	t.Run("init_otp", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-init",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Nonce"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		status, err := client.Sys().GenerateRootStatus()
		if err != nil {
			t.Fatal(err)
		}

		if !status.Started {
			t.Errorf("expected status to be started: %#v", status)
		}
	})

	t.Run("init_pgp", func(t *testing.T) {
		t.Parallel()

		pgpKey := "keybase:hashicorp"
		pgpFingerprint := "91a6e7f85d05c65630bef18951852d87348ffc4c"

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-init",
			"-pgp-key", pgpKey,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Nonce"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		status, err := client.Sys().GenerateRootStatus()
		if err != nil {
			t.Fatal(err)
		}

		if !status.Started {
			t.Errorf("expected status to be started: %#v", status)
		}
		if status.PGPFingerprint != pgpFingerprint {
			t.Errorf("expected %q to be %q", status.PGPFingerprint, pgpFingerprint)
		}
	})

	t.Run("status", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-status",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Nonce"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("provide_arg", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Initialize a generation
		status, err := client.Sys().GenerateRootInit("", "")
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce
		otp := status.OTP

		// Supply the first n-1 unseal keys
		for _, key := range keys[:len(keys)-1] {
			_, cmd := testOperatorGenerateRootCommand(t)
			cmd.client = client

			code := cmd.Run([]string{
				"-nonce", nonce,
				key,
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d", code, exp)
			}
		}

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-nonce", nonce,
			keys[len(keys)-1], // the last unseal key
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		reToken := regexp.MustCompile(`Encoded Token\s+(.+)`)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		match := reToken.FindAllStringSubmatch(combined, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		tokenBytes, err := base64.RawStdEncoding.DecodeString(match[0][1])
		if err != nil {
			t.Fatal(err)
		}

		token, err := xor.XORBytes(tokenBytes, []byte(otp))
		if err != nil {
			t.Fatal(err)
		}

		if l, exp := len(token), vault.TokenLength+2; l != exp {
			t.Errorf("expected %d to be %d: %s", l, exp, token)
		}
	})

	t.Run("provide_stdin", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Initialize a generation
		status, err := client.Sys().GenerateRootInit("", "")
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce
		otp := status.OTP

		// Supply the first n-1 unseal keys
		for _, key := range keys[:len(keys)-1] {
			stdinR, stdinW := io.Pipe()
			go func() {
				stdinW.Write([]byte(key))
				stdinW.Close()
			}()

			_, cmd := testOperatorGenerateRootCommand(t)
			cmd.client = client
			cmd.testStdin = stdinR

			code := cmd.Run([]string{
				"-nonce", nonce,
				"-",
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d", code, exp)
			}
		}

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(keys[len(keys)-1])) // the last unseal key
			stdinW.Close()
		}()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"-nonce", nonce,
			"-",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		reToken := regexp.MustCompile(`Encoded Token\s+(.+)`)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		match := reToken.FindAllStringSubmatch(combined, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		// encodedOTP := base64.RawStdEncoding.EncodeToString([]byte(otp))

		// tokenBytes, err := xor.XORBase64(match[0][1], encodedOTP)
		// if err != nil {
		// 	t.Fatal(err)
		// }
		// token, err := uuid.FormatUUID(tokenBytes)
		// if err != nil {
		// 	t.Fatal(err)
		// }

		tokenBytes, err := base64.RawStdEncoding.DecodeString(match[0][1])
		if err != nil {
			t.Fatal(err)
		}

		token, err := xor.XORBytes(tokenBytes, []byte(otp))
		if err != nil {
			t.Fatal(err)
		}

		if l, exp := len(token), vault.TokenLength+2; l != exp {
			t.Errorf("expected %d to be %d: %s", l, exp, token)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/foo",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error getting root generation status: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testOperatorGenerateRootCommand(t)
		assertNoTabs(t, cmd)
	})
}
