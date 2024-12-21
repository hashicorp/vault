// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package command

import "github.com/hashicorp/cli"

func NewPluginRegisterCommand(baseCommand *BaseCommand) cli.Command {
	return &PluginRegisterCommand{
		BaseCommand: baseCommand,
	}
}
