// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

import (
	"github.com/hashicorp/cli"
	base "github.com/hashicorp/vault/command/base"
)

func entInitCommands(ui, serverCmdUi cli.Ui, runOpts *base.RunOptions, commands map[string]cli.CommandFactory) {
}
