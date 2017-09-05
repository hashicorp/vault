package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*WriteCommand)(nil)
var _ cli.CommandAutocomplete = (*WriteCommand)(nil)

// WriteCommand is a Command that puts data into the Vault.
type WriteCommand struct {
	*BaseCommand

	flagForce bool

	testStdin io.Reader // for tests
}

func (c *WriteCommand) Synopsis() string {
	return "Writes data, configuration, and secrets"
}

func (c *WriteCommand) Help() string {
	helpText := `
Usage: vault write [options] PATH [DATA K=V...]

  Writes data to Vault at the given path. The data can be credentials, secrets,
  configuration, or arbitrary data. The specific behavior of this command is
  determined at the backend mounted at the path.

  Data is specified as "key=value" pairs. If the value begins with an "@", then
  it is loaded from a file. If the value is "-", Vault will read the value from
  stdin.

  Persist data in the static secret backend:

      $ vault write secret/my-secret foo=bar

  Create a new encryption key in the transit backend:

      $ vault write -f transit/keys/my-key

  Upload an AWS IAM policy from a file on disk:

      $ vault write aws/roles/ops policy=@policy.json

  Configure access to Consul by providing an access token:

      $ echo $MY_TOKEN | vault write consul/config/access token=-

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret backend in use.

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

	path, kvs, err := extractPath(f.Args())
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if len(kvs) == 0 && !c.flagForce {
		c.UI.Error("Missing DATA! Specify at least one K=V pair or use -force.")
		return 1
	}

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	data, err := parseArgsData(stdin, kvs)
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
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if c.flagFormat == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, c.flagFormat, secret)
}
