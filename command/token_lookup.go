package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*TokenLookupCommand)(nil)
var _ cli.CommandAutocomplete = (*TokenLookupCommand)(nil)

type TokenLookupCommand struct {
	*BaseCommand

	flagAccessor bool
}

func (c *TokenLookupCommand) Synopsis() string {
	return "Display information about a token"
}

func (c *TokenLookupCommand) Help() string {
	helpText := `
Usage: vault token lookup [options] [TOKEN | ACCESSOR]

  Displays information about a token or accessor. If a TOKEN is not provided,
  the locally authenticated token is used.

  Get information about the locally authenticated token (this uses the
  /auth/token/lookup-self endpoint and permission):

      $ vault token lookup

  Get information about a particular token (this uses the /auth/token/lookup
  endpoint and permission):

      $ vault token lookup 96ddf4bc-d217-f3ba-f9bd-017055595017

  Get information about a token via its accessor:

      $ vault token lookup -accessor 9793c9b3-e04a-46f3-e7b8-748d7da248da

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenLookupCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "accessor",
		Target:     &c.flagAccessor,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Treat the argument as an accessor instead of a token. When " +
			"this option is selected, the output will NOT include the token.",
	})

	return set
}

func (c *TokenLookupCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *TokenLookupCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TokenLookupCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	token := ""

	args = f.Args()
	switch {
	case c.flagAccessor && len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments with -accessor (expected 1, got %d)", len(args)))
		return 1
	case c.flagAccessor && len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments with -accessor (expected 1, got %d)", len(args)))
		return 1
	case len(args) == 0:
		// Use the local token
	case len(args) == 1:
		token = strings.TrimSpace(args[0])
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0-1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var secret *api.Secret
	switch {
	case token == "":
		secret, err = client.Auth().Token().LookupSelf()
	case c.flagAccessor:
		secret, err = client.Auth().Token().LookupAccessor(token)
	default:
		secret, err = client.Auth().Token().Lookup(token)
	}

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error looking up token: %s", err))
		return 2
	}

	return OutputSecret(c.UI, secret)
}
