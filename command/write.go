package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*WriteCommand)(nil)
	_ cli.CommandAutocomplete = (*WriteCommand)(nil)
)

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

  Persist data in the generic secrets engine:

      $ vault write secret/my-secret foo=bar

  Create a new encryption key in the transit secrets engine:

      $ vault write -f transit/keys/my-key

  Upload an AWS IAM policy from a file on disk:

      $ vault write aws/roles/ops policy=@policy.json

  Configure access to Consul by providing an access token:

      $ echo $MY_TOKEN | vault write consul/config/access token=-

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
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	if secret != nil && secret.Auth != nil && secret.Auth.MFARequirement != nil {
		if len(secret.Auth.MFARequirement.MFAConstraints) == 1 && !c.flagNonInteractive {
			// Currently, if there is only one MFA method configured, the login
			// request is validated interactively
			methodID, usePasscode := c.getMethodIDUsePasscode(secret.Auth.MFARequirement.MFAConstraints)
			if methodID != "" {
				return c.validateMFA(secret.Auth.MFARequirement.MFARequestID, methodID, usePasscode)
			}
		}
		c.UI.Warn(wrapAtLength("A login request was issued that is subject to "+
			"MFA validation. Please make sure to validate the login by sending another "+
			"request to mfa/validate endpoint.") + "\n")
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}

func (c *WriteCommand) getMethodIDUsePasscode(mfaConstraintAny map[string]*logical.MFAConstraintAny) (string, bool) {
	for _, mfaConstraint := range mfaConstraintAny {
		if len(mfaConstraint.Any) != 1 {
			return "", false
		}
		return mfaConstraint.Any[0].ID, mfaConstraint.Any[0].UsesPasscode
	}
	return "", false
}

func (c *WriteCommand) validateMFA(reqID, mfaMethodID string, usePasscode bool) int {
	var passcode string
	var err error
	if usePasscode {
		passcode, err = c.UI.AskSecret(fmt.Sprintf("Enter the passphrase for methodID %q:", mfaMethodID))
		if err != nil {
			c.UI.Error(fmt.Sprintf("failed to read the passphrase with error %q. please validate the login by sending a request to mfa/validate", err.Error()))
			return 2
		}
	} else {
		// TODO: not sure about printing this message; an error might occur before or during the validate request
		c.UI.Warn("Please acknowledge the push notification in your authenticator app")
	}

	// passcode could be an empty string
	mfaPayload := map[string][]string{
		mfaMethodID: {passcode},
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	path := "sys/mfa/validate"

	secret, err := client.Logical().Write(path, map[string]interface{}{
		"mfa_request_id": reqID,
		"mfa_payload":    mfaPayload,
	})
	if err != nil {
		c.UI.Error(err.Error())
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
