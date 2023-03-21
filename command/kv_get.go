// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"fmt"
	"path"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*KVGetCommand)(nil)
	_ cli.CommandAutocomplete = (*KVGetCommand)(nil)
)

type KVGetCommand struct {
	*BaseCommand

	flagVersion int
	flagMount   string
}

func (c *KVGetCommand) Synopsis() string {
	return "Retrieves data from the KV store"
}

func (c *KVGetCommand) Help() string {
	helpText := `
Usage: vault kv get [options] KEY

  Retrieves the value from Vault's key-value store at the given key name. If no
  key exists with that name, an error is returned. If a key exists with that
  name but has no data, nothing is returned.

      $ vault kv get -mount=secret foo

  The deprecated path-like syntax can also be used, but this should be avoided 
  for KV v2, as the fact that it is not actually the full API path to 
  the secret (secret/data/foo) can cause confusion: 
  
      $ vault kv get secret/foo

  To view the given key name at a specific version in time, specify the "-version"
  flag:

      $ vault kv get -mount=secret -version=1 foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *KVGetCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat | FlagSetClipboard)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.IntVar(&IntVar{
		Name:    "version",
		Target:  &c.flagVersion,
		Default: 0,
		Usage:   `If passed, the value at the version number will be returned.`,
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

func (c *KVGetCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *KVGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVGetCommand) Run(args []string) int {
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

	if !validateClipboardFlag(c.BaseCommand) {
		return 1
	}
	if !c.flagClipboard && c.flagClipboardTTL > 0 {
		c.UI.Error("-clipboard flag must be set to use -clipboard-ttl flag")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// If true, we're working with "-mount=secret foo" syntax.
	// If false, we're using "secret/foo" syntax.
	mountFlagSyntax := c.flagMount != ""

	var (
		mountPath string
		v2        bool
	)

	// Ignore leading slash
	partialPath := strings.TrimPrefix(args[0], "/")

	// Parse the paths and grab the KV version
	if mountFlagSyntax {
		// In this case, this arg is the secret path (e.g. "foo").
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
		mountPath, v2, err = isKVv2(partialPath, client)
		if err != nil {
			c.UI.Error(err.Error())
			return 2
		}
	}

	var versionParam map[string]string
	var fullPath string
	// Add /data to v2 paths only
	if v2 {
		fullPath = addPrefixToKVPath(partialPath, mountPath, "data")

		if c.flagVersion > 0 {
			versionParam = map[string]string{
				"version": fmt.Sprintf("%d", c.flagVersion),
			}
		}
	} else {
		// v1
		if mountFlagSyntax {
			fullPath = path.Join(mountPath, partialPath)
		} else {
			fullPath = partialPath
		}
	}

	secret, err := kvReadRequest(client, fullPath, versionParam)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading %s: %s", fullPath, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}
	if secret == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", fullPath))
		return 2
	}

	if c.flagField != "" {
		if v2 {
			// This is a v2, pass in the data field
			if data, ok := secret.Data["data"]; ok && data != nil {
				// If they requested a literal "data" see if they meant actual
				// value or the data block itself
				if c.flagField == "data" {
					if dataMap, ok := data.(map[string]interface{}); ok {
						if _, ok := dataMap["data"]; ok {
							return PrintRawField(c.UI, dataMap, c.flagField, c.flagClipboard, c.flagClipboardTTL)
						}
					}
					return PrintRawField(c.UI, secret, c.flagField, c.flagClipboard, c.flagClipboardTTL)
				}
				return PrintRawField(c.UI, data, c.flagField, c.flagClipboard, c.flagClipboardTTL)
			} else {
				c.UI.Error(fmt.Sprintf("No data found at %s", fullPath))
				return 2
			}
		} else {
			return PrintRawField(c.UI, secret, c.flagField, c.flagClipboard, c.flagClipboardTTL)
		}
	}

	// If we have wrap info print the secret normally.
	if secret.WrapInfo != nil || c.flagFormat != "table" {
		return OutputSecret(c.UI, secret)
	}

	if len(secret.Warnings) > 0 {
		tf := TableFormatter{}
		tf.printWarnings(c.UI, secret)
	}

	if v2 {
		outputPath(c.UI, fullPath, "Secret Path")
	}

	if metadata, ok := secret.Data["metadata"]; ok && metadata != nil {
		c.UI.Info(getHeaderForMap("Metadata", metadata.(map[string]interface{})))
		OutputData(c.UI, metadata)
		c.UI.Info("")
	}

	data := secret.Data
	if v2 && data != nil {
		data = nil
		dataRaw := secret.Data["data"]
		if dataRaw != nil {
			data = dataRaw.(map[string]interface{})
		}
	}

	if data != nil {
		c.UI.Info(getHeaderForMap("Data", data))
		OutputData(c.UI, data)
	}

	return 0
}
