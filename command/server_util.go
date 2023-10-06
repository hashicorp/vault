// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/vault"
)

var (
	// TODO remove once entAdjustCoreConfig has replaced it
	adjustCoreConfigForEnt = adjustCoreConfigForEntNoop
	// TODO remove once entCheckStorageType has replaced it
	storageSupportedForEnt = checkStorageTypeForEntNoop
)

func adjustCoreConfigForEntNoop(config *server.Config, coreConfig *vault.CoreConfig) {
}

// TODO remove once entGetFIPSInfoKey has replaced it
var getFIPSInfoKey = getFIPSInfoKeyNoop

func getFIPSInfoKeyNoop() string {
	return ""
}

func checkStorageTypeForEntNoop(coreConfig *vault.CoreConfig) bool {
	return true
}
