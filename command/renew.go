package command

import (
	"fmt"
	"strconv"
	"strings"
)

// RenewCommand is a Command that mounts a new mount.
type RenewCommand struct {
	Meta
}

func (c *RenewCommand) Run(args []string) int {
	var format string
	flags := c.Meta.FlagSet("renew", FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 || len(args) >= 3 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\nRenew expects at least one argument: the lease ID to renew"))
		return 1
	}

	var increment int
	leaseId := args[0]
	if len(args) > 1 {
		parsed, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Invalid increment, must be an int: %s", err))
			return 1
		}

		increment = int(parsed)
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err := client.Sys().Renew(leaseId, increment)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Renew error: %s", err))
		return 1
	}

	return OutputSecret(c.Ui, format, secret)
}

func (c *RenewCommand) Synopsis() string {
	return "Renew the lease of a secret"
}

func (c *RenewCommand) Help() string {
	helpText := `
Usage: vault renew [options] id [increment]

  Renew the lease on a secret, extending the time that it can be used
  before it is revoked by Vault.

  Every secret in Vault has a lease associated with it. If the user of
  the secret wants to use it longer than the lease, then it must be
  renewed. Renewing the lease will not change the contents of the secret.

  To renew a secret, run this command with the lease ID returned when it
  was read. Optionally, request a specific increment in seconds. Vault may
  is not required to honor this request.

General Options:

  -address=addr           The address of the Vault server.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-path is used.

  -tls-skip-verify        Do not verify TLS certificate. This is highly
                          not recommended. This is especially not recommended
                          for unsealing a vault.

Renew Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json.
`
	return strings.TrimSpace(helpText)
}
