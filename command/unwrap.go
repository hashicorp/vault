package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*UnwrapCommand)(nil)
var _ cli.CommandAutocomplete = (*UnwrapCommand)(nil)

// UnwrapCommand is a Command that behaves like ReadCommand but specifically for
// unwrapping cubbyhole-wrapped secrets
type UnwrapCommand struct {
	*BaseCommand
}

func (c *UnwrapCommand) Synopsis() string {
	return "Unwrap a wrapped secret"
}

func (c *UnwrapCommand) Help() string {
	helpText := `
Usage: vault unwrap [options] [TOKEN]

  Unwraps a wrapped secret from Vault by the given token. The result is the
  same as the "vault read" operation on the non-wrapped secret. If no token
  is given, the data in the currently authenticated token is unwrapped.

  Unwrap the data in the cubbyhole secrets engine for a token:

      $ vault unwrap 3de9ece1-b347-e143-29b0-dc2dc31caafd

  Unwrap the data in the active token:

      $ vault login 848f9ccf-7176-098c-5e2b-75a0689d41cd
      $ vault unwrap # unwraps 848f9ccf...

  For a full list of examples and paths, please see the online documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *UnwrapCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *UnwrapCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *UnwrapCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *UnwrapCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	token := ""
	switch len(args) {
	case 0:
		// Leave token as "", that will use the local token
	case 1:
		token = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0-1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().Unwrap(token)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error unwrapping: %s", err))
		return 2
	}
	if secret == nil {
		c.UI.Error("Could not find wrapped response")
		return 2
	}

	// Handle single field output
	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	// Check if the original was a list response and format as a list
	if _, ok := extractListData(secret); ok {
		return OutputList(c.UI, secret)
	}
	return OutputSecret(c.UI, secret)
}
