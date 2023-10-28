// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/command/config"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*ConfigUnset)(nil)

type ConfigUnset struct {
	*BaseCommand
}

func (c *ConfigUnset) Synopsis() string {
	return "unset an entry in a context or unset the current-context"
}

func (c *ConfigUnset) Help() string {
	helpText := `
Usage: vault config unset current-context

To remove an entry in a context:
      $ vault config unset contexts.vault_1.namespace

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigUnset) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *ConfigUnset) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigUnset) Run(args []string) int {
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

	ctxEntry := strings.TrimSpace(args[0])
	if ctxEntry == "" {
		c.UI.Error("cannot unset an a field in a client config context with an empty entry")
		return 2
	}

	splitEntry := strings.Split(ctxEntry, ".")
	var name, entry string
	switch len(splitEntry) {
	case 3:
		for i, e := range splitEntry {
			if strings.TrimSpace(e) == "" {
				c.UI.Error(fmt.Sprintf("invalid entry at index %d. entry cannot be empty", i))
				return 2
			}
		}
		name = splitEntry[1]
		entry = strings.ToLower(splitEntry[2])
	case 1:
		if splitEntry[0] != "current-context" {
			c.UI.Error("only 'current-context' can be unset as a whole")
			return 2
		}
	default:
		c.UI.Error("invalid entry to unset")
		return 2
	}

	currentConfig, err := config.LoadClientContextConfig("")
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to load client context config, %v", err))
		return 2
	}

	if len(splitEntry) == 1 {
		// we have validated the entry, and we need to unset the current-context
		currentConfig.CurrentContext = config.ContextInfo{}
	} else {
		index, found := config.FindContextInfoIndexByName(currentConfig.ClientContexts, name)
		if !found {
			c.UI.Error(fmt.Sprintf("failed to find the given context, %v", name))
			return 2
		}

		switch entry {
		// case "name": this is invalid. a context should always have a name
		case "token":
			currentConfig.ClientContexts[index].ClusterToken = ""
		case "address":
			currentConfig.ClientContexts[index].VaultAddr = ""
		case "namespace":
			currentConfig.ClientContexts[index].NamespacePath = ""
		default:
			c.UI.Error("invalid entry name to unset")
			return 2
		}
	}

	if err = config.WriteClientContextConfig("", currentConfig); err != nil {
		c.UI.Error(fmt.Sprintf("faile to write client context configuration, %v", err))
		return 2
	}

	c.UI.Output("Success! Unset the given entry in the client config context")
	return 0
}
