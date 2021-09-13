package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVPatchCommand)(nil)
	_ cli.CommandAutocomplete = (*KVPatchCommand)(nil)
)

type KVPatchCommand struct {
	*BaseCommand

	flagMethod string
	flagCAS int
	testStdin io.Reader // for tests
}

func (c *KVPatchCommand) Synopsis() string {
	return "Sets or updates data in the KV store without overwriting"
}

func (c *KVPatchCommand) Help() string {
	helpText := `
Usage: vault kv patch [options] KEY [DATA]

  *NOTE*: This is only supported for KV v2 engine mounts.

  Writes the data to the given path in the key-value store. The data can be of
  any type.

      $ vault kv patch secret/foo bar=baz

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ vault kv patch secret/foo @data.json

  Or it can be read from stdin using the "-" symbol:

      $ echo "abcd1234" | vault kv patch secret/foo bar=-

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVPatchCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringVar(&StringVar{
		Name:    "method",
		Aliases: []string{"m"},
		Target:  &c.flagMethod,
		Usage: `Specifies the patch method to use. The two options are "patch" and "rw" (read-then-write).
        If not provided, an HTTP PATCH will be attempted. If the request fails, a fall back to a read
        then write will take place.`,
	})

	f.IntVar(&IntVar{
		Name:    "cas",
		Target:  &c.flagCAS,
		Default: -1,
		Usage: `Specifies to use a Check-And-Set operation to be used for HTTP PATCH. Note in a
        read-then-write approach, the CAS value will be taken from the "version" of the read secret.
        If not set the write will be allowed. If set to 0 a write will only be allowed if the key
		doesn’t exist. If the index is non-zero the write will only be allowed if the key’s current
        version matches the version specified in the cas parameter.`,
	})

	return set
}

func (c *KVPatchCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *KVPatchCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVPatchCommand) Run(args []string) int {
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

	if c.flagMethod == "rw" && c.flagCAS > -1 {
		c.UI.Warn("Attempting to use cas flag with rw method but will use version from read secret")
	}

	var err error
	path := sanitizePath(args[0])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	newData, err := parseArgsData(stdin, args[1:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	mountPath, v2, err := isKVv2(path, client)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if !v2 {
		c.UI.Error(fmt.Sprintf("K/V engine mount must be version 2 for patch support"))
		return 2
	}

	path = addPrefixToVKVPath(path, mountPath, "data")

	if c.flagMethod == "" || c.flagMethod == "patch" {

		data := map[string]interface{}{
			"data":    newData,
			"options": map[string]interface{}{},
		}

		if c.flagCAS > -1 {
			data["options"].(map[string]interface{})["cas"] = c.flagCAS
		}

		secret, err := client.Logical().JSONMergePatch(path, data)

		if err != nil {
			c.UI.Warn(fmt.Sprintf("Unable to perform HTTP PATCH at %s\n\n%s\nFalling back to read then write\n", path, err))
		} else {
			return OutputSecret(c.UI, secret)
		}
	}

	if c.flagMethod == "" || c.flagMethod == "rw" {
		// First, do a read.
		// Note that we don't want to see curl output for the read request.
		curOutputCurl := client.OutputCurlString()
		client.SetOutputCurlString(false)
		secret, err := kvReadRequest(client, path, nil)
		client.SetOutputCurlString(curOutputCurl)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error doing pre-read at %s: %s", path, err))
			return 2
		}

		// Make sure a value already exists
		if secret == nil || secret.Data == nil {
			c.UI.Error(fmt.Sprintf("No value found at %s", path))
			return 2
		}

		// Verify metadata found
		rawMeta, ok := secret.Data["metadata"]
		if !ok || rawMeta == nil {
			c.UI.Error(fmt.Sprintf("No metadata found at %s; patch only works on existing data", path))
			return 2
		}
		meta, ok := rawMeta.(map[string]interface{})
		if !ok {
			c.UI.Error(fmt.Sprintf("Metadata found at %s is not the expected type (JSON object)", path))
			return 2
		}
		if meta == nil {
			c.UI.Error(fmt.Sprintf("No metadata found at %s; patch only works on existing data", path))
			return 2
		}

		// Verify old data found
		rawData, ok := secret.Data["data"]
		if !ok || rawData == nil {
			c.UI.Error(fmt.Sprintf("No data found at %s; patch only works on existing data", path))
			return 2
		}
		data, ok := rawData.(map[string]interface{})
		if !ok {
			c.UI.Error(fmt.Sprintf("Data found at %s is not the expected type (JSON object)", path))
			return 2
		}
		if data == nil {
			c.UI.Error(fmt.Sprintf("No data found at %s; patch only works on existing data", path))
			return 2
		}

		// Copy new data over
		for k, v := range newData {
			data[k] = v
		}

		secret, err = client.Logical().Write(path, map[string]interface{}{
			"data": data,
			"options": map[string]interface{}{
				"cas": meta["version"],
			},
		})
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

		if c.flagField != "" {
			return PrintRawField(c.UI, secret, c.flagField)
		}

		return OutputSecret(c.UI, secret)
	}

	c.UI.Error(fmt.Sprintf("Unknown method flag value %q", c.flagMethod))
	return 2
}
