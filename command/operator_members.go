package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorMembersCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorMembersCommand)(nil)

type OperatorMembersCommand struct {
	*BaseCommand
}

func (c *OperatorMembersCommand) Synopsis() string {
	return "Returns the list of Nodes in the cluster"
}

func (c *OperatorMembersCommand) Help() string {
	helpText := `
Usage: vault operator members
  Provides the details of all the nodes in the cluster.
	  $ vault operator members
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorMembersCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	return set
}

func (c *OperatorMembersCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorMembersCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorMembersCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	r := client.NewRequest("GET", "/v1/sys/ha-status")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := client.RawRequestWithContext(ctx, r)
	if err != nil {
		return 1
	}
	defer resp.Body.Close()

	var result HaStatusResponse
	err = resp.DecodeJSON(&result)

	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	switch Format(c.UI) {
	case "table":
		leader := result.Leader
		nodes := result.Standby
		out := []string{"Host Name | API Address | Cluster Address | Type"}
		out = append(out, fmt.Sprintf("%s | %s | %s | %s", leader.HostName, leader.ApiAddress, leader.ClusterAddress, "Active"))
		for _, node := range nodes {
			out = append(out, fmt.Sprintf("%s | %s | %s | %s", node.HostName, node.ApiAddress, node.ClusterAddress, "Standby"))
		}
		c.UI.Output(tableOutput(out, nil))
		return 0
	default:
		return OutputData(c.UI, result)
	}
}

type NodeAddr struct {
	HostName       string `json:"host_name"`
	ApiAddress     string `json:"api_addr"`
	ClusterAddress string `json:"cluster_addr"`
}

type HaStatusResponse struct {
	Leader      NodeAddr   `json:"leader"`
	PerfStandby []NodeAddr `json:"performance_standby"`
	Standby     []NodeAddr `json:"standby"`
}
