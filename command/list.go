package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

// ListCommand is a Command that lists data from the Vault.
type ListCommand struct {
	meta.Meta
}

func (c *ListCommand) Run(args []string) int {
	var format string
	var err error
	var secret *api.Secret
	var flags *flag.FlagSet
	flags = c.Meta.FlagSet("list", meta.FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 || len(args[0]) == 0 {
		c.Ui.Error("list expects one argument")
		flags.Usage()
		return 1
	}

	path := args[0]
	if path[0] == '/' {
		path = path[1:]
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err = client.Logical().List(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading %s: %s", path, err))
		return 1
	}
	if secret == nil {
		c.Ui.Error(fmt.Sprintf(
			"No value found at %s", path))
		return 1
	}
	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.Ui, format, secret)
	}

	if secret.Data["keys"] == nil {
		c.Ui.Error("No entries found")
		return 0
	}

	return OutputList(c.Ui, format, secret)
}

func (c *ListCommand) Synopsis() string {
	return "List data or secrets in Vault"
}

func (c *ListCommand) Help() string {
	helpText :=
		`
Usage: vault list [options] path

  List data from Vault.

  Retrieve a listing of available data. The data returned, if any, is backend-
  and endpoint-specific.

General Options:
` + meta.GeneralOptionsUsage() + `
Read Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.
`
	return strings.TrimSpace(helpText)
}
