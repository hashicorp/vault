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

var _ cli.Command = (*ConfigSetContext)(nil)

type ConfigSetContext struct {
	*BaseCommand

	flagClientToken string
}

func (c *ConfigSetContext) Synopsis() string {
	return "Set client config context"
}

func (c *ConfigSetContext) Help() string {
	helpText := `
Usage: vault config set-context NAME [options]

  Set a context:

      $ vault config set-context vault_1 -addr=http://127.0.0.1:8200 -token=hvs. -namespace=ns1

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *ConfigSetContext) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "token",
		Target: &c.flagClientToken,
		Usage:  "Address of the Vault server.",
	})

	return set
}

func (c *ConfigSetContext) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *ConfigSetContext) Run(args []string) int {
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

	address := c.flagAddress
	if address == "" {
		c.UI.Warn("address is empty")
	}
	namespace := c.flagNamespace
	token := c.flagClientToken
	if token == "" {
		c.UI.Warn("token is empty")
	}

	newCtx := config.ContextInfo{
		Name:          name,
		VaultAddr:     address,
		ClusterToken:  token,
		NamespacePath: namespace,
	}

	currentConfig, err := config.LoadClientContextConfig("")
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to load client context config, %v", err))
		return 2
	}

	if len(currentConfig.ClientContexts) == 0 {
		currentConfig.ClientContexts = make([]config.ContextInfo, 0)
	}

	// update the client context with the new entry
	index, found := config.FindContextInfoIndexByName(currentConfig.ClientContexts, name)
	if found {
		currentConfig.ClientContexts[index] = newCtx
	} else {
		currentConfig.ClientContexts = append(currentConfig.ClientContexts, newCtx)
	}

	// set the current context to the newly set context
	currentConfig.CurrentContext = newCtx

	if err = config.WriteClientContextConfig("", currentConfig); err != nil {
		c.UI.Error(fmt.Sprintf("faile to write client context configuration, %v", err))
		return 2
	}

	c.UI.Output(fmt.Sprintf("Success! Set client config context: %s", name))
	return 0
}
