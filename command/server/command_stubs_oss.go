// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

import (
	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-hclog"
	base "github.com/hashicorp/vault/command/base"
	server "github.com/hashicorp/vault/helper/serverconfig"
	"github.com/hashicorp/vault/vault"
)

func entAdjustCoreConfig(config *server.Config, coreConfig *vault.CoreConfig) {
}

func entCheckStorageType(coreConfig *vault.CoreConfig) bool {
	return true
}

func entGetFIPSInfoKey() string {
	return ""
}

func entCheckRequestLimiter(_cmd *ServerCommand, _config *server.Config) {
}

func entExtendAddonHandlers(handlers *vaultHandlers) {}

func entCheckListenerConfig(_ hclog.Logger, core *vault.Core, config *server.Config) error {
	return nil
}

func entAugmentInfoKeys(config *server.Config, info map[string]string, infoKeys []string) {}

func entInitServerOnlyCommands(_ cli.Ui, _ *base.RunOptions, _ map[string]cli.CommandFactory) {}
