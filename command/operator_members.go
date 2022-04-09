package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorMembersCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorMembersCommand)(nil)
)

type OperatorMembersCommand struct {
	*BaseCommand
}

func (c *OperatorMembersCommand) Synopsis() string {
	return "Returns the list of nodes in the cluster"
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

	resp, err := client.Sys().HAStatus()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	switch Format(c.UI) {
	case "table":
		out := []string{"Host Name | API Address | Cluster Address | ActiveNode | Last Echo"}
		for _, node := range resp.Nodes {
			out = append(out, fmt.Sprintf("%s | %s | %s | %t | %s", node.Hostname, node.APIAddress, node.ClusterAddress, node.ActiveNode, node.LastEcho))
		}
		c.UI.Output(tableOutput(out, nil))
		return 0
	default:
		return OutputData(c.UI, resp)
	}
}
