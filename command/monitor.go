package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*MonitorCommand)(nil)
var _ cli.CommandAutocomplete = (*MonitorCommand)(nil)

type MonitorCommand struct {
	*BaseCommand

	logLevel string

	// ShutdownCh is used to capture interrupt signal and end streaming
	ShutdownCh chan struct{}
}

func (c *MonitorCommand) Synopsis() string {
	return "Stream log messages from a Vault server"
}

func (c *MonitorCommand) Help() string {
	helpText := `
Usage: vault monitor [options]

	Stream log messages of a Vault server. The monitor command lets you listen
	for log levels that may be filtered out of the server logs. For example,
	the server may be logging at the INFO level, but with the monitor command
	you can set -log-level=DEBUG.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *MonitorCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Monitor Options")
	f.StringVar(&StringVar{
		Name:       "log-level",
		Target:     &c.logLevel,
		Default:    "INFO",
		Completion: complete.PredictSet("TRACE", "DEBUG", "INFO", "WARN", "ERROR"),
		Usage: "If passed, the log level to monitor logs. Supported values" +
			"(in order of detail) are \"TRACE\", \"DEBUG\", \"INFO\", \"WARN\"" +
			" and \"ERROR\".",
	})

	return set
}

func (c *MonitorCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *MonitorCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *MonitorCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	parsedArgs := f.Args()
	if len(parsedArgs) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(parsedArgs)))
		return 1
	}

	c.logLevel = strings.ToUpper(c.logLevel)
	validLevels := []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"}
	if !strutil.StrListContains(validLevels, c.logLevel) {
		c.UI.Error(fmt.Sprintf("%s is an unknown log level. Valid log levels are: %s", c.logLevel, validLevels))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	// Remove the default 60 second timeout so we can stream indefinitely
	client.SetClientTimeout(0)

	var logCh chan string
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logCh, err = client.Sys().Monitor(ctx, c.logLevel)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error starting monitor: %s", err))
		return 1
	}

	for {
		select {
		case log, ok := <-logCh:
			if !ok {
				return 0
			}
			c.UI.Info(log)
		case <-c.ShutdownCh:
			return 0
		}
	}
}
