package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorRaftAutopilotSetConfigCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorRaftAutopilotSetConfigCommand)(nil)

type OperatorRaftAutopilotSetConfigCommand struct {
	*BaseCommand
	flagCleanupDeadServers          bool
	flagLastContactThreshold        time.Duration
	flagLastContactFailureThreshold time.Duration
	flagMaxTrailingLogs             uint64
	flagMinQuorum                   uint
	flagServerStabilizationTime     time.Duration
}

func (c *OperatorRaftAutopilotSetConfigCommand) Synopsis() string {
	return "Modify the configuration of the autopilot subsystem under integrated storage"
}

func (c *OperatorRaftAutopilotSetConfigCommand) Help() string {
	helpText := `
Usage: vault operator raft snapshot save <snapshot_file>

  Modify the configuration of the autopilot subsystem under integrated storage.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftAutopilotSetConfigCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Common Options")

	f.BoolVar(&BoolVar{
		Name:    "cleanup-dead-servers",
		Target:  &c.flagCleanupDeadServers,
		Default: false,
	})

	f.DurationVar(&DurationVar{
		Name:    "last-contact-threshold",
		Target:  &c.flagLastContactThreshold,
		Default: 10 * time.Second,
	})

	f.DurationVar(&DurationVar{
		Name:    "last-contact-failure-threshold",
		Target:  &c.flagLastContactFailureThreshold,
		Default: 24 * time.Hour,
	})

	f.Uint64Var(&Uint64Var{
		Name:    "max-trailing-logs",
		Target:  &c.flagMaxTrailingLogs,
		Default: 1000,
	})

	f.UintVar(&UintVar{
		Name:    "min-quorum",
		Target:  &c.flagMinQuorum,
		Default: 3,
	})

	f.DurationVar(&DurationVar{
		Name:    "server-stabilization-time",
		Target:  &c.flagServerStabilizationTime,
		Default: 10 * time.Second,
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

	secret, err := client.Logical().Write("sys/storage/raft/autopilot/configuration", map[string]interface{}{
		"cleanup_dead_servers":           c.flagCleanupDeadServers,
		"max_trailing_logs":              c.flagMaxTrailingLogs,
		"min_quorum":                     c.flagMinQuorum,
		"last_contact_threshold":         c.flagLastContactThreshold.String(),
		"last_contact_failure_threshold": c.flagLastContactFailureThreshold.String(),
		"server_stabilization_time":      c.flagServerStabilizationTime.String(),
	})
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	if secret == nil {
		return 0
	}

	return OutputSecret(c.UI, secret)
}
