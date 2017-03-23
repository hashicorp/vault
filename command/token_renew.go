package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/meta"
)

// TokenRenewCommand is a Command that mounts a new mount.
type TokenRenewCommand struct {
	meta.Meta
}

func (c *TokenRenewCommand) Run(args []string) int {
	var format, increment string
	flags := c.Meta.FlagSet("token-renew", meta.FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&increment, "increment", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 2 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-renew expects at most two arguments"))
		return 1
	}

	var token string
	if len(args) > 0 {
		token = args[0]
	}

	var inc int
	// If both are specified prefer the argument
	if len(args) == 2 {
		increment = args[1]
	}
	if increment != "" {
		dur, err := parseutil.ParseDurationSecond(increment)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Invalid increment: %s", err))
			return 1
		}

		inc = int(dur / time.Second)
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	// If the given token is the same as the client's, use renew-self instead
	// as this is far more likely to be allowed via policy
	var secret *api.Secret
	if token == "" {
		secret, err = client.Auth().Token().RenewSelf(inc)
	} else {
		secret, err = client.Auth().Token().Renew(token, inc)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error renewing token: %s", err))
		return 1
	}

	return OutputSecret(c.Ui, format, secret)
}

func (c *TokenRenewCommand) Synopsis() string {
	return "Renew an auth token if there is an associated lease"
}

func (c *TokenRenewCommand) Help() string {
	helpText := `
Usage: vault token-renew [options] [token] [increment]

  Renew an auth token, extending the amount of time it can be used. If a token
  is given to the command, '/auth/token/renew' will be called with the given
  token; otherwise, '/auth/token/renew-self' will be called with the client
  token.

  This command is similar to "renew", but "renew" is only for leases; this
  command is only for tokens.

  An optional increment can be given to request a certain number of seconds to
  increment the lease. This request is advisory; Vault may not adhere to it at
  all. If a token is being passed in on the command line, the increment can as
  well; otherwise it must be passed in via the '-increment' flag.

General Options:
` + meta.GeneralOptionsUsage() + `
Token Renew Options:

  -increment=3600         The desired increment. If not supplied, Vault will
                          use the default TTL. If supplied, it may still be
                          ignored. This can be submitted as an integer number
                          of seconds or a string duration (e.g. "72h").

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

`
	return strings.TrimSpace(helpText)
}
