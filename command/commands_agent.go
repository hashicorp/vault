// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build agent

package command

import "github.com/hashicorp/cli"

func extendServerCommands(commands map[string]cli.CommandFactory, serverCmdUi cli.Ui, runOpts *RunOptions, handlers *vaultHandlers) {
	// No-op
}
