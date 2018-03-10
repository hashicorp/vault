package command

import (
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVMetadataPutCommand)(nil)
var _ cli.CommandAutocomplete = (*KVMetadataPutCommand)(nil)

type KVMetadataPutCommand struct {
	*BaseCommand

	flagMaxVersions int
	flagCASRequired bool
	testStdin       io.Reader // for tests
}

func (c *KVMetadataPutCommand) Synopsis() string {
	return "Sets or updates data in the KV store"
}

func (c *KVMetadataPutCommand) Help() string {
	helpText := `
Usage: vault kv put [options] KEY [DATA]

  Writes the data to the given path in the key-value store. The data can be of
  any type.

      $ vault kv put secret/foo bar=baz

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ vault kv put secret/foo @data.json

  Or it can be read from stdin using the "-" symbol:

      $ echo "abcd1234" | vault kv put secret/foo bar=-

  To perform a Check-And-Set operation, specify the -cas flag with the
  appropriate version numer corresponding to the key you want to perform
  the CAS operation on:

      $ vault kv put -cas=1 secret/foo bar=baz

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVMetadataPutCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.IntVar(&IntVar{
		Name:    "max-versions",
		Target:  &c.flagMaxVersions,
		Default: 0,
		Usage:   `The number of versions to keep. If not set, the backend’s configured max version is used.`,
	})

	f.BoolVar(&BoolVar{
		Name:    "cas-required",
		Target:  &c.flagCASRequired,
		Default: false,
		Usage:   `If true the key will require the cas parameter to be set on all write requests. If false, the backend’s configuration will be used.`,
	})

	return set
}

func (c *KVMetadataPutCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVMetadataPutCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataPutCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	var err error
	path := sanitizePath(args[0])
	path, err = addPrefixToVKVPath(path, "metadata")
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data := map[string]interface{}{
		"max_versions": c.flagMaxVersions,
		"cas_required": c.flagCASRequired,
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	return OutputSecret(c.UI, secret)
}
