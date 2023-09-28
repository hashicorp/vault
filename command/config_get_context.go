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

var _ cli.Command = (*ConfigGetContext)(nil)

type ConfigGetContext struct {
	*BaseCommand
}

func (c *ConfigGetContext) Synopsis() string {
	return "get the client config context for the given context name"
}

func (c *ConfigGetContext) Help() string {
	helpText := `
Usage: vault config get-context NAME

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigGetContext) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *ConfigGetContext) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigGetContext) Run(args []string) int {
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

	ctxInfo := currentConfig.ClientContexts[index]

	c.UI.Output(fmt.Sprintf("name: %s, address: %s, token: %s, namespace: %s",
		ctxInfo.Name,
		ctxInfo.VaultAddr,
		ctxInfo.ClusterToken,
		ctxInfo.NamespacePath))
	return 0
}
