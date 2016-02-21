package command

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
)

// TokenLookupCommand is a Command that outputs details about the
// provided.
type TokenLookupCommand struct {
	Meta
}

func (c *TokenLookupCommand) Run(args []string) int {
	var format string
	flags := c.Meta.FlagSet("token-lookup", FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) > 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-lookup expects at most one argument"))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	secret, err := doTokenLookup(args, client)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error looking up token: %s", err))
		return 1
	}
	return OutputSecret(c.Ui, format, secret)
}

func doTokenLookup(args []string, client *api.Client) (*api.Secret, error) {
	if len(args) == 0 {
		return client.Auth().Token().LookupSelf()
	}

	token := args[0]
	return client.Auth().Token().Lookup(token)
}

func (c *TokenLookupCommand) Synopsis() string {
	return "Display information about the specified token"
}

func (c *TokenLookupCommand) Help() string {
	helpText := `
Usage: vault token-lookup [options] [token]

  Displays information about the specified token.
  If no token is specified, the operation is performed on the currently
  authenticated token i.e. lookup-self.

General Options:

  ` + generalOptionsUsage() + `

Token Lookup Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

`
	return strings.TrimSpace(helpText)
}
