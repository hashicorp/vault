// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*OperatorRaftJoinCommand)(nil)
	_ cli.CommandAutocomplete = (*OperatorRaftJoinCommand)(nil)
)

type OperatorRaftJoinCommand struct {
	flagRetry            bool
	flagNonVoter         bool
	flagLeaderCACert     string
	flagLeaderClientCert string
	flagLeaderClientKey  string
	flagAutoJoinScheme   string
	flagAutoJoinPort     uint
	*BaseCommand
}

func (c *OperatorRaftJoinCommand) Synopsis() string {
	return "Joins a node to the Raft cluster"
}

func (c *OperatorRaftJoinCommand) Help() string {
	helpText := `
Usage: vault operator raft join [options] <leader-api-addr|auto-join-configuration>

  Join the current node as a peer to the Raft cluster by providing the address
  of the Raft leader node.

      $ vault operator raft join "http://127.0.0.2:8200"

  Join the current node as a peer to the Raft cluster by providing cloud auto-join
  configuration.

      $ vault operator raft join "provider=aws region=eu-west-1 ..."
			
  Join the current node as a peer to the Raft cluster by providing cloud auto-join
  configuration with an explicit URI scheme and port.

			$ vault operator raft join -auto-join-scheme="http" -auto-join-port=8201 \
			  "provider=aws region=eu-west-1 ..."

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
		Name:       "auto-join-scheme",
		Target:     &c.flagAutoJoinScheme,
		Completion: complete.PredictNothing,
		Default:    "https",
		Usage:      "An optional URI protocol scheme used for addresses discovered via auto-join.",
	})

	f.UintVar(&UintVar{
		Name:       "auto-join-port",
		Target:     &c.flagAutoJoinPort,
		Completion: complete.PredictNothing,
		Default:    8200,
		Usage:      "An optional port used for addresses discovered via auto-join.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-ca-cert",
		Target:     &c.flagLeaderCACert,
		Completion: complete.PredictNothing,
		Usage:      "CA cert to use when verifying the Raft leader certificate.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-client-cert",
		Target:     &c.flagLeaderClientCert,
		Completion: complete.PredictNothing,
		Usage:      "Client cert to use when authenticating with the Raft leader.",
	})

	f.StringVar(&StringVar{
		Name:       "leader-client-key",
		Target:     &c.flagLeaderClientKey,
		Completion: complete.PredictNothing,
		Usage:      "Client key to use when authenticating with the Raft leader.",
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

	leaderInfo := ""

	args = f.Args()
	switch len(args) {
	case 0:
		// No-op: This is acceptable if we're using raft for HA-only
	case 1:
		leaderInfo = strings.TrimSpace(args[0])
	default:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0-1, got %d)", len(args)))
		return 1
	}

	leaderCACert, err := parseFlagFile(c.flagLeaderCACert)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse leader CA certificate: %s", err))
		return 1
	}

	leaderClientCert, err := parseFlagFile(c.flagLeaderClientCert)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse leader client certificate: %s", err))
		return 1
	}

	leaderClientKey, err := parseFlagFile(c.flagLeaderClientKey)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse leader client key: %s", err))
		return 1
	}

	if c.flagAutoJoinScheme != "" && (c.flagAutoJoinScheme != "http" && c.flagAutoJoinScheme != "https") {
		c.UI.Error(fmt.Sprintf("invalid scheme %q; must either be http or https", c.flagAutoJoinScheme))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	joinReq := &api.RaftJoinRequest{
		LeaderCACert:     leaderCACert,
		LeaderClientCert: leaderClientCert,
		LeaderClientKey:  leaderClientKey,
		Retry:            c.flagRetry,
		NonVoter:         c.flagNonVoter,
	}

	if strings.Contains(leaderInfo, "provider=") {
		joinReq.AutoJoin = leaderInfo
		joinReq.AutoJoinScheme = c.flagAutoJoinScheme
		joinReq.AutoJoinPort = c.flagAutoJoinPort
	} else {
		joinReq.LeaderAPIAddr = leaderInfo
	}

	resp, err := client.Sys().RaftJoin(joinReq)
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
