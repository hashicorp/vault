// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/mitchellh/go-homedir"

	"github.com/hashicorp/cli"
)

var _ cli.Command = (*HCPDisconnectCommand)(nil)

type HCPDisconnectCommand struct {
	Ui cli.Ui
}

func (c *HCPDisconnectCommand) Help() string {
	helpText := `
Usage: vault hcp disconnect [options]
  
  Cleans up the cache with the HCP credentials used to connect to a HCP Vault cluster. 

      $ vault hcp disconnect
`
	return strings.TrimSpace(helpText)
}

func (c *HCPDisconnectCommand) Run(_ []string) int {
	path, err := homedir.Dir()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("\nFailed to find home directory: %s", err))
		return 1
	}

	if err := eraseConfig(path); err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to disconnect from HCP Vault Cluster: %s", err))
		return 1
	}

	opts := []config.HCPConfigOption{config.FromEnv()}
	cfg, err := config.NewHCPConfig(opts...)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to disconnect from HCP Vault Cluster: %s", err))
		return 1
	}

	err = cfg.Logout()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to disconnect from HCP Vault Cluster: %s", err))
		return 1
	}

	return 0
}

func (c *HCPDisconnectCommand) Synopsis() string {
	return "Disconnect from the HCP Vault Cluster"
}
