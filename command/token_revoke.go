package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/meta"
)

// TokenRevokeCommand is a Command that mounts a new mount.
type TokenRevokeCommand struct {
	meta.Meta
}

func (c *TokenRevokeCommand) Run(args []string) int {
	var mode string
	var accessor bool
	var self bool
	var token string
	flags := c.Meta.FlagSet("token-revoke", meta.FlagSetDefault)
	flags.BoolVar(&accessor, "accessor", false, "")
	flags.BoolVar(&self, "self", false, "")
	flags.StringVar(&mode, "mode", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	switch {
	case len(args) == 1 && !self:
		token = args[0]
	case len(args) != 0 && self:
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-revoke expects no arguments when revoking self"))
		return 1
	case len(args) != 1 && !self:
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-revoke expects one argument or the 'self' flag"))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	var fn func(string) error
	// Handle all 6 possible combinations
	switch {
	case !accessor && self && mode == "":
		fn = client.Auth().Token().RevokeSelf
	case !accessor && !self && mode == "":
		fn = client.Auth().Token().RevokeTree
	case !accessor && !self && mode == "orphan":
		fn = client.Auth().Token().RevokeOrphan
	case !accessor && !self && mode == "path":
		fn = client.Sys().RevokePrefix
	case accessor && !self && mode == "":
		fn = client.Auth().Token().RevokeAccessor
	case accessor && self:
		c.Ui.Error("token-revoke cannot be run on self when 'accessor' flag is set")
		return 1
	case self && mode != "":
		c.Ui.Error("token-revoke cannot be run on self when 'mode' flag is set")
		return 1
	case accessor && mode == "orphan":
		c.Ui.Error("token-revoke cannot be run for 'orphan' mode when 'accessor' flag is set")
		return 1
	case accessor && mode == "path":
		c.Ui.Error("token-revoke cannot be run for 'path' mode when 'accessor' flag is set")
		return 1
	}

	if err := fn(token); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error revoking token: %s", err))
		return 2
	}

	c.Ui.Output("Success! Token revoked if it existed.")
	return 0
}

func (c *TokenRevokeCommand) Synopsis() string {
	return "Revoke one or more auth tokens"
}

func (c *TokenRevokeCommand) Help() string {
	helpText := `
Usage: vault token-revoke [options] [token|accessor]

  Revoke one or more auth tokens.

  This command revokes auth tokens. Use the "revoke" command for
  revoking secrets.

  Depending on the flags used, auth tokens can be revoked in multiple ways
  depending on the "-mode" flag:

    * Without any value, the token specified and all of its children
      will be revoked.

    * With the "orphan" value, only the specific token will be revoked.
      All of its children will be orphaned.

    * With the "path" value, tokens created from the given auth path
      prefix will be deleted, along with all their children. In this case
      the "token" arg above is actually a "path". This mode does *not*
      work with token values or parts of token values.

  Token can be revoked using the token accessor. This can be done by
  setting the '-accessor' flag. Note that when '-accessor' flag is set,
  '-mode' should not be set for 'orphan' or 'path'. This is because,
  a token accessor always revokes the token along with its child tokens.

General Options:
` + meta.GeneralOptionsUsage() + `
Token Options:

  -accessor               A boolean flag, if set, treats the argument as an accessor of the token.
                          Note that accessor can also be used for looking up the token properties
                          via '/auth/token/lookup-accessor/<accessor>' endpoint.
                          Accessor is used when there is no access to token ID.

  -self                   A boolean flag, if set, the operation is performed on the currently
                          authenticated token i.e. lookup-self.

  -mode=value             The type of revocation to do. See the documentation
                          above for more information.

`
	return strings.TrimSpace(helpText)
}
