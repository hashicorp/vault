package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/ryanuber/columnize"
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

	return c.output(format, secret)
}

func (c *ReadCommand) output(format string, secret *api.Secret) int {
	switch format {
	case "json":
		return c.formatJSON(secret)
	case "table":
		fallthrough
	default:
		return c.formatTable(secret, true)
	}
}

func (c *ReadCommand) formatJSON(s *api.Secret) int {
	b, err := json.Marshal(s)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error formatting secret: %s", err))
		return 1
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	c.Ui.Output(out.String())
	return 0
}

func (c *ReadCommand) formatTable(s *api.Secret, whitespace bool) int {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)
	input = append(input, fmt.Sprintf("Key %s Value", config.Delim))

	if s.LeaseID != "" && s.LeaseDuration > 0 {
		input = append(input, fmt.Sprintf("lease_id %s %s", config.Delim, s.LeaseID))
		input = append(input, fmt.Sprintf(
			"lease_duration %s %d", config.Delim, s.LeaseDuration))
	}

	for k, v := range s.Data {
		input = append(input, fmt.Sprintf("%s %s %v", k, config.Delim, v))
	}

	c.Ui.Output(columnize.Format(input, config))
	return 0
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
