package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*TokenRevokeCommand)(nil)
var _ cli.CommandAutocomplete = (*TokenRevokeCommand)(nil)

type TokenRevokeCommand struct {
	*BaseCommand

	flagAccessor bool
	flagSelf     bool
	flagMode     string
}

func (c *TokenRevokeCommand) Synopsis() string {
	return "Revoke a token and its children"
}

func (c *TokenRevokeCommand) Help() string {
	helpText := `
Usage: vault token revoke [options] [TOKEN | ACCESSOR]

  Revokes authentication tokens and their children. If a TOKEN is not provided,
  the locally authenticated token is used. The "-mode" flag can be used to
  control the behavior of the revocation. See the "-mode" flag documentation
  for more information.

  Revoke a token and all the token's children:

      $ vault token revoke 96ddf4bc-d217-f3ba-f9bd-017055595017

  Revoke a token leaving the token's children:

      $ vault token revoke -mode=orphan 96ddf4bc-d217-f3ba-f9bd-017055595017

  Revoke a token by accessor:

      $ vault token revoke -accessor 9793c9b3-e04a-46f3-e7b8-748d7da248da

  For a full list of examples, please see the documentation.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenRevokeCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "accessor",
		Target:     &c.flagAccessor,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage:      "Treat the argument as an accessor instead of a token.",
	})

	f.BoolVar(&BoolVar{
		Name:       "self",
		Target:     &c.flagSelf,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage:      "Perform the revocation on the currently authenticated token.",
	})

	f.StringVar(&StringVar{
		Name:       "mode",
		Target:     &c.flagMode,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictSet("orphan", "path"),
		Usage: "Type of revocation to perform. If unspecified, Vault will revoke " +
			"the token and all of the token's children. If \"orphan\", Vault will " +
			"revoke only the token, leaving the children as orphans. If \"path\", " +
			"tokens created from the given authentication path prefix are deleted " +
			"along with their children.",
	})

	return set
}

func (c *TokenRevokeCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TokenRevokeCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TokenRevokeCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	token := ""
	if len(args) > 0 {
		token = strings.TrimSpace(args[0])
	}

	switch c.flagMode {
	case "", "orphan", "path":
	default:
		c.UI.Error(fmt.Sprintf("Invalid mode: %s", c.flagMode))
		return 1
	}

	switch {
	case c.flagSelf && len(args) > 0:
		c.UI.Error(fmt.Sprintf("Too many arguments with -self (expected 0, got %d)", len(args)))
		return 1
	case !c.flagSelf && len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1 or -self, got %d)", len(args)))
		return 1
	case !c.flagSelf && len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1 or -self, got %d)", len(args)))
		return 1
	case c.flagSelf && c.flagAccessor:
		c.UI.Error("Cannot use -self with -accessor!")
		return 1
	case c.flagSelf && c.flagMode != "":
		c.UI.Error("Cannot use -self with -mode!")
		return 1
	case c.flagAccessor && c.flagMode == "orphan":
		c.UI.Error("Cannot use -accessor with -mode=orphan!")
		return 1
	case c.flagAccessor && c.flagMode == "path":
		c.UI.Error("Cannot use -accessor with -mode=path!")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var revokeFn func(string) error
	// Handle all 6 possible combinations
	switch {
	case !c.flagAccessor && c.flagSelf && c.flagMode == "":
		revokeFn = client.Auth().Token().RevokeSelf
	case !c.flagAccessor && !c.flagSelf && c.flagMode == "":
		revokeFn = client.Auth().Token().RevokeTree
	case !c.flagAccessor && !c.flagSelf && c.flagMode == "orphan":
		revokeFn = client.Auth().Token().RevokeOrphan
	case !c.flagAccessor && !c.flagSelf && c.flagMode == "path":
		revokeFn = client.Sys().RevokePrefix
	case c.flagAccessor && !c.flagSelf && c.flagMode == "":
		revokeFn = client.Auth().Token().RevokeAccessor
	}

	if err := revokeFn(token); err != nil {
		c.UI.Error(fmt.Sprintf("Error revoking token: %s", err))
		return 2
	}

	c.UI.Output("Success! Revoked token (if it existed)")
	return 0
}
