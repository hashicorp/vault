package command

import (
	"fmt"
	"strings"
)

// TokenRevokeCommand is a Command that mounts a new mount.
type TokenRevokeCommand struct {
	Meta
}

func (c *TokenRevokeCommand) Run(args []string) int {
	var mode string
	flags := c.Meta.FlagSet("token-revoke", FlagSetDefault)
	flags.StringVar(&mode, "mode", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf(
			"\ntoken-revoke expects one argument"))
		return 1
	}

	token := args[0]

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	var fn func(string) error
	switch mode {
	case "":
		fn = client.Auth().Token().RevokeTree
	case "orphan":
		fn = client.Auth().Token().RevokeOrphan
	case "path":
		fn = client.Auth().Token().RevokePrefix
	default:
		c.Ui.Error(fmt.Sprintf(
			"Unknown revocation mode: %s", mode))
		return 1
	}

	if err := fn(token); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error revoking token: %s", err))
		return 2
	}

	c.Ui.Output("Revocation successful.")
	return 0
}

func (c *TokenRevokeCommand) Synopsis() string {
	return "Revoke one or more auth tokens"
}

func (c *TokenRevokeCommand) Help() string {
	helpText := `
Usage: vault token-revoke [options] token

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
      the "token" arg above is actually a "path".

General Options:

  ` + generalOptionsUsage() + `

Token Options:

  -mode=value             The type of revocation to do. See the documentation
                          above for more information.

`
	return strings.TrimSpace(helpText)
}
