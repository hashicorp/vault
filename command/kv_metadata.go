package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*KVMetadataCommand)(nil)

type KVMetadataCommand struct {
	*BaseCommand
}

func (c *KVMetadataCommand) Synopsis() string {
	return "Interact with Vault's Key-Value storage"
}

func (c *KVMetadataCommand) Help() string {
	helpText := `
Usage: vault kv metadata <subcommand> [options] [args]

  This command has subcommands for interacting with the metadata endpoint in
  Vault's key-value store. Here are some simple examples, and more detailed
  examples are  available in the subcommands or the documentation.

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

func (c *KVMetadataCommand) Run(args []string) int {
	return cli.RunResultHelp
}
