// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"fmt"

	slackpkg "github.com/hashicorp/vault/tools/pipeline/internal/pkg/slack"
	"github.com/spf13/cobra"
)

var sendSlackMsgReq = &slackpkg.SendSlackMsgReq{}

func newSlackPackageSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send --channel <channel-id> --message <text> [--slack-bot-token <token>] [--emoji <emoji>] [--thread-ts <timestamp>]",
		Short: "Send a Slack message",
		Long:  "Send a message to a Slack channel with optional threading and custom emoji",
		Args:  cobra.NoArgs,
		RunE:  runSendSlackMsg,
	}

	cmd.Flags().StringVar(&sendSlackMsgReq.ChannelID, "channel", "", "The Slack channel ID to send the message to")
	cmd.Flags().StringVar(&sendSlackMsgReq.MessageText, "message", "", "The text content of the Slack message")
	cmd.Flags().StringVar(&sendSlackMsgReq.Token, "slack-bot-token", "", "Slack bot token with chat:write permissions (defaults to SLACK_BOT_TOKEN)")
	cmd.Flags().StringVar(&sendSlackMsgReq.Emoji, "emoji", "", "The emoji to use as the message icon (defaults to :vault:)")
	cmd.Flags().StringVar(&sendSlackMsgReq.ThreadTimestamp, "thread-ts", "", "The timestamp of the parent message to reply in a thread")
	cmd.Flags().BoolVar(&sendSlackMsgReq.WriteToGithubOutput, "github-output", false, "Whether or not to write 'thread-ts' to $GITHUB_OUTPUT")

	_ = cmd.MarkFlagRequired("channel")
	_ = cmd.MarkFlagRequired("message")

	return cmd
}

func runSendSlackMsg(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	res, err := sendSlackMsgReq.Run(cmd.Context())
	if err != nil {
		return err
	}

	b, err := res.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	if sendSlackMsgReq.WriteToGithubOutput {
		output, err := res.ToGithubOutput()
		if err != nil {
			return err
		}

		return writeToGithubOutput("thread-ts", output)
	}

	return nil
}
