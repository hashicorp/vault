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

var _ cli.Command = (*ConfigUseContext)(nil)

type ConfigUseContext struct {
	*BaseCommand
}

func (c *ConfigUseContext) Synopsis() string {
	return "Disables an auth method"
}

func (c *ConfigUseContext) Help() string {
	helpText := `
Usage: vault config use-context name

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigUseContext) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *ConfigUseContext) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigUseContext) Run(args []string) int {
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
		c.UI.Error("cannot set a client config context with an empty name")
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

	// set the current context to the found client config context
	currentConfig.CurrentContext = currentConfig.ClientContexts[index]

	if err = config.WriteClientContextConfig("", currentConfig); err != nil {
		c.UI.Error(fmt.Sprintf("faile to write client context configuration, %v", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Set current client config context to: %s", name))
	return 0
}
