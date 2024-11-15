// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package vaulthcplib

import "github.com/hashicorp/cli"

func InitHCPCommand(ui cli.Ui) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"hcp connect": func() (cli.Command, error) {
			return &HCPConnectCommand{
				Ui: ui,
			}, nil
		},
		"hcp disconnect": func() (cli.Command, error) {
			return &HCPDisconnectCommand{
				Ui: ui,
			}, nil
		},
	}
}
