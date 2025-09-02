// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVPutCommand)(nil)
	_ cli.CommandAutocomplete = (*KVPutCommand)(nil)
)

type KVPutCommand struct {
	*BaseCommand

	flagCAS   int
	flagMount string
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

      $ vault kv put -mount=secret foo bar=baz

  The deprecated path-like syntax can also be used, but this should be avoided 
  for KV v2, as the fact that it is not actually the full API path to 
  the secret (secret/data/foo) can cause confusion: 
  
      $ vault kv put secret/foo bar=baz

  The data can also be consumed from a file on disk by prefixing with the "@"
  symbol. For example:

      $ vault kv put -mount=secret foo @data.json

  Or it can be read from stdin using the "-" symbol:

      $ echo "abcd1234" | vault kv put -mount=secret foo bar=-

  To perform a Check-And-Set operation, specify the -cas flag with the
  appropriate version number corresponding to the key you want to perform
  the CAS operation on:

      $ vault kv put -mount=secret -cas=1 foo bar=baz

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

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "", // no default, because the handling of the next arg is determined by whether this flag has a value
		Usage: `Specifies the path where the KV backend is mounted. If specified, 
		the next argument will be interpreted as the secret path. If this flag is 
		not specified, the next argument will be interpreted as the combined mount 
		path and secret path, with /data/ automatically appended between KV 
		v2 secrets.`,
	})

	return set
}

func (c *KVPutCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
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

	// If true, we're working with "-mount=secret foo" syntax.
	// If false, we're using "secret/foo" syntax.
	mountFlagSyntax := c.flagMount != ""

	var (
		mountPath   string
		partialPath string
		v2          bool
	)

	// Parse the paths and grab the KV version
	if mountFlagSyntax {
		// In this case, this arg is the secret path (e.g. "foo").
		partialPath = sanitizePath(args[0])
		mountPath, v2, err = isKVv2(sanitizePath(c.flagMount), client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}

		if v2 {
			partialPath = path.Join(mountPath, partialPath)
		}
	} else {
		// In this case, this arg is a path-like combination of mountPath/secretPath.
		// (e.g. "secret/foo")
		partialPath = sanitizePath(args[0])
		mountPath, v2, err = isKVv2(partialPath, client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
	}

	// Add /data to v2 paths only
	var fullPath string
	if v2 {
		fullPath = addPrefixToKVPath(partialPath, mountPath, "data", false)
		data = map[string]interface{}{
			"data":    data,
			"options": map[string]interface{}{},
		}

		if c.flagCAS > -1 {
			data["options"].(map[string]interface{})["cas"] = c.flagCAS
		}
	} else {
		// v1
		if mountFlagSyntax {
			fullPath = path.Join(mountPath, partialPath)
		} else {
			fullPath = partialPath
		}
	}

	secret, err := client.Logical().Write(fullPath, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error writing data to %s: %s", fullPath, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		// Don't output anything unless using the "table" format
		if Format(c.UI) == "table" {
			c.UI.Info(fmt.Sprintf("Success! Data written to: %s", fullPath))
		}
		return 0
	}

	if c.flagField != "" {
		return PrintRawField(c.UI, secret, c.flagField)
	}

	// If the secret is wrapped, return the wrapped response.
	if secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, secret)
	}

	if Format(c.UI) == "table" {
		outputPath(c.UI, fullPath, "Secret Path")
		metadata := secret.Data
		c.UI.Info(getHeaderForMap("Metadata", metadata))
		return OutputData(c.UI, metadata)
	}

	return OutputSecret(c.UI, secret)
}
