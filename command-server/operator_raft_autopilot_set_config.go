// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command_server

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/command"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftAutopilotSetConfigCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftAutopilotSetConfigCommand)(nil)
)

type OperatorRaftAutopilotSetConfigCommand struct {
	*command.BaseCommand
	flagCleanupDeadServers             command.BoolPtr
	flagLastContactThreshold           time.Duration
	flagDeadServerLastContactThreshold time.Duration
	flagMaxTrailingLogs                uint64
	flagMinQuorum                      uint
	flagServerStabilizationTime        time.Duration
	flagDisableUpgradeMigration        command.BoolPtr
	flagDRToken                        string
}

func (c *OperatorRaftAutopilotSetConfigCommand) Synopsis() string {
	return "Modify the configuration of the autopilot subsystem under integrated storage"
}

func (c *OperatorRaftAutopilotSetConfigCommand) Help() string {
	helpText := `
Usage: vault operator raft autopilot set-config [options]

  Modify the configuration of the autopilot subsystem under integrated storage.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftAutopilotSetConfigCommand) Flags() *command.FlagSets {
	set := c.FlagSet(command.FlagSetHTTP | command.FlagSetOutputFormat)

	f := set.NewFlagSet("Common Options")

	f.BoolPtrVar(&command.BoolPtrVar{
		Name:   "cleanup-dead-servers",
		Target: &c.flagCleanupDeadServers,
		Usage:  "Controls whether to remove dead servers from the Raft peer list periodically or when a new server joins.",
	})

	f.DurationVar(&command.DurationVar{
		Name:   "last-contact-threshold",
		Target: &c.flagLastContactThreshold,
		Usage:  "Limit on the amount of time a server can go without leader contact before being considered unhealthy.",
	})

	f.DurationVar(&command.DurationVar{
		Name:   "dead-server-last-contact-threshold",
		Target: &c.flagDeadServerLastContactThreshold,
		Usage:  "Limit on the amount of time a server can go without leader contact before being considered failed. This takes effect only when cleanup_dead_servers is set.",
	})

	f.Uint64Var(&command.Uint64Var{
		Name:   "max-trailing-logs",
		Target: &c.flagMaxTrailingLogs,
		Usage:  "Amount of entries in the Raft Log that a server can be behind before being considered unhealthy.",
	})

	f.UintVar(&command.UintVar{
		Name:   "min-quorum",
		Target: &c.flagMinQuorum,
		Usage:  "Minimum number of servers allowed in a cluster before autopilot can prune dead servers. This should at least be 3.",
	})

	f.DurationVar(&command.DurationVar{
		Name:   "server-stabilization-time",
		Target: &c.flagServerStabilizationTime,
		Usage:  "Minimum amount of time a server must be in a stable, healthy state before it can be added to the cluster.",
	})

	f.BoolPtrVar(&command.BoolPtrVar{
		Name:   "disable-upgrade-migration",
		Target: &c.flagDisableUpgradeMigration,
		Usage:  "Whether or not to perform automated version upgrades.",
	})

	f.StringVar(&command.StringVar{
		Name:       "dr-token",
		Target:     &c.flagDRToken,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage:      "DR operation token used to authorize this request (if a DR secondary node).",
	})

	return set
}

func (c *OperatorRaftAutopilotSetConfigCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftAutopilotSetConfigCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftAutopilotSetConfigCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch len(args) {
	case 0:
	default:
		c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	data := make(map[string]interface{})
	if c.flagCleanupDeadServers.IsSet() {
		data["cleanup_dead_servers"] = c.flagCleanupDeadServers.Get()
	}
	if c.flagMaxTrailingLogs > 0 {
		data["max_trailing_logs"] = c.flagMaxTrailingLogs
	}
	if c.flagMinQuorum > 0 {
		data["min_quorum"] = c.flagMinQuorum
	}
	if c.flagLastContactThreshold > 0 {
		data["last_contact_threshold"] = c.flagLastContactThreshold.String()
	}
	if c.flagDeadServerLastContactThreshold > 0 {
		data["dead_server_last_contact_threshold"] = c.flagDeadServerLastContactThreshold.String()
	}
	if c.flagServerStabilizationTime > 0 {
		data["server_stabilization_time"] = c.flagServerStabilizationTime.String()
	}
	if c.flagDisableUpgradeMigration.IsSet() {
		data["disable_upgrade_migration"] = c.flagDisableUpgradeMigration.Get()
	}
	if c.flagDRToken != "" {
		data["dr_operation_token"] = c.flagDRToken
	}

	secret, err := client.Logical().Write("sys/storage/raft/autopilot/configuration", data)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	if secret == nil {
		return 0
	}

	return command.OutputSecret(c.UI, secret)
}
