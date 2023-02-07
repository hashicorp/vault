package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftAutopilotSetConfigCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftAutopilotSetConfigCommand)(nil)
)

type OperatorRaftAutopilotSetConfigCommand struct {
	*BaseCommand
	flagCleanupDeadServers             BoolPtr
	flagLastContactThreshold           time.Duration
	flagDeadServerLastContactThreshold time.Duration
	flagMaxTrailingLogs                uint64
	flagMinQuorum                      uint
	flagServerStabilizationTime        time.Duration
	flagDisableUpgradeMigration        BoolPtr
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

func (c *OperatorRaftAutopilotSetConfigCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Common Options")

	f.BoolPtrVar(&BoolPtrVar{
		Name:   "cleanup-dead-servers",
		Target: &c.flagCleanupDeadServers,
	})

	f.DurationVar(&DurationVar{
		Name:   "last-contact-threshold",
		Target: &c.flagLastContactThreshold,
	})

	f.DurationVar(&DurationVar{
		Name:   "dead-server-last-contact-threshold",
		Target: &c.flagDeadServerLastContactThreshold,
	})

	f.Uint64Var(&Uint64Var{
		Name:   "max-trailing-logs",
		Target: &c.flagMaxTrailingLogs,
	})

	f.UintVar(&UintVar{
		Name:   "min-quorum",
		Target: &c.flagMinQuorum,
	})

	f.DurationVar(&DurationVar{
		Name:   "server-stabilization-time",
		Target: &c.flagServerStabilizationTime,
	})

	f.BoolPtrVar(&BoolPtrVar{
		Name:   "disable-upgrade-migration",
		Target: &c.flagDisableUpgradeMigration,
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

	secret, err := client.Logical().Write("sys/storage/raft/autopilot/configuration", data)
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	if secret == nil {
		return 0
	}

	return OutputSecret(c.UI, secret)
}
