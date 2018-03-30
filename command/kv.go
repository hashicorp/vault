package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*KVCommand)(nil)

type KVCommand struct {
	*BaseCommand
}

func (c *KVCommand) Synopsis() string {
	return "Interact with Vault's Key-Value storage"
}

func (c *KVCommand) Help() string {
	helpText := `
Usage: vault kv <subcommand> [options] [args]

  This command has subcommands for interacting with Vault's key-value
  store. Here are some simple examples, and more detailed examples are
  available in the subcommands or the documentation.

  Create or update the key named "foo" in the "secret" mount with the value
  "bar=baz":

      $ vault kv put secret/foo bar=baz

  Read this value back:

      $ vault kv get secret/foo

  Get metadata for the key:

      $ vault kv metadata get secret/foo
	  
  Get a specific version of the key:

      $ vault kv get -version=1 secret/foo

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *KVCommand) Run(args []string) int {
	return cli.RunResultHelp
}
