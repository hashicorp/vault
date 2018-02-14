// +build !race

package command

import (
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/xor"
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
			"init_no_args",
			[]string{
				"-init",
			},
			"must specify either -otp or -pgp-key",
			1,
		},
		{
			"init_invalid_otp",
			[]string{
				"-init",
				"-otp", "not-a-valid-otp",
			},
			"Error initializing: invalid OTP:",
			1,
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

				ui, cmd := testOperatorGenerateRootCommand(t)

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

	t.Run("generate_otp", func(t *testing.T) {
		t.Parallel()

		ui, cmd := testOperatorGenerateRootCommand(t)

		code := cmd.Run([]string{
			"-generate-otp",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		output := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if err := cmd.verifyOTP(output); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("decode", func(t *testing.T) {
		t.Parallel()

		encoded := "L9MaZ/4mQanpOV6QeWd84g=="
		otp := "dIeeezkjpDUv3fy7MYPOLQ=="

		ui, cmd := testOperatorGenerateRootCommand(t)

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

		expected := "5b54841c-c705-e59c-c6e4-a22b48e4b2cf"
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if combined != expected {
			t.Errorf("expected %q to be %q", combined, expected)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		t.Parallel()

		otp := "dIeeezkjpDUv3fy7MYPOLQ=="

		client, closer := testVaultServer(t)
		defer closer()

		// Initialize a generation
		if _, err := client.Sys().GenerateRootInit(otp, ""); err != nil {
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

		otp := "dIeeezkjpDUv3fy7MYPOLQ=="

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorGenerateRootCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-init",
			"-otp", otp,
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

		otp := "dIeeezkjpDUv3fy7MYPOLQ=="

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Initialize a generation
		status, err := client.Sys().GenerateRootInit(otp, "")
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce

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

		reToken := regexp.MustCompile(`Root Token\s+(.+)`)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		match := reToken.FindAllStringSubmatch(combined, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		tokenBytes, err := xor.XORBase64(match[0][1], otp)
		if err != nil {
			t.Fatal(err)
		}
		token, err := uuid.FormatUUID(tokenBytes)
		if err != nil {
			t.Fatal(err)
		}

		if l, exp := len(token), 36; l != exp {
			t.Errorf("expected %d to be %d: %s", l, exp, token)
		}
	})

	t.Run("provide_stdin", func(t *testing.T) {
		t.Parallel()

		otp := "dIeeezkjpDUv3fy7MYPOLQ=="

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Initialize a generation
		status, err := client.Sys().GenerateRootInit(otp, "")
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce

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

		reToken := regexp.MustCompile(`Root Token\s+(.+)`)
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		match := reToken.FindAllStringSubmatch(combined, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("no match: %#v", match)
		}

		tokenBytes, err := xor.XORBase64(match[0][1], otp)
		if err != nil {
			t.Fatal(err)
		}
		token, err := uuid.FormatUUID(tokenBytes)
		if err != nil {
			t.Fatal(err)
		}

		if l, exp := len(token), 36; l != exp {
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
