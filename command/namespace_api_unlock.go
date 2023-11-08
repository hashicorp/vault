// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*NamespaceAPIUnlockCommand)(nil)
	_ cli.CommandAutocomplete = (*NamespaceAPIUnlockCommand)(nil)
)

type NamespaceAPIUnlockCommand struct {
	*BaseCommand
}

func (c *NamespaceAPIUnlockCommand) Synopsis() string {
	return "Unlock the API for particular namespaces"
}

func (c *NamespaceAPIUnlockCommand) Help() string {
	helpText := `
Usage: vault namespace unlock [options] PATH

	Unlock the current namespace, and all descendants, with unlock key:

		$ vault namespace unlock -unlock-key=<key>

	Unlock the current namespace, and all descendants (from a root token):

		$ vault namespace unlock

	Unlock a child namespace, and all of its descendants (e.g. ns1/ns2/):

		$ vault namespace lock -unlock-key=<key> ns1/ns2

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceAPIUnlockCommand) Flags() *FlagSets {
	return c.FlagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *NamespaceAPIUnlockCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultNamespaces()
}

func (c *NamespaceAPIUnlockCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceAPIUnlockCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0 or 1, got %d)", len(args)))
		return 1
	}

	// current namespace is already encoded in the :ApiClient:
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	optionalChildNSPath := ""
	if len(args) == 1 {
		optionalChildNSPath = fmt.Sprintf("/%s", namespace.Canonicalize(args[0]))
	}

	_, err = client.Logical().Write(fmt.Sprintf("sys/namespaces/api-lock/unlock%s", optionalChildNSPath), map[string]interface{}{
		"unlock_key": c.flagUnlockKey,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error unlocking namespace: %v", err))
		return 2
	}

	return 0
}
