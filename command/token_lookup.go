package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

// TokenLookupCommand is a Command that outputs details about the
// provided.
type TokenLookupCommand struct {
	meta.Meta
}

func (c *TokenLookupCommand) Run(args []string) int {
	var format string
	var accessor bool
	flags := c.Meta.FlagSet("token-lookup", meta.FlagSetDefault)
	flags.BoolVar(&accessor, "accessor", false, "")
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
			"error initializing client: %s", err))
		return 2
	}

	var secret *api.Secret
	switch {
	case !accessor && len(args) == 0:
		secret, err = client.Auth().Token().LookupSelf()
	case !accessor && len(args) == 1:
		secret, err = client.Auth().Token().Lookup(args[0])
	case accessor && len(args) == 1:
		secret, err = client.Auth().Token().LookupAccessor(args[0])
	default:
		// This happens only when accessor is set and no argument is passed
		c.Ui.Error(fmt.Sprintf("token-lookup expects an argument when accessor flag is set"))
		return 1
	}

	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"error looking up token: %s", err))
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
Usage: vault token-lookup [options] [token|accessor]

  Displays information about the specified token. If no token is specified, the
  operation is performed on the currently authenticated token i.e. lookup-self.
  Information about the token can be retrieved using the token accessor via the 
  '-accessor' flag.

General Options:
` + meta.GeneralOptionsUsage() + `
Token Lookup Options:
  -accessor               A boolean flag, if set, treats the argument as an accessor of the token.
                          Note that the response of the command when this is set, will not contain
                          the token ID. Accessor is only meant for looking up the token properties
                          (and for revocation via '/auth/token/revoke-accessor/<accessor>' endpoint).

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

`
	return strings.TrimSpace(helpText)
}
