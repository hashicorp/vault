// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import "github.com/spf13/cobra"

func newSlackCmd() *cobra.Command {
	slackCmd := &cobra.Command{
		Use:   "slack",
		Short: "Slack related tasks",
		Long:  "Slack related tasks",
	}
	slackCmd.AddCommand(newSlackLookupIDCmd())
	slackCmd.AddCommand(newSlackPackageSendCmd())

	return slackCmd
}
