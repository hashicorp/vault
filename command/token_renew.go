package command

import (
	"fmt"
	"strconv"
	"strings"
)

// TokenRenewCommand is a Command that mounts a new mount.
type TokenRenewCommand struct {
	Meta
}

func (c *TokenRenewCommand) Run(args []string) int {
	var format string
	flags := c.Meta.FlagSet("token-renew", FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-renew expects at least one argument"))
		return 1
	}

	var increment int
	token := args[0]
	if len(args) > 1 {
		value, err := strconv.ParseInt(args[1], 10, 0)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Invalid increment: %s", err))
			return 1
		}

		increment = int(value)
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err := client.Auth().Token().Renew(token, increment)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error renewing token: %s", err))
		return 1
	}

	return OutputSecret(c.Ui, format, secret)
}

func (c *TokenRenewCommand) Synopsis() string {
	return "Renew an auth token"
}

func (c *TokenRenewCommand) Help() string {
	helpText := `
Usage: vault token-renew [options] token [increment]

  Renew an auth token, extending the amount of time it can be used.

  This command is similar to "renew", but "renew" is only for lease IDs.
  This command is only for tokens.

  An optional increment can be given to request a certain number of
  seconds to increment the lease. This request is advisory; Vault may not
  adhere to it at all.

General Options:

  ` + generalOptionsUsage() + `

Token Renew Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json.

`
	return strings.TrimSpace(helpText)
}
