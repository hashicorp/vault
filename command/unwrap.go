package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

// UnwrapCommand is a Command that behaves like ReadCommand but specifically
// for unwrapping cubbyhole-wrapped secrets
type UnwrapCommand struct {
	meta.Meta
}

func (c *UnwrapCommand) Run(args []string) int {
	var format string
	var field string
	var err error
	var secret *api.Secret
	var flags *flag.FlagSet
	flags = c.Meta.FlagSet("unwrap", meta.FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&field, "field", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	var tokenID string

	args = flags.Args()
	switch len(args) {
	case 0:
	case 1:
		tokenID = args[0]
	default:
		c.Ui.Error("unwrap expects zero or one argument (the ID of the wrapping token)")
		flags.Usage()
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err = client.Logical().Unwrap(tokenID)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	if secret == nil {
		c.Ui.Error("Server gave empty response or secret returned was empty")
		return 1
	}

	// Handle single field output
	if field != "" {
		return PrintRawField(c.Ui, secret, field)
	}

	// Check if the original was a list response and format as a list if so
	if secret.Data != nil &&
		len(secret.Data) == 1 &&
		secret.Data["keys"] != nil {
		_, ok := secret.Data["keys"].([]interface{})
		if ok {
			return OutputList(c.Ui, format, secret)
		}
	}
	return OutputSecret(c.Ui, format, secret)
}

func (c *UnwrapCommand) Synopsis() string {
	return "Unwrap a wrapped secret"
}

func (c *UnwrapCommand) Help() string {
	helpText := `
Usage: vault unwrap [options] <wrapping token ID>

  Unwrap a wrapped secret.

  Unwraps the data wrapped by the given token ID. The returned result is the
  same as a 'read' operation on a non-wrapped secret.

General Options:
` + meta.GeneralOptionsUsage() + `
Read Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

  -field=field            If included, the raw value of the specified field
                          will be output raw to stdout.

`
	return strings.TrimSpace(helpText)
}
