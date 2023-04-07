// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*AgentGenerateConfigCommand)(nil)
	_ cli.CommandAutocomplete = (*AgentGenerateConfigCommand)(nil)
)

type AgentGenerateConfigCommand struct {
	*BaseCommand

	flagType  string
	flagPaths []string
	flagExec  string
}

func (c *AgentGenerateConfigCommand) Synopsis() string {
	return "Generate a Vault Agent configuration file."
}

func (c *AgentGenerateConfigCommand) Help() string {
	helpText := `
Usage: vault agent generate-config [options]
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *AgentGenerateConfigCommand) Flags() *FlagSets {
	set := NewFlagSets(c.UI)

	// Common Options
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "type",
		Target:  &c.flagType,
		Default: "env-template",
		Usage:   "The type of configuration file to generate; currently, only 'env-template' is supported.",
		Completion: complete.PredictSet(
			"env-template",
		),
	})

	f.StringSliceVar(&StringSliceVar{
		Name:       "path",
		Target:     &c.flagPaths,
		Usage:      "Path to a KV v1/v2 secret (e.g. secret/data/foo, secret/prefix/*).",
		Completion: c.PredictVaultFolders(),
	})

	f.StringVar(&StringVar{
		Name:    "exec",
		Target:  &c.flagExec,
		Default: "env",
		Usage:   "The command to execute for in env-template mode.",
	})

	return set
}

func (c *AgentGenerateConfigCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *AgentGenerateConfigCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *AgentGenerateConfigCommand) Run(args []string) int {
	ctx, cancelContextFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelContextFunc()

	flags := c.Flags()

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = flags.Args()

	if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected at most 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	var envTemplates []*config.EnvTemplateConfig

	for _, path := range c.flagPaths {
		pathSanitized := sanitizePath(path)
		pathMount, v2, err := isKVv2(pathSanitized, client)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Could not validate secret path %q: %v", path, err))
			return 2
		}

		var pathFull string
		if v2 {
			pathFull = addPrefixToKVPath(pathSanitized, pathMount, "data")
		} else {
			pathFull = pathSanitized
		}

		resp, err := client.Logical().ReadWithContext(ctx, pathFull)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error querying %q: %v", pathFull, err))
			return 2
		}
		if resp == nil {
			c.UI.Error(fmt.Sprintf("Secret not found at %q", pathFull))
			return 2
		}

		for field := range resp.Data {
			envTemplates = append(envTemplates, &config.EnvTemplateConfig{
				Name:     strings.ToUpper(field),
				Contents: fmt.Sprintf("{{ with secret %s }}{{ Data.data.%s }}", pathFull, field),
			})
		}
	}

	var execCommand string
	if c.flagExec != "" {
		execCommand = c.flagExec
	} else {
		execCommand = "env"
	}

	agentConfig := config.ConfigGen{
		Vault: &config.VaultGen{
			Address: client.Address(),
		},
		AutoAuth: &config.AutoAuthGen{
			Method: &config.AutoAuthMethodGen{
				Type: "token_file",
				Config: config.AutoAuthMethodConfigGen{
					TokenFilePath: "$HOME/.vault-token",
				},
			},
		},
		EnvTemplates: envTemplates,
		Exec: &config.ExecConfig{
			Command:            execCommand,
			Args:               []string{},
			RestartOnNewSecret: "always",
			RestartKillSignal:  "SIGTERM",
		},
	}

	var configPath string
	if len(args) == 1 {
		configPath = args[0]
	} else {
		configPath = "agent.hcl"
	}

	contents := hclwrite.NewEmptyFile()

	gohcl.EncodeIntoBody(&agentConfig, contents.Body())

	f, err := os.Create(configPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not create configuration file %q: %v", configPath, err))
		return 1
	}
	defer func() {
		if err := f.Close(); err != nil {
			c.UI.Error(fmt.Sprintf("Could not close configuration file %q: %v", configPath, err))
		}
	}()

	if _, err := contents.WriteTo(f); err != nil {
		c.UI.Error(fmt.Sprintf("Could not write to configuration file %q: %v", configPath, err))
		return 1
	}

	c.UI.Info(fmt.Sprintf("Successfully generated %q configuration file!", configPath))

	return 0
}
