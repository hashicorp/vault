// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !race

package command

import (
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/roottoken"
)

func testOperatorRekeyCommand(tb testing.TB) (*cli.MockUi, *OperatorRekeyCommand) {
	tb.Helper()

	ui := cli.NewMockUi()
	return ui, &OperatorRekeyCommand{
		BaseCommand: &BaseCommand{
			UI: ui,
		},
	}
}

func TestOperatorRekeyCommand_Run(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		args []string
		out  string
		code int
	}{
		{
			"pgp_keys_multi",
			[]string{
				"-init",
				"-pgp-keys", "keybase:hashicorp",
				"-pgp-keys", "keybase:jefferai",
			},
			"can only be specified once",
			1,
		},
		{
			"key_shares_pgp_less",
			[]string{
				"-init",
				"-key-shares", "10",
				"-pgp-keys", "keybase:jefferai,keybase:sethvargo",
			},
			"incorrect number",
			2,
		},
		{
			"key_shares_pgp_more",
			[]string{
				"-init",
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

				ui, cmd := testOperatorRekeyCommand(t)
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

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		// Verify the non-init response
		code := cmd.Run([]string{
			"-status",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		expected := "Nonce"
		combined := ui.OutputWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		// Now init to verify the init response
		if _, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		}); err != nil {
			t.Fatal(err)
		}

		// Verify the init response
		ui, cmd = testOperatorRekeyCommand(t)
		cmd.client = client
		code = cmd.Run([]string{
			"-status",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		expected = "Progress"
		combined = ui.OutputWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		// Initialize a rekey
		if _, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		}); err != nil {
			t.Fatal(err)
		}

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-cancel",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Success! Canceled rekeying"
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

	t.Run("init", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-init",
			"-key-shares", "1",
			"-key-threshold", "1",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		expected := "Nonce"
		combined := ui.OutputWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		status, err := client.Sys().RekeyStatus()
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
		pgpFingerprints := []string{"c874011f0ab405110d02105534365d9472d7468f"}

		client, closer := testVaultServer(t)
		defer closer()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-init",
			"-key-shares", "1",
			"-key-threshold", "1",
			"-pgp-keys", pgpKey,
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		expected := "Nonce"
		combined := ui.OutputWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}

		status, err := client.Sys().RekeyStatus()
		if err != nil {
			t.Fatal(err)
		}
		if !status.Started {
			t.Errorf("expected status to be started: %#v", status)
		}
		if !reflect.DeepEqual(status.PGPFingerprints, pgpFingerprints) {
			t.Errorf("expected %#v to be %#v", status.PGPFingerprints, pgpFingerprints)
		}
	})

	t.Run("provide_arg_recovery_keys", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerAutoUnseal(t)
		defer closer()

		// Initialize a rekey
		status, err := client.Sys().RekeyRecoveryKeyInit(&api.RekeyInitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		})
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce

		// Supply the first n-1 recovery keys
		for _, key := range keys[:len(keys)-1] {
			ui, cmd := testOperatorRekeyCommand(t)
			cmd.client = client

			code := cmd.Run([]string{
				"-nonce", nonce,
				"-target", "recovery",
				key,
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
			}
		}

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-nonce", nonce,
			"-target", "recovery",
			keys[len(keys)-1], // the last recovery key
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		re := regexp.MustCompile(`Key 1: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("bad match: %#v", match)
		}
		recoveryKey := match[0][1]

		if strings.Contains(strings.ToLower(output), "unseal key") {
			t.Fatalf(`output %s shouldn't contain "unseal key"`, output)
		}

		// verify that we can perform operations with the recovery key
		// below we generate a root token using the recovery key
		rootStatus, err := client.Sys().GenerateRootStatus()
		if err != nil {
			t.Fatal(err)
		}
		otp, err := roottoken.GenerateOTP(rootStatus.OTPLength)
		if err != nil {
			t.Fatal(err)
		}
		genRoot, err := client.Sys().GenerateRootInit(otp, "")
		if err != nil {
			t.Fatal(err)
		}
		r, err := client.Sys().GenerateRootUpdate(recoveryKey, genRoot.Nonce)
		if err != nil {
			t.Fatal(err)
		}
		if !r.Complete {
			t.Fatal("expected root update to be complete")
		}
	})
	t.Run("provide_arg", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Initialize a rekey
		status, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		})
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce

		// Supply the first n-1 unseal keys
		for _, key := range keys[:len(keys)-1] {
			ui, cmd := testOperatorRekeyCommand(t)
			cmd.client = client

			code := cmd.Run([]string{
				"-nonce", nonce,
				key,
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
			}
		}

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-nonce", nonce,
			keys[len(keys)-1], // the last unseal key
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		re := regexp.MustCompile(`Key 1: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("bad match: %#v", match)
		}

		// Grab the unseal key and try to unseal
		unsealKey := match[0][1]
		if err := client.Sys().Seal(); err != nil {
			t.Fatal(err)
		}
		sealStatus, err := client.Sys().Unseal(unsealKey)
		if err != nil {
			t.Fatal(err)
		}
		if sealStatus.Sealed {
			t.Errorf("expected vault to be unsealed: %#v", sealStatus)
		}
	})

	t.Run("provide_stdin", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		// Initialize a rekey
		status, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		})
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

			ui, cmd := testOperatorRekeyCommand(t)
			cmd.client = client
			cmd.testStdin = stdinR

			code := cmd.Run([]string{
				"-nonce", nonce,
				"-",
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
			}
		}

		stdinR, stdinW := io.Pipe()
		go func() {
			stdinW.Write([]byte(keys[len(keys)-1])) // the last unseal key
			stdinW.Close()
		}()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"-nonce", nonce,
			"-",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		re := regexp.MustCompile(`Key 1: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("bad match: %#v", match)
		}

		// Grab the unseal key and try to unseal
		unsealKey := match[0][1]
		if err := client.Sys().Seal(); err != nil {
			t.Fatal(err)
		}
		sealStatus, err := client.Sys().Unseal(unsealKey)
		if err != nil {
			t.Fatal(err)
		}
		if sealStatus.Sealed {
			t.Errorf("expected vault to be unsealed: %#v", sealStatus)
		}
	})

	t.Run("provide_stdin_recovery_keys", func(t *testing.T) {
		t.Parallel()

		client, keys, closer := testVaultServerAutoUnseal(t)
		defer closer()

		// Initialize a rekey
		status, err := client.Sys().RekeyRecoveryKeyInit(&api.RekeyInitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		})
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce
		for _, key := range keys[:len(keys)-1] {
			stdinR, stdinW := io.Pipe()
			go func() {
				_, _ = stdinW.Write([]byte(key))
				_ = stdinW.Close()
			}()

			ui, cmd := testOperatorRekeyCommand(t)
			cmd.client = client
			cmd.testStdin = stdinR

			code := cmd.Run([]string{
				"-target", "recovery",
				"-nonce", nonce,
				"-",
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
			}
		}

		stdinR, stdinW := io.Pipe()
		go func() {
			_, _ = stdinW.Write([]byte(keys[len(keys)-1])) // the last recovery key
			_ = stdinW.Close()
		}()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client
		cmd.testStdin = stdinR

		code := cmd.Run([]string{
			"-nonce", nonce,
			"-target", "recovery",
			"-",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		re := regexp.MustCompile(`Key 1: (.+)`)
		output := ui.OutputWriter.String()
		match := re.FindAllStringSubmatch(output, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			t.Fatalf("bad match: %#v", match)
		}
		recoveryKey := match[0][1]

		if strings.Contains(strings.ToLower(output), "unseal key") {
			t.Fatalf(`output %s shouldn't contain "unseal key"`, output)
		}
		// verify that we can perform operations with the recovery key
		// below we generate a root token using the recovery key
		rootStatus, err := client.Sys().GenerateRootStatus()
		if err != nil {
			t.Fatal(err)
		}
		otp, err := roottoken.GenerateOTP(rootStatus.OTPLength)
		if err != nil {
			t.Fatal(err)
		}
		genRoot, err := client.Sys().GenerateRootInit(otp, "")
		if err != nil {
			t.Fatal(err)
		}
		r, err := client.Sys().GenerateRootUpdate(recoveryKey, genRoot.Nonce)
		if err != nil {
			t.Fatal(err)
		}
		if !r.Complete {
			t.Fatal("expected root update to be complete")
		}
	})
	t.Run("backup", func(t *testing.T) {
		t.Parallel()

		pgpKey := "keybase:hashicorp"
		// pgpFingerprints := []string{"c874011f0ab405110d02105534365d9472d7468f"}

		client, keys, closer := testVaultServerUnseal(t)
		defer closer()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"-init",
			"-key-shares", "1",
			"-key-threshold", "1",
			"-pgp-keys", pgpKey,
			"-backup",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		// Get the status for the nonce
		status, err := client.Sys().RekeyStatus()
		if err != nil {
			t.Fatal(err)
		}
		nonce := status.Nonce

		var combined string
		// Supply the unseal keys
		for _, key := range keys {
			ui, cmd := testOperatorRekeyCommand(t)
			cmd.client = client

			code := cmd.Run([]string{
				"-nonce", nonce,
				key,
			})
			if exp := 0; code != exp {
				t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
			}

			// Append to our output string
			combined += ui.OutputWriter.String()
		}

		re := regexp.MustCompile(`Key 1 fingerprint: (.+); value: (.+)`)
		match := re.FindAllStringSubmatch(combined, -1)
		if len(match) < 1 || len(match[0]) < 3 {
			t.Fatalf("bad match: %#v", match)
		}

		// Grab the output fingerprint and encrypted key
		fingerprint, encryptedKey := match[0][1], match[0][2]

		// Get the backup
		ui, cmd = testOperatorRekeyCommand(t)
		cmd.client = client

		code = cmd.Run([]string{
			"-backup-retrieve",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		output := ui.OutputWriter.String()
		if !strings.Contains(output, fingerprint) {
			t.Errorf("expected %q to contain %q", output, fingerprint)
		}
		if !strings.Contains(output, encryptedKey) {
			t.Errorf("expected %q to contain %q", output, encryptedKey)
		}

		// Delete the backup
		ui, cmd = testOperatorRekeyCommand(t)
		cmd.client = client

		code = cmd.Run([]string{
			"-backup-delete",
		})
		if exp := 0; code != exp {
			t.Errorf("expected %d to be %d: %s", code, exp, ui.ErrorWriter.String())
		}

		secret, err := client.Sys().RekeyRetrieveBackup()
		if err == nil {
			t.Errorf("expected error: %#v", secret)
		}
	})

	t.Run("communication_failure", func(t *testing.T) {
		t.Parallel()

		client, closer := testVaultServerBad(t)
		defer closer()

		ui, cmd := testOperatorRekeyCommand(t)
		cmd.client = client

		code := cmd.Run([]string{
			"secret/foo",
		})
		if exp := 2; code != exp {
			t.Errorf("expected %d to be %d", code, exp)
		}

		expected := "Error getting rekey status: "
		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()
		if !strings.Contains(combined, expected) {
			t.Errorf("expected %q to contain %q", combined, expected)
		}
	})

	t.Run("no_tabs", func(t *testing.T) {
		t.Parallel()

		_, cmd := testOperatorRekeyCommand(t)
		assertNoTabs(t, cmd)
	})
}
