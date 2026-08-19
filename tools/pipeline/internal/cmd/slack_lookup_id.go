// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"
	"strings"

	slackpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/slack"
	"github.com/spf13/cobra"
)

var lookupIDReq = &slackpkg.LookupIDByEmailReq{
	AllowedDomain: "ibm.com",
}

func newSlackLookupIDCmd() *cobra.Command {
	var emails string

	cmd := &cobra.Command{
		Use:   "lookup-id --emails \"user1@ibm.com user2@ibm.com\" [--slack-token <token>]",
		Short: "Lookup Slack user IDs by IBM email address",
		Long:  "Lookup Slack user IDs by IBM email address",
		Args:  cobra.NoArgs,
		RunE:  runLookupIDReq,
	}

	cmd.Flags().StringVar(&emails, "emails", "", "A space-separated list of @ibm.com email addresses to lookup")
	cmd.Flags().StringVar(&lookupIDReq.Token, "slack-token", "", "Slack API token with permissions to read user info (defaults to SLACK_TOKEN)")
	cmd.Flags().BoolVar(&lookupIDReq.WriteToGithubOutput, "github-output", false, "Whether or not to write 'email-slack-id-map' to $GITHUB_OUTPUT")

	_ = cmd.MarkFlagRequired("emails")

	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		lookupIDReq.Emails = strings.Fields(emails)
	}

	return cmd
}

func runLookupIDReq(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, err := lookupIDReq.Run(cmd.Context())
	if err != nil {
		return err
	}

	b, err := res.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	if lookupIDReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("email-slack-id-map", output)
	}

	return nil
}
