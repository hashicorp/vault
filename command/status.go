package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

// Ensure we are implementing the right interfaces.
var _ cli.Command = (*StatusCommand)(nil)
var _ cli.CommandAutocomplete = (*StatusCommand)(nil)

// StatusCommand is a Command that outputs the status of whether Vault is sealed
// or not as well as HA information.
type StatusCommand struct {
	*BaseCommand
}

func (c *StatusCommand) Synopsis() string {
	return "Prints seal and HA status"
}

func (c *StatusCommand) Help() string {
	helpText := `
Usage: vault status [options]

  Prints the current state of Vault including whether it is sealed and if HA
  mode is enabled. This command prints regardless of whether the Vault is
  sealed.

  The exit code reflects the seal status:

      - 0 - unsealed
      - 1 - error
      - 2 - sealed

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *StatusCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *StatusCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *StatusCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *StatusCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		// We return 2 everywhere else, but 2 is reserved for "sealed" here
		return 1
	}

	sealStatus, err := client.Sys().SealStatus()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error checking seal status: %s", err))
		return 1
	}

	outStr := fmt.Sprintf(
		"Sealed: %v\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Unseal Progress: %d\n"+
			"Unseal Nonce: %v\n"+
			"Version: %s",
		sealStatus.Sealed,
		sealStatus.N,
		sealStatus.T,
		sealStatus.Progress,
		sealStatus.Nonce,
		sealStatus.Version)

	if sealStatus.ClusterName != "" && sealStatus.ClusterID != "" {
		outStr = fmt.Sprintf("%s\nCluster Name: %s\nCluster ID: %s", outStr, sealStatus.ClusterName, sealStatus.ClusterID)
	}

	c.UI.Output(outStr)

	// Mask the 'Vault is sealed' error, since this means HA is enabled, but that
	// we cannot query for the leader since we are sealed.
	leaderStatus, err := client.Sys().Leader()
	if err != nil && strings.Contains(err.Error(), "Vault is sealed") {
		leaderStatus = &api.LeaderResponse{HAEnabled: true}
		err = nil
	}
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error checking leader status: %s", err))
		return 1
	}

	// Output if HA is enabled
	c.UI.Output("")
	c.UI.Output(fmt.Sprintf("High-Availability Enabled: %v", leaderStatus.HAEnabled))
	if leaderStatus.HAEnabled {
		if sealStatus.Sealed {
			c.UI.Output("\tMode: sealed")
		} else {
			mode := "standby"
			if leaderStatus.IsSelf {
				mode = "active"
			}
			c.UI.Output(fmt.Sprintf("\tMode: %s", mode))

			if leaderStatus.LeaderAddress == "" {
				leaderStatus.LeaderAddress = "<none>"
			}
			if leaderStatus.LeaderClusterAddress == "" {
				leaderStatus.LeaderClusterAddress = "<none>"
			}
			c.UI.Output(fmt.Sprintf("\tLeader Cluster Address: %s", leaderStatus.LeaderClusterAddress))
		}
	}

	if sealStatus.Sealed {
		return 2
	}

	return 0
}
