package command

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorUsageCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorUsageCommand)(nil)

type OperatorUsageCommand struct {
	*BaseCommand
	flagStartTime time.Time
	flagEndTime   time.Time
}

func (c *OperatorUsageCommand) Synopsis() string {
	return "Lists historical client counts"
}

func (c *OperatorUsageCommand) Help() string {
	helpText := `
Usage: vault operator usage

  List the client counts for the default reporting period.

	  $ vault operator usage

  List the client counts for a specific time period.

          $ vault operator usage -start-time=2020-10 -end-time=2020-11

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorUsageCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.TimeVar(&TimeVar{
		Name:       "start-time",
		Usage:      "Start of report period (defaults to default_reporting_period before end time.)",
		Target:     &c.flagStartTime,
		Completion: complete.PredictNothing,
		Default:    time.Time{},
		Formats:    TimeVar_TimeOrDay | TimeVar_Month,
	})
	f.TimeVar(&TimeVar{
		Name:       "end-time",
		Usage:      "End of report period (defaults to end of last month.)",
		Target:     &c.flagEndTime,
		Completion: complete.PredictNothing,
		Default:    time.Time{},
		Formats:    TimeVar_TimeOrDay | TimeVar_Month,
	})

	return set
}

func (c *OperatorUsageCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorUsageCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorUsageCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	data := make(map[string][]string)
	if !c.flagStartTime.IsZero() {
		data["start_time"] = []string{c.flagStartTime.Format(time.RFC3339)}
	}
	if !c.flagEndTime.IsZero() {
		data["end_time"] = []string{c.flagEndTime.Format(time.RFC3339)}
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Logical().ReadWithData("sys/internal/counters/activity", data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving client counts: %v", err))
		return 2
	}

	if resp == nil || resp.Data == nil {
		c.UI.Warn("No data is available for the given time range.")
		// No further output
		// TODO: report if any data at all is available?
		return 0
	}

	switch Format(c.UI) {
	case "table":
	default:
		// Handle JSON, YAML, etc.
		return OutputData(c.UI, resp)
	}

	// Show this before the headers
	c.outputTimestamps(resp.Data)

	out := []string{
		"Namespace path | Distinct entities | Non-Entity tokens | Clients",
	}

	out = c.addNamespacesToOutput(out, resp.Data)
	out = c.addTotalToOutput(out, resp.Data)

	c.UI.Output(tableOutput(out, nil))
	return 0
}

func (c *OperatorUsageCommand) outputTimestamps(data map[string]interface{}) {
	c.UI.Output(fmt.Sprintf("Period start: %v\nPeriod end: %v\n",
		data["start_time"].(string),
		data["end_time"].(string)))
}

func (c *OperatorUsageCommand) addNamespacesToOutput(out []string, data map[string]interface{}) []string {
	byNs := data["by_namespace"].([]interface{})

	// TODO: provide a function in the API module for doing this conversion?
	for _, rawVal := range byNs {
		val := rawVal.(map[string]interface{})
		namespacePath := val["namespace_path"].(string)
		counts := val["counts"].(map[string]interface{})

		// TODO: check errors
		entityCount, _ := counts["distinct_entities"].(json.Number).Int64()
		tokenCount, _ := counts["non_entity_tokens"].(json.Number).Int64()
		clientCount, _ := counts["clients"].(json.Number).Int64()
		if namespacePath == "" {
			namespacePath = "[root]"
		}
		out = append(out, fmt.Sprintf("%s | %d | %d | %d", namespacePath, entityCount, tokenCount, clientCount))
	}
	return out
}

func (c *OperatorUsageCommand) addTotalToOutput(out []string, data map[string]interface{}) []string {
	// blank line separating it from namespaces
	out = append(out, " - | - | - | - ")

	total := data["total"].(map[string]interface{})
	entityCount, _ := total["distinct_entities"].(json.Number).Int64()
	tokenCount, _ := total["non_entity_tokens"].(json.Number).Int64()
	clientCount, _ := total["clients"].(json.Number).Int64()
	out = append(out, fmt.Sprintf("Total | %d | %d | %d", entityCount, tokenCount, clientCount))
	return out
}
