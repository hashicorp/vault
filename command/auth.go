package command

import (
	"flag"
	"io"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*AuthCommand)(nil)

type AuthCommand struct {
	*BaseCommand

	Handlers map[string]LoginHandler

	testStdin io.Reader // for tests
}

func (c *AuthCommand) Synopsis() string {
	return "Interact with auth methods"
}

func (c *AuthCommand) Help() string {
	return strings.TrimSpace(`
Usage: vault auth <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's auth methods.
  Users can list, enable, disable, and get help for different auth methods.

  To authenticate to Vault as a user or machine, use the "vault login" command
  instead. This command is for interacting with the auth methods themselves, not
  authenticating to Vault.

  List all enabled auth methods:

      $ vault auth list

  Enable a new auth method "userpass";

      $ vault auth enable userpass

  Get detailed help information about how to authenticate to a particular auth
  method:

      $ vault auth help github

  Please see the individual subcommand help for detailed usage information.
`)
}

func (c *AuthCommand) Run(args []string) int {
	// If we entered the run method, none of the subcommands picked up. This
	// means the user is still trying to use auth as "vault auth TOKEN" or
	// similar, so direct them to vault login instead.
	//
	// This run command is a bit messy to maintain BC for a bit. In the future,
	// it will just be a tiny function, but for now we have to maintain bc.
	//
	// Deprecation
	// TODO: remove in 0.9.0

	// Parse the args for our deprecations and defer to the proper areas.
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "-methods"):
			c.UI.Warn(wrapAtLength(
				"WARNING! The -methods flag is deprecated. Please use "+
					"\"vault auth list\" instead. This flag will be removed in "+
					"Vault 0.11 (or later).") + "\n")
			return (&AuthListCommand{
				BaseCommand: &BaseCommand{
					UI:     c.UI,
					client: c.client,
				},
			}).Run(nil)
		case strings.HasPrefix(arg, "-method-help"):
			c.UI.Warn(wrapAtLength(
				"WARNING! The -method-help flag is deprecated. Please use "+
					"\"vault auth help\" instead. This flag will be removed in "+
					"Vault 0.11 (or later).") + "\n")
			// Parse the args to pull out the method, surpressing any errors because
			// there could be other flags that we don't care about.
			f := flag.NewFlagSet("", flag.ContinueOnError)
			f.Usage = func() {}
			f.SetOutput(ioutil.Discard)
			flagMethod := f.String("method", "", "")
			f.Parse(args)

			return (&AuthHelpCommand{
				BaseCommand: &BaseCommand{
					UI:     c.UI,
					client: c.client,
				},
				Handlers: c.Handlers,
			}).Run([]string{*flagMethod})
		}
	}

	// If we got this far, we have an arg or a series of args that should be
	// passed directly to the new "vault login" command.
	c.UI.Warn(wrapAtLength(
		"WARNING! The \"vault auth ARG\" command is deprecated and is now a "+
			"subcommand for interacting with auth methods. To authenticate "+
			"locally to Vault, use \"vault login\" instead. This backwards "+
			"compatability will be removed in Vault 0.11 (or later).") + "\n")
	return (&LoginCommand{
		BaseCommand: &BaseCommand{
			UI:     c.UI,
			client: c.client,
		},
		Handlers: c.Handlers,
	}).Run(args)
}
