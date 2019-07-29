package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorRaftJoinCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorRaftJoinCommand)(nil)

type OperatorRaftJoinCommand struct {
	flagRaftRetry        bool
	flagLeaderCACert     string
	flagLeaderClientCert string
	flagLeaderClientKey  string
	*BaseCommand
}

func (c *OperatorRaftJoinCommand) Synopsis() string {
	return "Joins a node to the raft cluster"
}

func (c *OperatorRaftJoinCommand) Help() string {
	helpText := `
Usage: vault operator raft join [options] <leader-api-addr>

  Join the current node as a peer to the raft cluster by providing the address
  of the raft leader node.

	  $ vault operator raft join "http://127.0.0.2:8200"

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *OperatorRaftJoinCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "leader-ca-cert",
		Target:     &c.flagLeaderCACert,
		Completion: complete.PredictNothing,
		Usage:      "CA cert to communicate with raft leader.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-client-cert",
		Target:     &c.flagLeaderClientCert,
		Completion: complete.PredictNothing,
		Usage:      "Client cert to to authenticate to raft leader.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-client-key",
		Target:     &c.flagLeaderClientKey,
		Completion: complete.PredictNothing,
		Usage:      "Client key to to authenticate to raft leader.",
	})

	f.BoolVar(&BoolVar{
		Name:    "retry",
		Target:  &c.flagRaftRetry,
		Default: false,
		Usage:   "Continuously retry joining the raft cluster upon failures.",
	})

	return set
}

func (c *OperatorRaftJoinCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRaftJoinCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRaftJoinCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	leaderAPIAddr := ""

	args = f.Args()
	switch len(args) {
	case 1:
		leaderAPIAddr = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Incorrect arguments (expected 1, got %d)", len(args)))
		return 1
	}

	if len(leaderAPIAddr) == 0 {
		c.UI.Error("leader api address is required")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Sys().RaftJoin(&api.RaftJoinRequest{
		LeaderAPIAddr:    leaderAPIAddr,
		LeaderCACert:     c.flagLeaderCACert,
		LeaderClientCert: c.flagLeaderClientCert,
		LeaderClientKey:  c.flagLeaderClientKey,
		Retry:            c.flagRaftRetry,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error joining the node to the raft cluster: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
	default:
		return OutputData(c.UI, resp)
	}

	out := []string{}
	out = append(out, "Key | Value")
	out = append(out, fmt.Sprintf("Joined | %t", resp.Joined))
	c.UI.Output(tableOutput(out, nil))

	return 0
}
