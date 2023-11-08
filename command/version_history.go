// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/ryanuber/columnize"
)

var (
	_ cli.Command             = (*VersionHistoryCommand)(nil)
	_ cli.CommandAutocomplete = (*VersionHistoryCommand)(nil)
)

// VersionHistoryCommand is a Command implementation prints the version.
type VersionHistoryCommand struct {
	*BaseCommand
}

func (c *VersionHistoryCommand) Synopsis() string {
	return "Prints the version history of the target Vault server"
}

func (c *VersionHistoryCommand) Help() string {
	helpText := `
Usage: vault version-history

  Prints the version history of the target Vault server.

  Print the version history:

      $ vault version-history
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *VersionHistoryCommand) Flags() *FlagSets {
	return c.FlagSet(FlagSetOutputFormat)
}

func (c *VersionHistoryCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *VersionHistoryCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

const versionTrackingWarning = `Note:
Use of this command requires a server running Vault 1.10.0 or greater.
Version tracking was added in 1.9.0. Earlier versions have not been tracked.
`

func (c *VersionHistoryCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Logical().List("sys/version-history")
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading version history: %s", err))
		return 2
	}

	if resp == nil || resp.Data == nil {
		c.UI.Error("Invalid response returned from Vault")
		return 2
	}

	if c.FlagFormat == "json" {
		c.UI.Warn("")
		c.UI.Warn(versionTrackingWarning)
		c.UI.Warn("")

		return OutputData(c.UI, resp)
	}

	var keyInfo map[string]interface{}

	keys, ok := extractListData(resp)
	if !ok {
		c.UI.Error("Expected keys in response to be an array")
		return 2
	}

	keyInfo, ok = resp.Data["key_info"].(map[string]interface{})
	if !ok {
		c.UI.Error("Expected key_info in response to be a map")
		return 2
	}

	table := []string{"Version | Installation Time | Build Date"}
	columnConfig := columnize.DefaultConfig()

	for _, versionRaw := range keys {
		version, ok := versionRaw.(string)

		if !ok {
			c.UI.Error("Expected version to be string")
			return 2
		}

		versionInfoRaw := keyInfo[version]

		versionInfo, ok := versionInfoRaw.(map[string]interface{})
		if !ok {
			c.UI.Error(fmt.Sprintf("Expected version info for %q to be map", version))
			return 2
		}

		table = append(table, fmt.Sprintf("%s | %s | %s", version, versionInfo["timestamp_installed"], versionInfo["build_date"]))
	}

	c.UI.Warn("")
	c.UI.Warn(versionTrackingWarning)
	c.UI.Warn("")
	c.UI.Output(TableOutput(table, columnConfig))

	return 0
}
