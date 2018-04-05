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
  examples are available in the subcommands or the documentation.

  Create or update a metadata entry for a key: 

      $ vault kv metadata put -max-versions=5 secret/foo

  Get the metadata for a key, this provides information about each existing
  version:

      $ vault kv metadata get secret/foo

  Delete a key and all existing versions:

      $ vault kv metadata delete secret/foo

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *KVMetadataCommand) Run(args []string) int {
	return cli.RunResultHelp
}
