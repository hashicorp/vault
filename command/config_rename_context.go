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

var _ cli.Command = (*ConfigRenameContext)(nil)

type ConfigRenameContext struct {
	*BaseCommand
}

func (c *ConfigRenameContext) Synopsis() string {
	return "Renames a client config context with the given name"
}

func (c *ConfigRenameContext) Help() string {
	helpText := `
Usage: vault config rename-context NAME NEW-NAME
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigRenameContext) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *ConfigRenameContext) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigRenameContext) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}

	currentName := strings.TrimSpace(args[0])
	if currentName == "" {
		c.UI.Error("cannot get a client config context with an empty name")
		return 2
	}

	newName := strings.TrimSpace(args[1])
	if newName == "" {
		c.UI.Error("cannot replace a client config context with an empty new name")
		return 2
	}

	currentConfig, err := config.LoadClientContextConfig("")
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to load client context config, %v", err))
		return 2
	}

	index, found := config.FindContextInfoIndexByName(currentConfig.ClientContexts, currentName)
	if !found {
		c.UI.Error(fmt.Sprintf("failed to find a client config context with the given name %q", currentName))
		return 2
	}

	// rename the client config context
	currentConfig.ClientContexts[index].Name = newName

	// make sure to rename the current context name if its name matches the current name
	if currentConfig.CurrentContext.Name == currentName {
		currentConfig.CurrentContext.Name = newName
	}

	if err = config.WriteClientContextConfig("", currentConfig); err != nil {
		c.UI.Error(fmt.Sprintf("failed to write client config context, error: %v", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Renamed a client config context name (if it existed): old name:%s, new name: %s", currentName, newName))
	return 0
}
