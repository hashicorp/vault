package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorRaftJoinCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorRaftJoinCommand)(nil)

type OperatorRaftJoinCommand struct {
	flagRetry            bool
	flagLeaderCACert     string
	flagLeaderClientCert string
	flagLeaderClientKey  string
	flagNonVoter         bool
	*BaseCommand
}

func (c *OperatorRaftJoinCommand) Synopsis() string {
	return "Joins a node to the Raft cluster"
}

func (c *OperatorRaftJoinCommand) Help() string {
	helpText := `
Usage: vault operator raft join [options] <leader-api-addr>

  Join the current node as a peer to the Raft cluster by providing the address
  of the Raft leader node.

      $ vault operator raft join "http://127.0.0.2:8200"

  TLS certificate data can also be consumed from a file on disk by prefixing with
  the "@" symbol. For example:

      $ vault operator raft join "http://127.0.0.2:8200" \
        -leader-ca-cert=@leader_ca.crt \
        -leader-client-cert=@leader_client.crt \
        -leader-client-key=@leader.key

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
		Usage:      "CA cert to communicate with Raft leader.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-client-cert",
		Target:     &c.flagLeaderClientCert,
		Completion: complete.PredictNothing,
		Usage:      "Client cert to to authenticate to Raft leader.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-client-key",
		Target:     &c.flagLeaderClientKey,
		Completion: complete.PredictNothing,
		Usage:      "Client key to to authenticate to Raft leader.",
	})

	f.BoolVar(&BoolVar{
		Name:    "retry",
		Target:  &c.flagRetry,
		Default: false,
		Usage:   "Continuously retry joining the Raft cluster upon failures.",
	})

	f.BoolVar(&BoolVar{
		Name:    "non-voter",
		Target:  &c.flagNonVoter,
		Default: false,
		Usage:   "(Enterprise-only) This flag is used to make the server not participate in the Raft quorum, and have it only receive the data replication stream. This can be used to add read scalability to a cluster in cases where a high volume of reads to servers are needed.",
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

	leaderCACert, err := parseArg(c.flagLeaderCACert)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse leader CA certificate: %s", err))
		return 1
	}

	leaderClientCert, err := parseArg(c.flagLeaderClientCert)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse leader client certificate: %s", err))
		return 1
	}

	leaderClientKey, err := parseArg(c.flagLeaderClientKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse leader client key: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	resp, err := client.Sys().RaftJoin(&api.RaftJoinRequest{
		LeaderAPIAddr:    leaderAPIAddr,
		LeaderCACert:     leaderCACert,
		LeaderClientCert: leaderClientCert,
		LeaderClientKey:  leaderClientKey,
		Retry:            c.flagRetry,
		NonVoter:         c.flagNonVoter,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error joining the node to the Raft cluster: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
	default:
		return OutputData(c.UI, resp)
	}

	out := []string{
		"Key | Value",
		fmt.Sprintf("Joined | %t", resp.Joined),
	}
	c.UI.Output(tableOutput(out, nil))

	return 0
}

func parseArg(raw string) (string, error) {
	// check if the provided argument should be read from file
	if len(raw) > 0 && raw[0] == '@' {
		contents, err := ioutil.ReadFile(raw[1:])
		if err != nil {
			return "", errwrap.Wrapf("error reading file: {{err}}", err)
		}

		return string(contents), nil
	}

	return raw, nil
}
