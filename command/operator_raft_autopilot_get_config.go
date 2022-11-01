package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftAutopilotGetConfigCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftAutopilotGetConfigCommand)(nil)
)

type OperatorRaftAutopilotGetConfigCommand struct {
	*BaseCommand
}

func (c *OperatorRaftAutopilotGetConfigCommand) Synopsis() string {
	return "Returns the configuration of the autopilot subsystem under integrated storage"
}

func (c *OperatorRaftAutopilotGetConfigCommand) Help() string {
	helpText := `
Usage: vault operator raft autopilot get-config

 Returns the configuration of the autopilot subsystem under integrated storage.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftAutopilotGetConfigCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *OperatorRaftAutopilotGetConfigCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftAutopilotGetConfigCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftAutopilotGetConfigCommand) Run(args []string) int {
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

	config, err := client.Sys().RaftAutopilotConfiguration()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if config == nil {
		return 0
	}

	if Format(c.UI) != "table" {
		return OutputData(c.UI, config)
	}

	entries := []string{"Key | Value"}
	entries = append(entries, fmt.Sprintf("%s | %t", "Cleanup Dead Servers", config.CleanupDeadServers))
	entries = append(entries, fmt.Sprintf("%s | %s", "Last Contact Threshold", config.LastContactThreshold.String()))
	entries = append(entries, fmt.Sprintf("%s | %s", "Dead Server Last Contact Threshold", config.DeadServerLastContactThreshold.String()))
	entries = append(entries, fmt.Sprintf("%s | %s", "Server Stabilization Time", config.ServerStabilizationTime.String()))
	entries = append(entries, fmt.Sprintf("%s | %d", "Min Quorum", config.MinQuorum))
	entries = append(entries, fmt.Sprintf("%s | %d", "Max Trailing Logs", config.MaxTrailingLogs))
	entries = append(entries, fmt.Sprintf("%s | %t", "Disable Upgrade Migration", config.DisableUpgradeMigration))

	return OutputData(c.UI, entries)
}
