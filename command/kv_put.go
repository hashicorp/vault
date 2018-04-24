package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*KVPutCommand)(nil)
var _ cli.CommandAutocomplete = (*KVPutCommand)(nil)

type KVPutCommand struct {
	*BaseCommand

	flagCAS   int
	testStdin io.Reader // for tests
}

func (c *KVPutCommand) Synopsis() string {
	return "Sets or updates data in the KV store"
}

func (c *KVPutCommand) Help() string {
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
  appropriate version number corresponding to the key you want to perform
  the CAS operation on:

      $ vault kv put -cas=1 secret/foo bar=baz

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVPutCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.IntVar(&IntVar{
		Name:    "cas",
		Target:  &c.flagCAS,
		Default: -1,
		Usage: `Specifies to use a Check-And-Set operation. If not set the write
		will be allowed. If set to 0 a write will only be allowed if the key
		doesn’t exist. If the index is non-zero the write will only be allowed
		if the key’s current version matches the version specified in the cas
		parameter.`,
	})

	return set
}

func (c *KVPutCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVPutCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVPutCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)
	if c.testStdin != nil {
		stdin = c.testStdin
	}

	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected >1, got %d)", len(args)))
		return 1
	case len(args) == 1:
		c.UI.Error("Must supply data")
		return 1
	}

	var err error
	path := sanitizePath(args[0])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data, err := parseArgsData(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	mountPath, v2, err := isKVv2(path, client)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if v2 {
		path = addPrefixToVKVPath(path, mountPath, "data")
		data = map[string]interface{}{
			"data":    data,
			"options": map[string]interface{}{},
		}

		if c.flagCAS > -1 {
			data["options"].(map[string]interface{})["cas"] = c.flagCAS
		}
	}

	secret, err := client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", path, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", path))
		}
		return 0
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	return OutputSecret(c.UI, secret)
}
