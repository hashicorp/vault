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
	_ cli.Command             = (*KVListCommand)(nil)
	_ cli.CommandAutocomplete = (*KVListCommand)(nil)
)

type KVListCommand struct {
	*BaseCommand
	flagMount string
}

func (c *KVListCommand) Synopsis() string {
	return "List data or secrets"
}

func (c *KVListCommand) Help() string {
	helpText := `

Usage: vault kv list [options] PATH

  Lists data from Vault's key-value store at the given path.

  List values under the "my-app" folder of the key-value store:

      $ vault kv list secret/my-app/

  Additional flags and more advanced use cases are detailed below.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *KVListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

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

func (c *KVListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *KVListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *KVListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		if c.flagMount == "" {
			c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
			return 1
		}
		args = []string{""}
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

	// Add /metadata to v2 paths only
	var fullPath string
	if v2 {
		fullPath = addPrefixToKVPath(partialPath, mountPath, "metadata")
	} else {
		// v1
		if mountFlagSyntax {
			fullPath = path.Join(mountPath, partialPath)
		} else {
			fullPath = partialPath
		}
	}

	secret, err := client.Logical().List(fullPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing %s: %s", fullPath, err))
		return 2
	}

	// If the secret is wrapped, return the wrapped response.
	if secret != nil && secret.WrapInfo != nil && secret.WrapInfo.TTL != 0 {
		return OutputSecret(c.UI, secret)
	}

	_, ok := extractListData(secret)
	if Format(c.UI) != "table" {
		if secret == nil || secret.Data == nil || !ok {
			OutputData(c.UI, map[string]interface{}{})
			return 2
		}
	}

	if secret == nil || secret.Data == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", fullPath))
		return 2
	}

	if !ok {
		c.UI.Error(fmt.Sprintf("No entries found at %s", fullPath))
		return 2
	}

	return OutputList(c.UI, secret)
}
