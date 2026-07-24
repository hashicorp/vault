// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*ListCommand)(nil)
	_ cli.CommandAutocomplete = (*ListCommand)(nil)
)

type ListCommand struct {
	*BaseCommand
	flagRecursive   bool
	flagFuzzy       bool
	flagPattern     string
	flagPermissions string
	flagAllMounts   bool
	flagRecursiveMode string // "cli" for CLI recursive, "api" for API ListRecursive
}

func (c *ListCommand) Synopsis() string {
	return "List data or secrets"
}

func (c *ListCommand) Help() string {
	helpText := `

Usage: vault list [options] PATH

  Lists data from Vault at the given path. This can be used to list keys in a,
  given secret engine.

  List values under the "my-app" folder of the generic secret engine:

      $ vault list secret/my-app/

  For a full list of examples and paths, please see the documentation that
  corresponds to the secret engine in use. Not all engines support listing.

  New options for recursive listing across all mounts:

      $ vault list --all-mounts --fuzzy api
          Fuzzy search all mounts for keys containing "api"

      $ vault list --all-mounts --pattern "*db*" secret/
          Glob pattern search across secret mount

      $ vault list --recursive secret/
          Recursively list all keys under secret/

      $ vault list --all-mounts --permissions=read
          Only show keys with read permission

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat | FlagSetOutputDetailed | FlagSetSnapshot)

	f := set.NewFlagSet("Command Options")

	f.BoolVar(&BoolVar{
		Name:       "all-mounts",
		Target:     &c.flagAllMounts,
		Usage:      "List secrets across all secret engine mounts, not just the specified path",
	})

	f.BoolVar(&BoolVar{
		Name:       "fuzzy",
		Target:     &c.flagFuzzy,
		Usage:      "Enable fuzzy (substring, case-insensitive) pattern matching",
	})

	f.StringVar(&StringVar{
		Name:   "pattern",
		Target: &c.flagPattern,
		Usage:  "Glob pattern for key matching (e.g., *api*, config/*/db)",
	})

	f.StringVar(&StringVar{
		Name:   "permissions",
		Target: &c.flagPermissions,
		Usage:  "Filter by capability (e.g., read, list, read,list)",
	})

	f.BoolVar(&BoolVar{
		Name:   "recursive",
		Target: &c.flagRecursive,
		Usage:  "Recursively list all keys under the path across all mounts",
	})

	f.StringVar(&StringVar{
		Name:   "recursive-mode",
		Target: &c.flagRecursiveMode,
		Usage:  "Recursive listing mode: 'api' uses the new list-recursive endpoint, 'cli' walks locally",
	})

	return set
}

func (c *ListCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *ListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ListCommand) Run(args []string) int {
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

	path := sanitizePath(args[0])
	var secret *api.Secret

	// Determine if we should use the new list-recursive API
	useRecursive := c.flagRecursive || c.flagAllMounts || c.flagFuzzy || c.flagPattern != "" || c.flagPermissions != ""

	if useRecursive {
		// Use the new ListRecursive API
		input := &api.ListRecursiveInput{
			Path:        path,
			Pattern:     c.flagPattern,
			Fuzzy:       c.flagFuzzy,
			Permissions: c.flagPermissions,
			AllMounts:   c.flagAllMounts,
		}
		secret, err = client.Logical().ListRecursive(context.Background(), input)
	} else if c.flagSnapshotID != "" {
		// Use snapshot-based listing
		secret, err = client.Logical().ListFromSnapshot(path, c.flagSnapshotID)
	} else {
		// Use standard listing
		secret, err = client.Logical().List(path)
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing %s: %s", path, err))
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

	if secret == nil {
		c.UI.Error(fmt.Sprintf("No value found at %s", path))
		return 2
	}
	if secret.Data == nil {
		// If secret wasn't nil, we have warnings, so output them anyways. We
		// may also have non-keys info.
		return OutputSecret(c.UI, secret)
	}

	if !ok {
		c.UI.Error(fmt.Sprintf("No entries found at %s", path))
		return 2
	}

	return OutputList(c.UI, secret)
}
