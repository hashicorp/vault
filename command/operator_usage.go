// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
	"github.com/ryanuber/columnize"
)

var (
	_ cli.Command             = (*OperatorUsageCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorUsageCommand)(nil)
)

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
		Usage:      "Start of report period. Defaults to 'default_reporting_period' before end time.",
		Target:     &c.flagStartTime,
		Completion: complete.PredictNothing,
		Default:    time.Time{},
		Formats:    TimeVar_TimeOrDay | TimeVar_Month,
	})
	f.TimeVar(&TimeVar{
		Name:       "end-time",
		Usage:      "End of report period. Defaults to end of last month.",
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
		if c.noReportAvailable(client) {
			c.UI.Warn("Vault does not have any usage data available. A report will be available\n" +
				"after the first calendar month in which monitoring is enabled.")
		} else {
			c.UI.Warn("No data is available for the given time range.")
		}
		// No further output
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
		"Namespace path | Distinct entities | Non-Entity tokens | Secret syncs | ACME clients | Active clients",
	}

	out = append(out, c.namespacesOutput(resp.Data)...)
	out = append(out, c.totalOutput(resp.Data)...)

	colConfig := columnize.DefaultConfig()
	colConfig.Empty = " " // Do not show n/a on intentional blank lines
	colConfig.Glue = "   "
	c.UI.Output(tableOutput(out, colConfig))
	return 0
}

// noReportAvailable checks whether we can definitively say that no
// queries can be answered; if there's an error, just fall back to
// reporting that the response is empty.
func (c *OperatorUsageCommand) noReportAvailable(client *api.Client) bool {
	if c.flagOutputCurlString || c.flagOutputPolicy {
		// Don't mess up the original query string
		return false
	}

	resp, err := client.Logical().Read("sys/internal/counters/config")
	if err != nil || resp == nil || resp.Data == nil {
		c.UI.Warn("bad response from config")
		return false
	}

	qaRaw, ok := resp.Data["queries_available"]
	if !ok {
		c.UI.Warn("no queries_available key")
		return false
	}

	qa, ok := qaRaw.(bool)
	if !ok {
		c.UI.Warn("wrong type")
		return false
	}

	return !qa
}

func (c *OperatorUsageCommand) outputTimestamps(data map[string]interface{}) {
	c.UI.Output(fmt.Sprintf("Period start: %v\nPeriod end: %v\n",
		data["start_time"].(string),
		data["end_time"].(string)))
}

type UsageCommandNamespace struct {
	formattedLine string
	sortOrder     string

	// Sort order:
	// -- root first
	// -- namespaces in lexicographic order
	// -- deleted namespace "xxxxx" last
}

type UsageResponse struct {
	namespacePath string
	entityCount   int64
	// As per 1.9, the tokenCount field will contain the distinct non-entity
	// token clients instead of each individual token.
	tokenCount  int64
	secretSyncs int64
	acmeCount   int64
	clientCount int64
}

func jsonNumberOK(m map[string]interface{}, key string) (int64, bool) {
	val, ok := m[key].(json.Number)
	if !ok {
		return 0, false
	}
	intVal, err := val.Int64()
	if err != nil {
		return 0, false
	}
	return intVal, true
}

// TODO: provide a function in the API module for doing this conversion?
func (c *OperatorUsageCommand) parseNamespaceCount(rawVal interface{}) (UsageResponse, error) {
	var ret UsageResponse

	val, ok := rawVal.(map[string]interface{})
	if !ok {
		return ret, errors.New("value is not a map")
	}

	ret.namespacePath, ok = val["namespace_path"].(string)
	if !ok {
		return ret, errors.New("bad namespace path")
	}

	counts, ok := val["counts"].(map[string]interface{})
	if !ok {
		return ret, errors.New("missing counts")
	}

	ret.entityCount, ok = jsonNumberOK(counts, "distinct_entities")
	if !ok {
		return ret, errors.New("missing distinct_entities")
	}

	ret.tokenCount, ok = jsonNumberOK(counts, "non_entity_tokens")
	if !ok {
		return ret, errors.New("missing non_entity_tokens")
	}

	// don't error if the secret syncs key is missing
	ret.secretSyncs, _ = jsonNumberOK(counts, "secret_syncs")

	// don't error if acme clients is missing
	ret.acmeCount, _ = jsonNumberOK(counts, "acme_clients")

	ret.clientCount, ok = jsonNumberOK(counts, "clients")
	if !ok {
		return ret, errors.New("missing clients")
	}

	return ret, nil
}

func (c *OperatorUsageCommand) namespacesOutput(data map[string]interface{}) []string {
	byNs, ok := data["by_namespace"].([]interface{})
	if !ok {
		c.UI.Error("missing namespace breakdown in response")
		return nil
	}

	nsOut := make([]UsageCommandNamespace, 0, len(byNs))

	for _, rawVal := range byNs {
		val, err := c.parseNamespaceCount(rawVal)
		if err != nil {
			c.UI.Error(fmt.Sprintf("malformed namespace in response: %v", err))
			continue
		}

		sortOrder := "1" + val.namespacePath
		if val.namespacePath == "" {
			val.namespacePath = "[root]"
			sortOrder = "0"
		} else if strings.HasPrefix(val.namespacePath, "deleted namespace") {
			sortOrder = "2" + val.namespacePath
		}

		formattedLine := fmt.Sprintf("%s | %d | %d | %d | %d | %d",
			val.namespacePath, val.entityCount, val.tokenCount, val.secretSyncs, val.acmeCount, val.clientCount)
		nsOut = append(nsOut, UsageCommandNamespace{
			formattedLine: formattedLine,
			sortOrder:     sortOrder,
		})
	}

	sort.Slice(nsOut, func(i, j int) bool {
		return nsOut[i].sortOrder < nsOut[j].sortOrder
	})

	out := make([]string, len(nsOut))
	for i := range nsOut {
		out[i] = nsOut[i].formattedLine
	}

	return out
}

func (c *OperatorUsageCommand) totalOutput(data map[string]interface{}) []string {
	// blank line separating it from namespaces
	out := []string{"  |  |  |  |  |  "}

	total, ok := data["total"].(map[string]interface{})
	if !ok {
		c.UI.Error("missing total in response")
		return out
	}

	entityCount, ok := jsonNumberOK(total, "distinct_entities")
	if !ok {
		c.UI.Error("missing distinct_entities in total")
		return out
	}

	tokenCount, ok := jsonNumberOK(total, "non_entity_tokens")
	if !ok {
		c.UI.Error("missing non_entity_tokens in total")
		return out
	}
	// don't error if secret syncs key is missing
	secretSyncs, _ := jsonNumberOK(total, "secret_syncs")

	// don't error if acme clients is missing
	acmeCount, _ := jsonNumberOK(total, "acme_clients")

	clientCount, ok := jsonNumberOK(total, "clients")
	if !ok {
		c.UI.Error("missing clients in total")
		return out
	}

	out = append(out, fmt.Sprintf("Total | %d | %d | %d | %d | %d",
		entityCount, tokenCount, secretSyncs, acmeCount, clientCount))
	return out
}
