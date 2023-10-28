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

var _ cli.Command = (*ConfigCurrentContext)(nil)

type ConfigCurrentContext struct {
	*BaseCommand
}

func (c *ConfigCurrentContext) Synopsis() string {
	return "returns the current client context if set"
}

func (c *ConfigCurrentContext) Help() string {
	helpText := `
Usage: vault config current-context

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigCurrentContext) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *ConfigCurrentContext) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigCurrentContext) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) > 0:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	currentConfig, err := config.LoadClientContextConfig("")
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to load client context config, %v", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("name: %s, address: %s, token: %s, namespace: %s",
		currentConfig.CurrentContext.Name,
		currentConfig.CurrentContext.VaultAddr,
		currentConfig.CurrentContext.ClusterToken,
		currentConfig.CurrentContext.NamespacePath))
	return 0
}
