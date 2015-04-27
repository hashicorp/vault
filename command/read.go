package command

import (
	"fmt"
	"strings"
)

// ReadCommand is a Command that reads data from the Vault.
type ReadCommand struct {
	Meta
}

func (c *ReadCommand) Run(args []string) int {
	var format string
	flags := c.Meta.FlagSet("read", FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 || len(args) > 2 {
		c.Ui.Error("read expects one or two arguments")
		flags.Usage()
		return 1
	}
	path := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err := client.Logical().Read(path)
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

	return OutputSecret(c.Ui, format, secret)
}

func (c *ReadCommand) Synopsis() string {
	return "Read data or secrets from Vault"
}

func (c *ReadCommand) Help() string {
	helpText := `
Usage: vault read [options] path

  Read data from Vault.

  Read reads data at the given path from Vault. This can be used to
  read secrets and configuration as well as generate dynamic values from
  materialized backends. Please reference the documentation for the
  backends in use to determine key structure.

General Options:

  -address=TODO           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -insecure               Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

Read Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json.

`
	return strings.TrimSpace(helpText)
}
