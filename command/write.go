// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*WriteCommand)(nil)
	_ cli.CommandAutocomplete = (*WriteCommand)(nil)
)

// MFAMethodInfo contains the information about an MFA method
type MFAMethodInfo struct {
	methodID    string
	methodType  string
	usePasscode bool
}

// WriteCommand is a Command that puts data into the Vault.
type WriteCommand struct {
	*BaseCommand

	flagForce bool

	testStdin io.Reader // for tests
}

func (c *WriteCommand) Synopsis() string {
	return "Write data, configuration, and secrets"
}

func (c *WriteCommand) Help() string {
	helpText := `
Usage: vault write [options] PATH [DATA K=V...]

  Writes data to Vault at the given path. The data can be credentials, secrets,
  configuration, or arbitrary data. The specific behavior of this command is
  determined at the thing mounted at the path.

  Data is specified as "key=value" pairs. If the value begins with an "@", then
  it is loaded from a file. If the value is "-", Vault will read the value from
  stdin.

  Store an arbitrary secret in the token's cubbyhole.

      $ vault write cubbyhole/git-credentials username="student01" password="p@$$w0rd"

  Create a new encryption key in the transit secrets engine:

      $ vault write -force transit/keys/my-key

  The -force / -f flag allows a write operation without any input data.

  Upload an AWS IAM policy from a file on disk:

      $ vault write aws/roles/ops policy=@policy.json

  Configure access to Consul by providing an access token:

      $ echo $MY_TOKEN | vault write consul/config/access token=-

  Create a token

      $ vault write auth/token/create policies="admin" policies="secops" ttl=8h num_uses=3

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret engines in use.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *WriteCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "force",
		Aliases:    []string{"f"},
		Target:     &c.flagForce,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Allow the operation to continue with no key=value pairs. This " +
			"allows writing to keys that do not need or expect data.",
	})

	return set
}

func (c *WriteCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here. Without a way to access help
	// information, we don't know what paths we could write to.
	return complete.PredictAnything
}

func (c *WriteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *WriteCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) == 1 && !c.flagForce:
		c.UI.Error("Must supply data or use -force")
		return 1
	}

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	path := sanitizePath(args[0])

	data, err := parseArgsData(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().Write(path, data)
	return handleWriteSecretOutput(c.BaseCommand, path, secret, err)
}

func handleWriteSecretOutput(c *BaseCommand, path string, secret *api.Secret, err error) int {
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		// and even then, don't output anything if a specific field was requested
		if c.flagField == "" && Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	// Currently, if there is only one MFA method configured, the login
	// request is validated interactively
	methodInfo := c.getInteractiveMFAMethodInfo(secret)
	if methodInfo != nil {
		secret, err = c.validateMFA(secret.Auth.MFARequirement.MFARequestID, *methodInfo)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
	} else if c.getMFAValidationRequired(secret) {
		c.UI.Warn(wrapAtLength("A login request was issued that is subject to "+
			"MFA validation. Please make sure to validate the login by sending another "+
			"request to sys/mfa/validate endpoint.") + "\n")
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
