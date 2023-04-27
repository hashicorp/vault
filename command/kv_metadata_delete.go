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
	_ cli.Command             = (*KVMetadataDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*KVMetadataDeleteCommand)(nil)
)

type KVMetadataDeleteCommand struct {
	*BaseCommand
	flagMount string
}

func (c *KVMetadataDeleteCommand) Synopsis() string {
	return "Deletes all versions and metadata for a key in the KV store"
}

func (c *KVMetadataDeleteCommand) Help() string {
	helpText := `
Usage: vault kv metadata delete [options] PATH

  Deletes all versions and metadata for the provided key. 

      $ vault kv metadata delete -mount=secret foo

  The deprecated path-like syntax can also be used, but this should be avoided 
  for KV v2, as the fact that it is not actually the full API path to 
  the secret (secret/metadata/foo) can cause confusion: 
  
      $ vault kv metadata delete secret/foo

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVMetadataDeleteCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "", // no default, because the handling of the next arg is determined by whether this flag has a value
		Usage: `Specifies the path where the KV backend is mounted. If specified, 
		the next argument will be interpreted as the secret path. If this flag is 
		not specified, the next argument will be interpreted as the combined mount 
		path and secret path, with /metadata/ automatically appended between KV 
		v2 secrets.`,
	})

	return set
}

func (c *KVMetadataDeleteCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFiles()
}

func (c *KVMetadataDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVMetadataDeleteCommand) Run(args []string) int {
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

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
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

	if !v2 {
		c.UI.Error("Metadata not supported on KV Version 1")
		return 1
	}

	fullPath := addPrefixToKVPath(partialPath, mountPath, "metadata")
	if secret, err := client.Logical().Delete(fullPath); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting %s: %s", fullPath, err))
		if secret != nil {
			OutputSecret(c.UI, secret)
		}
		return 2
	}

	c.UI.Info(fmt.Sprintf("Success! Data deleted (if it existed) at: %s", fullPath))
	return 0
}
