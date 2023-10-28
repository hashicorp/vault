// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/vault/command/config"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*ConfigDeleteContext)(nil)

type ConfigDeleteContext struct {
	*BaseCommand
}

func (c *ConfigDeleteContext) Synopsis() string {
	return "Delete a client config context with the given name"
}

func (c *ConfigDeleteContext) Help() string {
	helpText := `
Usage: vault config delete-context NAME

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigDeleteContext) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *ConfigDeleteContext) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigDeleteContext) Run(args []string) int {
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

	name := strings.TrimSpace(args[0])
	if name == "" {
		c.UI.Error("cannot get a client config context with an empty name")
		return 2
	}

	currentConfig, err := config.LoadClientContextConfig("")
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to load client context config, %v", err))
		return 2
	}

	index, found := config.FindContextInfoIndexByName(currentConfig.ClientContexts, name)
	if !found {
		c.UI.Error(fmt.Sprintf("failed to find the given context, %v", name))
		return 2
	}

	// check if the current context is set to the one that we need to delete
	if currentConfig.CurrentContext.Name == name {
		currentConfig.CurrentContext = config.ContextInfo{}
	}

	currentConfig.ClientContexts = slices.Delete(currentConfig.ClientContexts, index, index+1)

	if err = config.WriteClientContextConfig("", currentConfig); err != nil {
		c.UI.Error(fmt.Sprintf("failed to write back the client config context, %v", name))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! removed the given client config context (if existed), %s", name))
	return 0
}
