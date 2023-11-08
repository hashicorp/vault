// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command_server

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/command"

	"github.com/hashicorp/go-secure-stdlib/password"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/sdk/helper/roottoken"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorGenerateRootCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorGenerateRootCommand)(nil)
)

type generateRootKind int

const (
	generateRootRegular generateRootKind = iota
	generateRootDR
	generateRootRecovery
)

type OperatorGenerateRootCommand struct {
	*command.BaseCommand

	flagInit          bool
	flagCancel        bool
	flagStatus        bool
	flagDecode        string
	flagOTP           string
	flagPGPKey        string
	flagNonce         string
	flagGenerateOTP   bool
	flagDRToken       bool
	flagRecoveryToken bool

	testStdin io.Reader // for tests
}

func (c *OperatorGenerateRootCommand) Synopsis() string {
	return "Generates a new root, DR operation, or recovery token"
}

func (c *OperatorGenerateRootCommand) Help() string {
	helpText := `
Usage: vault operator generate-root [options] -init [-otp=...] [-pgp-key=...]
       vault operator generate-root [options] [-nonce=... KEY]
       vault operator generate-root [options] -decode=... -otp=...
       vault operator generate-root [options] -generate-otp
       vault operator generate-root [options] -status
       vault operator generate-root [options] -cancel

  Generates a new root token by combining a quorum of share holders.

  This command is unusual, as it is effectively six separate subcommands,
  selected via the options -init, -decode, -generate-otp, -status, -cancel, 
  or the absence of any of the previous five options (which selects the 
  "provide a key share" form).

  With the -dr-token or -recovery-token options, a DR operation token or a
  recovery token is generated instead of a root token - the relevant option
  must be included in every form of the generate-root command.

  Form 1 (-init) - Start a token generation:

    When starting a root or privileged operation token generation, you must
    choose one of the following protection methods for how the token will be
    returned:

      - A base64-encoded one-time-password (OTP). The resulting token is XORed
        with this value when it is returned. Use the "-decode" form of this
        command to output the final value.

        The Vault server will generate a suitable OTP for you, and return it:

            $ vault operator generate-root -init

        Vault versions before 0.11.2, released in 2018, required you to
        generate your own OTP (see the "-generate-otp" form) and pass it in,
        but this is no longer necessary. The command is still supported for
        compatibility, though:

            $ vault operator generate-root -init -otp="..."

      - A PGP key. The resulting token is encrypted with this public key.
        The key may be specified as a path to a file, or a string of the
        form "keybase:<username>" to fetch the key from the keybase.io API.

            $ vault operator generate-root -init -pgp-key="..."

  Form 2 (no option) - Enter an unseal key to progress root token generation:

    In the sub-form intended for interactive use, the command will
    automatically look up the nonce of the currently active generation
    operation, and will prompt for the key to be entered:

        $ vault operator generate-root

    In the sub-form intended for automation, the operation nonce must be
    explicitly provided, and the key is provided directly on the command line

        $ vault operator generate-root -nonce=... KEY

    If key is specified as "-", the command will read from stdin.

  Form 3 (-decode) - Decode a generated token protected with an OTP:

        $ vault operator generate-root -decode=ENCODED_TOKEN -otp=OTP

    If encoded token is specified as "-", the command will read from stdin.

  Form 4 (-generate-otp) - Generate an OTP code for the final token:

        $ vault operator generate-root -generate-otp

    Since changes in Vault 0.11.2 in 2018, there is no longer any reason to
    use this form, as a suitable OTP will be returned as part of the "-init"
    command.

  Form 5 (-status) - Get the status of a token generation that is in progress:

        $ vault operator generate-root -status

    This form also returns the length of the a correct OTP, for the running
    version and configuration of Vault.

  Form 6 (-cancel) - Cancel a token generation that is in progress:

    This would be used to remove an in progress generation operation, so that
    a new one can be started with different parameters.

        $ vault operator generate-root -cancel

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *OperatorGenerateRootCommand) Flags() *command.FlagSets {
	set := c.FlagSet(command.FlagSetHTTP | command.FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&command.BoolVar{
		Name:       "init",
		Target:     &c.flagInit,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Start a root token generation. This can only be done if " +
			"there is not currently one in progress.",
	})

	f.BoolVar(&command.BoolVar{
		Name:       "cancel",
		Target:     &c.flagCancel,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Reset the root token generation progress. This will discard any " +
			"submitted unseal keys or configuration.",
	})

	f.BoolVar(&command.BoolVar{
		Name:       "status",
		Target:     &c.flagStatus,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Print the status of the current attempt without providing an " +
			"unseal key.",
	})

	f.StringVar(&command.StringVar{
		Name:       "decode",
		Target:     &c.flagDecode,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "The value to decode; setting this triggers a decode operation. " +
			" If the value is \"-\" then read the encoded token from stdin.",
	})

	f.BoolVar(&command.BoolVar{
		Name:       "generate-otp",
		Target:     &c.flagGenerateOTP,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Generate and print a high-entropy one-time-password (OTP) " +
			"suitable for use with the \"-init\" flag.",
	})

	f.BoolVar(&command.BoolVar{
		Name:       "dr-token",
		Target:     &c.flagDRToken,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Set this flag to do generate root operations on DR operation " +
			"tokens.",
	})

	f.BoolVar(&command.BoolVar{
		Name:       "recovery-token",
		Target:     &c.flagRecoveryToken,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Set this flag to do generate root operations on recovery " +
			"tokens.",
	})

	f.StringVar(&command.StringVar{
		Name:       "otp",
		Target:     &c.flagOTP,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage:      "OTP code to use with \"-decode\" or \"-init\".",
	})

	f.VarFlag(&command.VarFlag{
		Name:       "pgp-key",
		Value:      (*pgpkeys.PubKeyFileFlag)(&c.flagPGPKey),
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Path to a file on disk containing a binary or base64-encoded " +
			"public PGP key. This can also be specified as a Keybase username " +
			"using the format \"keybase:<username>\". When supplied, the generated " +
			"root token will be encrypted and base64-encoded with the given public " +
			"key. Must be used with \"-init\".",
	})

	f.StringVar(&command.StringVar{
		Name:       "nonce",
		Target:     &c.flagNonce,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Nonce value returned at initialization. The same nonce value " +
			"must be provided with each unseal or recovery key. Only needed " +
			"when providing an unseal or recovery key.",
	})

	return set
}

func (c *OperatorGenerateRootCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *OperatorGenerateRootCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorGenerateRootCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0-1, got %d)", len(args)))
		return 1
	}

	if c.flagDRToken && c.flagRecoveryToken {
		c.UI.Error("Both -recovery-token and -dr-token flags are set")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	kind := generateRootRegular
	switch {
	case c.flagDRToken:
		kind = generateRootDR
	case c.flagRecoveryToken:
		kind = generateRootRecovery
	}

	switch {
	case c.flagGenerateOTP:
		otp, code := c.generateOTP(client, kind)
		if code == 0 {
			switch command.Format(c.UI) {
			case "", "table":
				return command.PrintRaw(c.UI, otp)
			default:
				status := map[string]interface{}{
					"otp":        otp,
					"otp_length": len(otp),
				}
				return command.OutputData(c.UI, status)
			}
		}
		return code
	case c.flagDecode != "":
		return c.decode(client, c.flagDecode, c.flagOTP, kind)
	case c.flagCancel:
		return c.cancel(client, kind)
	case c.flagInit:
		return c.init(client, c.flagOTP, c.flagPGPKey, kind)
	case c.flagStatus:
		return c.status(client, kind)
	default:
		// If there are no other flags, prompt for an unseal key.
		key := ""
		if len(args) > 0 {
			key = strings.TrimSpace(args[0])
		}
		return c.provide(client, key, kind)
	}
}

// generateOTP generates a suitable OTP code for generating a root token.
func (c *OperatorGenerateRootCommand) generateOTP(client *api.Client, kind generateRootKind) (string, int) {
	f := client.Sys().GenerateRootStatus
	switch kind {
	case generateRootDR:
		f = client.Sys().GenerateDROperationTokenStatus
	case generateRootRecovery:
		f = client.Sys().GenerateRecoveryOperationTokenStatus
	}

	status, err := f()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting root generation status: %s", err))
		return "", 2
	}

	otp, err := roottoken.GenerateOTP(status.OTPLength)
	var retCode int
	if err != nil {
		retCode = 2
		c.UI.Error(err.Error())
	} else {
		retCode = 0
	}
	return otp, retCode
}

// decode decodes the given value using the otp.
func (c *OperatorGenerateRootCommand) decode(client *api.Client, encoded, otp string, kind generateRootKind) int {
	if encoded == "" {
		c.UI.Error("Missing encoded value: use -decode=<string> to supply it")
		return 1
	}
	if otp == "" {
		c.UI.Error("Missing otp: use -otp to supply it")
		return 1
	}

	if encoded == "-" {
		// Pull our fake stdin if needed
		stdin := (io.Reader)(os.Stdin)
		if c.testStdin != nil {
			stdin = c.testStdin
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdin); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to read from stdin: %s", err))
			return 1
		}

		encoded = buf.String()

		if encoded == "" {
			c.UI.Error("Missing encoded value. When using -decode=\"-\" value must be passed via stdin.")
			return 1
		}
	}

	f := client.Sys().GenerateRootStatus
	switch kind {
	case generateRootDR:
		f = client.Sys().GenerateDROperationTokenStatus
	case generateRootRecovery:
		f = client.Sys().GenerateRecoveryOperationTokenStatus
	}

	status, err := f()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting root generation status: %s", err))
		return 2
	}

	token, err := roottoken.DecodeToken(encoded, otp, status.OTPLength)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error decoding root token: %s", err))
		return 1
	}

	switch command.Format(c.UI) {
	case "", "table":
		return command.PrintRaw(c.UI, token)
	default:
		tokenJSON := map[string]interface{}{
			"token": token,
		}
		return command.OutputData(c.UI, tokenJSON)
	}
}

// init is used to start the generation process
func (c *OperatorGenerateRootCommand) init(client *api.Client, otp, pgpKey string, kind generateRootKind) int {
	// Validate incoming fields. Either OTP OR PGP keys must be supplied.
	if otp != "" && pgpKey != "" {
		c.UI.Error("Error initializing: cannot specify both -otp and -pgp-key")
		return 1
	}

	// Start the root generation
	f := client.Sys().GenerateRootInit
	switch kind {
	case generateRootDR:
		f = client.Sys().GenerateDROperationTokenInit
	case generateRootRecovery:
		f = client.Sys().GenerateRecoveryOperationTokenInit
	}
	status, err := f(otp, pgpKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing root generation: %s", err))
		return 2
	}

	switch command.Format(c.UI) {
	case "table":
		return c.printStatus(status)
	default:
		return command.OutputData(c.UI, status)
	}
}

// provide prompts the user for the seal key and posts it to the update root
// endpoint. If this is the last unseal, this function outputs it.
func (c *OperatorGenerateRootCommand) provide(client *api.Client, key string, kind generateRootKind) int {
	f := client.Sys().GenerateRootStatus
	switch kind {
	case generateRootDR:
		f = client.Sys().GenerateDROperationTokenStatus
	case generateRootRecovery:
		f = client.Sys().GenerateRecoveryOperationTokenStatus
	}
	status, err := f()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting root generation status: %s", err))
		return 2
	}

	// Verify a root token generation is in progress. If there is not one in
	// progress, return an error instructing the user to start one.
	if !status.Started {
		c.UI.Error(command.WrapAtLength(
			"No root generation is in progress. Start a root generation by " +
				"running \"vault operator generate-root -init\"."))
		c.UI.Warn(command.WrapAtLength(fmt.Sprintf(
			"If starting root generation using the OTP method and generating "+
				"your own OTP, the length of the OTP string needs to be %d "+
				"characters in length.", status.OTPLength)))
		return 1
	}

	var nonce string

	switch key {
	case "-": // Read from stdin
		nonce = c.flagNonce

		// Pull our fake stdin if needed
		stdin := (io.Reader)(os.Stdin)
		if c.testStdin != nil {
			stdin = c.testStdin
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdin); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to read from stdin: %s", err))
			return 1
		}

		key = buf.String()
	case "": // Prompt using the tty
		// Nonce value is not required if we are prompting via the terminal
		nonce = status.Nonce

		w := command.GetWriterFromUI(c.UI)
		fmt.Fprintf(w, "Operation nonce: %s\n", nonce)
		fmt.Fprintf(w, "Unseal Key (will be hidden): ")
		key, err = password.Read(os.Stdin)
		fmt.Fprintf(w, "\n")
		if err != nil {
			if err == password.ErrInterrupted {
				c.UI.Error("user canceled")
				return 1
			}

			c.UI.Error(command.WrapAtLength(fmt.Sprintf("An error occurred attempting to "+
				"ask for the unseal key. The raw error message is shown below, but "+
				"usually this is because you attempted to pipe a value into the "+
				"command or you are executing outside of a terminal (tty). If you "+
				"want to pipe the value, pass \"-\" as the argument to read from "+
				"stdin. The raw error was: %s", err)))
			return 1
		}
	default: // Supplied directly as an arg
		nonce = c.flagNonce
	}

	// Trim any whitespace from they key, especially since we might have prompted
	// the user for it.
	key = strings.TrimSpace(key)

	// Verify we have a nonce value
	if nonce == "" {
		c.UI.Error("Missing nonce value: specify it via the -nonce flag")
		return 1
	}

	// Provide the key, this may potentially complete the update
	fUpd := client.Sys().GenerateRootUpdate
	switch kind {
	case generateRootDR:
		fUpd = client.Sys().GenerateDROperationTokenUpdate
	case generateRootRecovery:
		fUpd = client.Sys().GenerateRecoveryOperationTokenUpdate
	}
	status, err = fUpd(key, nonce)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error posting unseal key: %s", err))
		return 2
	}
	switch command.Format(c.UI) {
	case "table":
		return c.printStatus(status)
	default:
		return command.OutputData(c.UI, status)
	}
}

// cancel cancels the root token generation
func (c *OperatorGenerateRootCommand) cancel(client *api.Client, kind generateRootKind) int {
	f := client.Sys().GenerateRootCancel
	switch kind {
	case generateRootDR:
		f = client.Sys().GenerateDROperationTokenCancel
	case generateRootRecovery:
		f = client.Sys().GenerateRecoveryOperationTokenCancel
	}
	if err := f(); err != nil {
		c.UI.Error(fmt.Sprintf("Error canceling root token generation: %s", err))
		return 2
	}
	c.UI.Output("Success! Root token generation canceled (if it was started)")
	return 0
}

// status is used just to fetch and dump the status
func (c *OperatorGenerateRootCommand) status(client *api.Client, kind generateRootKind) int {
	f := client.Sys().GenerateRootStatus
	switch kind {
	case generateRootDR:
		f = client.Sys().GenerateDROperationTokenStatus
	case generateRootRecovery:
		f = client.Sys().GenerateRecoveryOperationTokenStatus
	}

	status, err := f()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting root generation status: %s", err))
		return 2
	}
	switch command.Format(c.UI) {
	case "table":
		return c.printStatus(status)
	default:
		return command.OutputData(c.UI, status)
	}
}

// printStatus dumps the status to output
func (c *OperatorGenerateRootCommand) printStatus(status *api.GenerateRootStatusResponse) int {
	out := []string{}
	out = append(out, fmt.Sprintf("Nonce | %s", status.Nonce))
	out = append(out, fmt.Sprintf("Started | %t", status.Started))
	out = append(out, fmt.Sprintf("Progress | %d/%d", status.Progress, status.Required))
	out = append(out, fmt.Sprintf("Complete | %t", status.Complete))
	if status.PGPFingerprint != "" {
		out = append(out, fmt.Sprintf("PGP Fingerprint | %s", status.PGPFingerprint))
	}
	switch {
	case status.EncodedToken != "":
		out = append(out, fmt.Sprintf("Encoded Token | %s", status.EncodedToken))
	case status.EncodedRootToken != "":
		out = append(out, fmt.Sprintf("Encoded Root Token | %s", status.EncodedRootToken))
	}
	if status.OTP != "" {
		c.UI.Warn(command.WrapAtLength("A One-Time-Password has been generated for you and is shown in the OTP field. You will need this value to decode the resulting root token, so keep it safe."))
		out = append(out, fmt.Sprintf("OTP | %s", status.OTP))
	}
	if status.OTPLength != 0 {
		out = append(out, fmt.Sprintf("OTP Length | %d", status.OTPLength))
	}

	output := command.ColumnOutput(out, nil)
	c.UI.Output(output)
	return 0
}
