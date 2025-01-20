// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

import (
	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
)

func entInitCommands(ui, serverCmdUi cli.Ui, runOpts *RunOptions, commands map[string]cli.CommandFactory) {
}

func entEnableFourClusterDev(c *ServerCommand, base *vault.CoreConfig, info map[string]string, infoKeys []string, tempDir string) int {
	c.logger.Error("-dev-four-cluster only supported in enterprise Vault")
	return 1
}

func entAdjustCoreConfig(config *server.Config, coreConfig *vault.CoreConfig) {
}

func entCheckStorageType(coreConfig *vault.CoreConfig) bool {
	return true
}

func entGetFIPSInfoKey() string {
	return ""
}

func entGetRequestLimiterStatus(coreConfig vault.CoreConfig) string {
	return ""
}

func entExtendAddonHandlers(handlers *vaultHandlers) {}
